package metrics_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/goadesign/goa"
	metrics "github.com/rcrowley/go-metrics"
	"github.com/zenoss/zenkit/logging"
	. "github.com/zenoss/zenkit/metrics"
        "github.com/zenoss/zenkit/test"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TimedFunc(ctx context.Context) {
	defer MeasureTime(ctx)()
	time.Sleep(1 * time.Millisecond)
}

func TimedFuncWithName(ctx context.Context, prefix string) {
	defer MeasureTimeNamed(ctx, prefix)()
	time.Sleep(1 * time.Millisecond)
}

var _ = Describe("Metrics", func() {

	It("should be able to be registered on and retrieved from the context", func() {
		ctx := context.Background()
		reg := metrics.NewRegistry()
		newctx := WithMetrics(ctx, reg)
		newreg := ContextMetrics(newctx)
		Ω(newreg).Should(BeIdenticalTo(reg))
	})

	It("should not fail if no metric registry exists", func() {
		ctx := context.Background()
		reg := ContextMetrics(ctx)
		Ω(reg).Should(BeNil())
	})

	It("should be able to be registered by middleware", func() {
		var (
			m metrics.Registry
			resp = httptest.NewRecorder()
			req, _ = http.NewRequest("", "http://example.com", nil)
		)

		handler := func(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
			m = ContextMetrics(ctx)
			return nil
		}

		mw := MetricsMiddleware()
		wrapped := mw(handler)
		err := wrapped(context.Background(), resp, req)
		Ω(err).ShouldNot(HaveOccurred())
		Ω(m).ShouldNot(BeNil())
	})

	It("should be able to measure the time a function takes", func() {
		reg := metrics.NewRegistry()
		ctx := WithMetrics(context.Background(), reg)
		TimedFunc(ctx)
		metric := reg.Get("func.TimedFunc.time")
		Ω(metric).ShouldNot(BeNil())
	})

	It("should be able to measure the time for a named timer", func() {
		reg := metrics.NewRegistry()
		ctx := WithMetrics(context.Background(), reg)
		TimedFuncWithName(ctx, "my.timer")
		metric := reg.Get("my.timer.func.TimedFuncWithName.time")
		Ω(metric).ShouldNot(BeNil())
	})

	It("should not panic function metrics if no registry is defined", func() {
		ctx := WithMetrics(context.Background(), nil)
		TimedFunc(ctx)
	})

	It("should be able to increment a counter", func() {
		reg := metrics.NewRegistry()
		ctx := WithMetrics(context.Background(), reg)
		IncrementCounter(ctx, "my.test.counter", 1)
		metric := reg.Get("my.test.counter")
		Ω(metric).ShouldNot(BeNil())
	})

	It("should not panic counter increment if no registry is defined", func() {
		IncrementCounter(context.Background(), "my.test.counter", 1)
	})

	It("should be able to update a meter", func() {
		reg := metrics.NewRegistry()
		ctx := WithMetrics(context.Background(), reg)
		UpdateMeter(ctx, "my.test.meter", 1)
		metric := reg.Get("my.test.meter")
		Ω(metric).ShouldNot(BeNil())
	})

	It("should not panic meter update if no registry is defined", func() {
		UpdateMeter(context.Background(), "my.test.meter", 1)
	})

	Context("with CollectMetrics", func() {
		var (
			registry metrics.Registry
			svc      *goa.Service
			ctx      context.Context
			cancel   chan struct{}
		)

		BeforeEach(func() {
			svc = goa.New(test.RandString(8))
			svc.WithLogger(logging.ServiceLogger())
			registry = metrics.NewRegistry()
			ctx = WithMetrics(svc.Context, registry)
			cancel = make(chan struct{})
		})

		AfterEach(func() {
			close(cancel)
		})

		It("should not panic if no registry is defined", func() {
			err := CollectMetrics(svc.Context, time.Second, "prefix", cancel, func(metrics []Metric) {})
			Ω(err).Should(HaveOccurred())
			Ω(err).Should(MatchError("metrics registry not defined"))
		})

		It("should return no metrics if nothing has been added to metric registry", func() {
			nReceived := 0
			results   := []Metric{}
			receiver  := func(metrics []Metric) {
				nReceived += 1
				results = append(results, metrics...)
			}
			err := CollectMetrics(ctx, time.Second, "prefix", cancel, receiver)
			Ω(err).Should(Succeed())

			// wait for at least 1 collection cycle
			time.Sleep(2 * time.Second)

			// the receiver should have been called, but no metrics collected
			Ω(nReceived).Should(Equal(0))
			Ω(len(results)).Should(Equal(0))
		})

		It("should return some Counter metrics if a counter has been incremented", func() {
			nReceived := 0
			results   := []Metric{}
			receiver  := func(metrics []Metric) {
				nReceived += 1
				results = append(results, metrics...)
			}
			err := CollectMetrics(ctx, time.Second, "prefix", cancel, receiver)
			Ω(err).Should(Succeed())

			IncrementCounter(ctx, "my.counter", 1)

			// wait for at least 1 collection cycle
			time.Sleep(2 * time.Second)

			// the receiver should have been called, but no metrics collected
			Ω(nReceived).ShouldNot(Equal(0))
			Ω(len(results)).ShouldNot(Equal(0))
			Ω(results[0].Metric).Should(Equal("prefix.my.counter.count"))
		})

		It("should return some Meter metrics if a meter has been incremented", func() {
			nReceived := 0
			results   := []Metric{}
			receiver  := func(metrics []Metric) {
				nReceived += 1
				results = append(results, metrics...)
			}
			err := CollectMetrics(ctx, time.Second, "prefix", cancel, receiver)
			Ω(err).Should(Succeed())

			UpdateMeter(ctx, "my.meter", 1)

			// wait for at least 1 collection cycle
			time.Sleep(2 * time.Second)

			// the receiver should have been called, but no metrics collected
			Ω(nReceived).ShouldNot(Equal(0))
			Ω(len(results)).ShouldNot(Equal(0))
			Ω(results[0].Metric).Should(Equal("prefix.my.meter.count"))
		})

		It("should return some Timer metrics if a timer has been incremented", func() {
			nReceived := 0
			results   := []Metric{}
			receiver  := func(metrics []Metric) {
				nReceived += 1
				results = append(results, metrics...)
			}
			err := CollectMetrics(ctx, time.Second, "prefix", cancel, receiver)
			Ω(err).Should(Succeed())

			MeasureTime(ctx)()

			// wait for at least 1 collection cycle
			time.Sleep(2 * time.Second)

			// the receiver should have been called, but no metrics collected
			Ω(nReceived).ShouldNot(Equal(0))
			Ω(len(results)).ShouldNot(Equal(0))
			Ω(results[0].Metric).Should(ContainSubstring("prefix"))
			Ω(results[0].Metric).Should(ContainSubstring("time"))
		})

		It("should return some metrics even if prefix is empty", func() {
			nReceived := 0
			results   := []Metric{}
			receiver  := func(metrics []Metric) {
				nReceived += 1
				results = append(results, metrics...)
			}
			err := CollectMetrics(ctx, time.Second, "", cancel, receiver)
			Ω(err).Should(Succeed())

			IncrementCounter(ctx, "my.counter", 1)

			// wait for at least 1 collection cycle
			time.Sleep(2 * time.Second)

			// the receiver should have been called, but no metrics collected
			Ω(nReceived).ShouldNot(Equal(0))
			Ω(len(results)).ShouldNot(Equal(0))
			Ω(results[0].Metric).Should(Equal("my.counter.count"))
		})

	})

})
