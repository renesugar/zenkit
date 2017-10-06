package admin_test

import (
	"context"

	"github.com/goadesign/goa"
	"github.com/zenoss/zenkit"
	. "github.com/zenoss/zenkit/admin"
	"github.com/zenoss/zenkit/admin/app/test"

	. "github.com/onsi/ginkgo"
	//. "github.com/onsi/gomega"
)

var _ = Describe("Swagger", func() {

	var (
		t      = GinkgoT()
		ctx    context.Context
		parent *goa.Service
		svc    = goa.New("swagger-test")
		ctrl   = NewSwaggerController(svc)
	)

	BeforeEach(func() {
		ctx = context.Background()
		parent = zenkit.NewService("test-service", false)
	})

	JustBeforeEach(func() {
		ctx = WithParentService(ctx, parent)
	})

	Context("when the Swagger resource is requested", func() {
		It("should respond OK", func() {
			test.SwaggerSwaggerOK(t, ctx, svc, ctrl)
		})
	})

	Context("when the JSON resource is requested", func() {
		originalAsset := SwaggerJSONAsset
		Context("when the swagger.json asset is missing", func() {
			BeforeEach(func() {
				SwaggerJSONAsset = "none"
			})
			AfterEach(func() {
				SwaggerJSONAsset = originalAsset
			})
			It("should respond with an InternalServerError", func() {
				test.JSONSwaggerInternalServerError(t, ctx, svc, ctrl)
			})
		})
		Context("when the swagger.json asset is available", func() {
			It("should respond OK", func() {
				test.JSONSwaggerOK(t, ctx, svc, ctrl)
			})
		})
	})
})
