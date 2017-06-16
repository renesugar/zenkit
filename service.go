package zenkit

import (
	"github.com/goadesign/goa"
	"github.com/goadesign/goa/middleware"
)

func NewService(name string) *goa.Service {

	svc := goa.New(name)

	svc.WithLogger(ServiceLogger())

	svc.Use(middleware.RequestID())
	svc.Use(middleware.LogRequest(true))
	svc.Use(MetricsMiddleware())
	svc.Use(middleware.ErrorHandler(svc, true))
	svc.Use(middleware.Recover())

	return svc
}
