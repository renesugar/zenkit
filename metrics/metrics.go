package metrics

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/goadesign/goa"
	metrics "github.com/rcrowley/go-metrics"
	"github.com/zenoss/zenkit/funcname"
	"github.com/pkg/errors"
	"github.com/zenoss/zenkit/logging"
)

type key int

const (
	metricsKey key = iota + 1
)

func WithMetrics(ctx context.Context, registry metrics.Registry) context.Context {
	return context.WithValue(ctx, metricsKey, registry)
}

func ContextMetrics(ctx context.Context) metrics.Registry {
	if v := ctx.Value(metricsKey); v != nil {
		return v.(metrics.Registry)
	}
	return nil
}

func MeasureTime(ctx context.Context) func() {
	return measureTimeNamed(ctx, "")
}

func MeasureTimeNamed(ctx context.Context, prefix string) func() {
	return measureTimeNamed(ctx, prefix)
}

func measureTimeNamed(ctx context.Context, prefix string) func() {
	begin := TimeFunc()
	fn := funcname.FuncName(3)
	var name string
	if len(prefix) == 0 {
		name = fmt.Sprintf("func.%s.time", fn)
	} else {
		name = fmt.Sprintf("%s.func.%s.time", prefix, fn)
	}
	registry := ContextMetrics(ctx)
	exit := func() {
		if registry == nil {
			return
		}
		t := metrics.GetOrRegisterTimer(name, registry)
		t.UpdateSince(begin)
	}
	return exit
}

func IncrementCounter(ctx context.Context, name string, inc int64) {
	registry := ContextMetrics(ctx)
	if registry == nil {
		return
	}
	ctr := metrics.GetOrRegisterCounter(name, registry)
	ctr.Inc(inc)
}

func UpdateMeter(ctx context.Context, name string, n int64) {
	registry := ContextMetrics(ctx)
	if registry == nil {
		return
	}
	ctr := metrics.GetOrRegisterMeter(name, registry)
	ctr.Mark(n)
}

func MetricsMiddleware() goa.Middleware {
	return func(h goa.Handler) goa.Handler {
		return func(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
			// if the parent context has a metrics registry already, then use it
			m := ContextMetrics(ctx)
			if m == nil {
				ctx = WithMetrics(ctx, metrics.NewRegistry())
			}
			return h(ctx, rw, req)
		}
	}
}

// Metric defines the structure of a single metric returned from CollectMetrics
type Metric struct {
	Timestamp float64                `json:"timestamp"`
	Metric    string                 `json:"metric"`
	Value     float64                `json:"value"`
}

// MetricReceiver defines the callback function invoked periodically by CollectMetrics
type MetricReceiver func(metrics []Metric)

// CollectMetrics periodically collects all codahale metrics in the current metric registry,
// constructs a list of zenkit Metric values derived from each codahale Metric, and
// calls the receiver function with the list of zenkit Metric values.
//
// freq is the rate at which the metrics will be collected and the receiver function will be called
//
// prefix is an optional string that will be used to prefix each name for the zenkit Metric. Metric names in a
// zenkit Metric will have the format [prefix.]name.valueName where name is the name of the codahale Metric and
// valueName describes a specific value from that Metric (e.g. 'count', 'min', 'max' etc).
//
// cancel is a channel the caller can use to cancel the periodic collection of metrics.
//
// receiver is the user-defined function tha will be called with the generated list of zenkit Metric values
//
func CollectMetrics(ctx context.Context, freq time.Duration, prefix string, cancel <-chan struct{}, receiver MetricReceiver) error {
	logger := logging.ContextLogger(ctx).Logger

	registry := ContextMetrics(ctx)
	if registry == nil {
		logger.Error("CollectMetrics called before metrics Registry was added to Context")
		return errors.New("metrics registry not defined")
	}

	go func() {
		timer := time.NewTicker(freq)
		for {
			select {
			case <-timer.C:
				metricList := []Metric{}
				timestamp := float64(time.Now().Unix())
				registry.Each(func(name string, m interface{}) {
					metricList = addMetrics(prefix, name, m, timestamp, metricList)
				})
				if len(metricList) == 0 {
					logger.Debug("CollectMetrics - nothing to report")
					continue
				}
				receiver(metricList)
			case <-cancel:
				logger.Info("CollectMetrics cancelled")
				return
			}
		}
	}()
	return nil
}

