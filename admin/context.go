package admin

import (
	"context"

	"github.com/goadesign/goa"
	gometrics "github.com/rcrowley/go-metrics"
	"github.com/sirupsen/logrus"
	"github.com/zenoss/zenkit/logging"
	"github.com/zenoss/zenkit/metrics"
)

type key int

const (
	serviceKey key = iota + 1
)

func WithParentService(ctx context.Context, service *goa.Service) context.Context {
	return context.WithValue(ctx, serviceKey, service)
}

func ContextParentService(ctx context.Context) *goa.Service {
	v := ctx.Value(serviceKey)
	s, ok := v.(*goa.Service)
	if !ok {
		return nil
	}
	return s
}

func ContextLogger(ctx context.Context) *logrus.Entry {
	s := ContextParentService(ctx)
	return logging.ContextLogger(s.Context)
}

func ContextMetrics(ctx context.Context) gometrics.Registry {
	s := ContextParentService(ctx)
	return metrics.ContextMetrics(s.Context)
}
