package metrics

import (
	"context"
	"fmt"
	"net/http"

	"github.com/goadesign/goa"
	metrics "github.com/rcrowley/go-metrics"
	"github.com/zenoss/zenkit/funcname"
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
	begin := TimeFunc()
	fn := funcname.FuncName(2)
	registry := ContextMetrics(ctx)
	exit := func() {
		if registry == nil {
			return
		}
		t := metrics.GetOrRegisterTimer(fmt.Sprintf("func.%s.time", fn), registry)
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

func MetricsMiddleware() goa.Middleware {
	m := metrics.NewRegistry()
	return func(h goa.Handler) goa.Handler {
		return func(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
			return h(WithMetrics(ctx, m), rw, req)
		}
	}
}
