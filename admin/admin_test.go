package admin_test

import (
	"context"
	"errors"

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
		parent = zenkit.NewService("test-service")
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

})
