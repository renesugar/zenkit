package admin_test

import (
	"context"
	"errors"

	healthcheck "github.com/docker/go-healthcheck"
	"github.com/goadesign/goa"
	gometrics "github.com/rcrowley/go-metrics"
	"github.com/zenoss/zenkit"
	. "github.com/zenoss/zenkit/admin"
	"github.com/zenoss/zenkit/admin/app/test"
	"github.com/zenoss/zenkit/metrics"

	. "github.com/onsi/ginkgo"
	//. "github.com/onsi/gomega"
)

// We need a registry that refuses to Marshal
type Registry struct {
	gometrics.StandardRegistry
}

func (r *Registry) MarshalJSON() ([]byte, error) {
	return []byte(nil), errors.New("this is a test")
}

var _ = Describe("Admin", func() {

	var (
		t      = GinkgoT()
		ctx    context.Context
		parent *goa.Service
		svc    = goa.New("admin-test")
		ctrl   = NewAdminController(svc)
	)

	BeforeEach(func() {
		ctx = context.Background()
		parent = zenkit.NewService("test-service", false)
	})

	JustBeforeEach(func() {
		ctx = WithParentService(ctx, parent)
	})

	Context("when the Metrics resource is requested", func() {

		Context("when the metrics middleware is hooked up", func() {
			BeforeEach(func() {
				registry := gometrics.NewRegistry()
				parent.Context = metrics.WithMetrics(parent.Context, registry)
			})
			Context("when the registry cannot be encoded", func() {
				BeforeEach(func() {
					parent.Context = metrics.WithMetrics(parent.Context, &Registry{})
				})
				It("should produce an error", func() {
					test.MetricsAdminInternalServerError(t, ctx, svc, ctrl, true)
				})
			})
			Context("when the registry can be encoded", func() {
				It("should respond OK", func() {
					test.MetricsAdminOK(t, ctx, svc, ctrl, true)
				})
			})
		})

		Context("when the metrics middleware isn't hooked up", func() {
			It("should respond OK", func() {
				test.MetricsAdminOK(t, ctx, svc, ctrl, false)
			})
		})
	})

	Context("when the Ping resource is requested", func() {
		It("should respond OK", func() {
			test.PingAdminOK(t, ctx, svc, ctrl)
		})
	})

	Context("when the Swagger resource is requested", func() {
		It("should respond OK", func() {
			test.SwaggerAdminOK(t, ctx, svc, ctrl)
		})
	})

	Context("when the SwaggerJSON resource is requested", func() {
		originalAsset := SwaggerJSONAsset
		Context("when the swagger.json asset is missing", func() {
			BeforeEach(func() {
				SwaggerJSONAsset = "none"
			})
			AfterEach(func() {
				SwaggerJSONAsset = originalAsset
			})
			It("should respond with an InternalServerError", func() {
				test.SwaggerJSONAdminInternalServerError(t, ctx, svc, ctrl)
			})
		})
		Context("when the swagger.json asset is available", func() {
			It("should respond OK", func() {
				test.SwaggerJSONAdminOK(t, ctx, svc, ctrl)
			})
		})
	})

	Context("when the Health resource is requested", func() {

		AfterEach(func() {
			healthcheck.DefaultRegistry = healthcheck.NewRegistry()
		})

		It("should return OK if there are no failing health checks", func() {
			check := func() error { return nil }
			healthcheck.RegisterFunc("testOK", check)
			test.HealthAdminOK(t, ctx, svc, ctrl)
		})

		It("should return ServiceUnavailable if there are failing health checks", func() {
			check := func() error { return errors.New("he dead") }
			healthcheck.RegisterFunc("testDOWN", check)
			test.HealthAdminServiceUnavailable(t, ctx, svc, ctrl)
		})
	})
})
