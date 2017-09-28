package metrics

import (
	"context"
	"net/http"

	"github.com/goadesign/goa"
	metrics "github.com/rcrowley/go-metrics"
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

func MetricsMiddleware() goa.Middleware {
	m := metrics.NewRegistry()
	return func(h goa.Handler) goa.Handler {
		return func(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
			return h(WithMetrics(ctx, m), rw, req)
		}
	}
}