func addMetrics(prefix, name string, m interface{}, timestamp float64, metricList []Metric) []Metric {
	switch codahaleMetric := m.(type) {
	case metrics.Counter:
		metricList = append(metricList, toMetric(metricName(prefix, name, "count"), float64(codahaleMetric.Count()), timestamp))

	case metrics.Meter:
		metricList = append(metricList, toMetric(metricName(prefix, name, "count"), float64(codahaleMetric.Count()), timestamp))
		metricList = append(metricList, toMetric(metricName(prefix, name, "1MinuteRate"), codahaleMetric.Rate1(), timestamp))
		metricList = append(metricList, toMetric(metricName(prefix, name, "5MinuteRate"), codahaleMetric.Rate5(), timestamp))
		metricList = append(metricList, toMetric(metricName(prefix, name, "15MinuteRate"), codahaleMetric.Rate15(), timestamp))
		metricList = append(metricList, toMetric(metricName(prefix, name, "meanRate"), codahaleMetric.RateMean(), timestamp))

	case metrics.Timer:
		/*
		 * report the same values reported for Timers in MetricConsumer, adjusting the units for
		 *    time values to milliseconds instead of nanoseconds
		 */
		s := codahaleMetric.Snapshot()
		units := time.Millisecond
		unitDivisor := float64(units)
		metricList = append(metricList, toMetric(metricName(prefix, name, "min"), float64(s.Min())/unitDivisor, timestamp))
		metricList = append(metricList, toMetric(metricName(prefix, name, "max"), float64(s.Max())/unitDivisor, timestamp))
		metricList = append(metricList, toMetric(metricName(prefix, name, "mean"), float64(s.Mean()/unitDivisor), timestamp))
		metricList = append(metricList, toMetric(metricName(prefix, name, "stddev"), s.StdDev()/unitDivisor, timestamp))

		ps := s.Percentiles([]float64{0.5, 0.75, 0.95, 0.98, 0.99, 0.999})
		metricList = append(metricList, toMetric(metricName(prefix, name, "median"),  ps[0]/unitDivisor, timestamp))
		metricList = append(metricList, toMetric(metricName(prefix, name, "p75"),  ps[1]/unitDivisor, timestamp))
		metricList = append(metricList, toMetric(metricName(prefix, name, "p95"),  ps[2]/unitDivisor, timestamp))
		metricList = append(metricList, toMetric(metricName(prefix, name, "p98"),  ps[3]/unitDivisor, timestamp))
		metricList = append(metricList, toMetric(metricName(prefix, name, "p99"),  ps[4]/unitDivisor, timestamp))
		metricList = append(metricList, toMetric(metricName(prefix, name, "p999"),  ps[5]/unitDivisor, timestamp))

		// Add the value which overlap with metrics.Meter
		metricList = append(metricList, toMetric(metricName(prefix, name, "count"), float64(s.Count()), timestamp))
		metricList = append(metricList, toMetric(metricName(prefix, name, "meanRate"), s.RateMean(), timestamp))
		metricList = append(metricList, toMetric(metricName(prefix, name, "1MinuteRate"), s.Rate1(), timestamp))
		metricList = append(metricList, toMetric(metricName(prefix, name, "5MinuteRate"), s.Rate5(), timestamp))
		metricList = append(metricList, toMetric(metricName(prefix, name, "15MinuteRate"), s.Rate15(), timestamp))
	}
	return metricList
}

func metricName(prefix, name, valueName string) string {
	if len(prefix) == 0 {
		return fmt.Sprintf("%s.%s", name, valueName)
	}
	return fmt.Sprintf("%s.%s.%s", prefix, name, valueName)
}

// toMetric creates a zenkit Metric from a name and value
func toMetric(name string, value float64, timestamp float64) Metric {
	metric := Metric{}
	metric.Metric = name
	metric.Timestamp = timestamp
	metric.Value = value
	return metric
}
