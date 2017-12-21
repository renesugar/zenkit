package zenkit

import (
	"github.com/goadesign/goa"
	"github.com/goadesign/goa/middleware"
	"github.com/zenoss/zenkit/admin"
	"github.com/zenoss/zenkit/admin/app"
	"github.com/zenoss/zenkit/auth"
	"github.com/zenoss/zenkit/logging"
	"github.com/zenoss/zenkit/metrics"
)

func NewService(name string, authDisabled bool) *goa.Service {

	svc := goa.New(name)
	svc.WithLogger(logging.ServiceLogger())

	if authDisabled {
		svc.Use(auth.DevModeMiddleware)
	}
	svc.Use(middleware.RequestID())
	svc.Use(middleware.LogRequest(false))
	svc.Use(metrics.MetricsMiddleware())
	svc.Use(middleware.ErrorHandler(svc, true))
	svc.Use(logging.LogErrorResponse())
	svc.Use(middleware.Recover())

	return svc
}

func NewAdminService(parent *goa.Service) *goa.Service {

	svc := goa.New("admin")
	svc.Context = admin.WithParentService(svc.Context, parent)

	c := admin.NewAdminController(svc)
	app.MountAdminController(svc, c)

	c2 := admin.NewHealthController(svc)
	app.MountHealthController(svc, c2)

	c3 := admin.NewSwaggerController(svc)
	app.MountSwaggerController(svc, c3)

	return svc
}
