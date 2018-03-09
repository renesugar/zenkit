package zenkit

import (
	"github.com/goadesign/goa"
	"github.com/goadesign/goa/middleware"
	gometrics "github.com/rcrowley/go-metrics"
	"github.com/zenoss/zenkit/admin"
	"github.com/zenoss/zenkit/admin/app"
	"github.com/zenoss/zenkit/logging"
	"github.com/zenoss/zenkit/metrics"
)

func NewService(name string) *goa.Service {

	svc := goa.New(name)
	svc.WithLogger(logging.ServiceLogger())
	svc.Context = metrics.WithMetrics(svc.Context, gometrics.NewRegistry())
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
	// Assuming NewService was called before NewAdminService, this will link the
	// admin metrics controller with the metrics collected by parent service so that
	// any metrics collected by the parent service are reported by the AdminService.
	svc.Context = metrics.WithMetrics(svc.Context, metrics.ContextMetrics(parent.Context))

	c := admin.NewAdminController(svc)
	app.MountAdminController(svc, c)

	c2 := admin.NewHealthController(svc)
	app.MountHealthController(svc, c2)

	c3 := admin.NewSwaggerController(svc)
	app.MountSwaggerController(svc, c3)

	return svc
}
