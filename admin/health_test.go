package admin_test

import (
	"context"
	"errors"

	"github.com/goadesign/goa"
	"github.com/zenoss/zenkit"
	. "github.com/zenoss/zenkit/admin"
	"github.com/zenoss/zenkit/admin/app"
	"github.com/zenoss/zenkit/admin/app/test"
	"github.com/zenoss/zenkit/healthcheck"

	. "github.com/onsi/ginkgo"
	//. "github.com/onsi/gomega"
)

var _ = Describe("Health", func() {
	var (
		t      = GinkgoT()
		ctx    context.Context
		parent *goa.Service
		svc    = goa.New("admin-test")
		ctrl   = NewHealthController(svc)
	)

	BeforeEach(func() {
		ctx = context.Background()
		parent = zenkit.NewService("test-service")
	})

	JustBeforeEach(func() {
		ctx = WithParentService(ctx, parent)
	})
	AfterEach(func() {
		ResetRegistry()
	})

	Context("when the Health resource is requested", func() {
		It("should return OK if there are no failing health checks", func() {
			check := func() error { return nil }
			healthcheck.RegisterFunc("testOK", check)
			test.HealthHealthOK(t, ctx, svc, ctrl)
		})

		It("should return ServiceUnavailable if there are failing health checks", func() {
			check := func() error { return errors.New("he dead") }
			healthcheck.RegisterFunc("testDOWN", check)
			test.HealthHealthServiceUnavailable(t, ctx, svc, ctrl)
		})
	})

	It("should change the response of the healthcheck", func() {

		By("applying the DOWN state to the service")

		test.DownHealthOK(t, ctx, svc, ctrl, &app.DownHealthPayload{
			Reason: "testing",
		})
		test.HealthHealthServiceUnavailable(t, ctx, svc, ctrl)

		By("applying the UP state to the service")

		test.UpHealthOK(t, ctx, svc, ctrl)
		test.HealthHealthOK(t, ctx, svc, ctrl)
	})
})
