package zenkit

import (
	"github.com/goadesign/goa"
	"github.com/goadesign/goa/middleware"
)

func NewService(name string, devMode bool) *goa.Service {

	svc := goa.New(name)
	svc.WithLogger(ServiceLogger())

	if devMode {
		svc.Use(DevModeMiddleware)
	}
	svc.Use(middleware.RequestID())
	svc.Use(middleware.LogRequest(false))
	svc.Use(MetricsMiddleware())
	svc.Use(middleware.ErrorHandler(svc, true))
	svc.Use(middleware.Recover())

	return svc
}
