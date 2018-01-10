package admin_test

import (
	"context"

	"github.com/zenoss/zenkit"
	. "github.com/zenoss/zenkit/admin"
	"github.com/zenoss/zenkit/logging"
	"github.com/zenoss/zenkit/metrics"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Context", func() {

	var ctx context.Context

	BeforeEach(func() {
		ctx = context.Background()
	})

	It("should return a nil service if no parent is on the context", func() {
		Ω(ContextParentService(ctx)).Should(BeNil())
	})

	Context("with a parent service on the context", func() {

		var service = zenkit.NewService("test-service")

		BeforeEach(func() {
			service.Context = metrics.WithMetrics(ctx, &Registry{})
			ctx = WithParentService(ctx, service)
		})

		It("should return the service from the context", func() {
			Ω(ContextParentService(ctx)).Should(Equal(service))
		})

		It("should return the logger from the parent serivice", func() {
			logger := logging.ContextLogger(service.Context)
			Ω(ContextLogger(ctx)).Should(Equal(logger))
		})

		It("should return the metrics from the parent service", func() {
			registry := metrics.ContextMetrics(service.Context)
			Ω(ContextMetrics(ctx)).Should(Equal(registry))
		})
	})
})
