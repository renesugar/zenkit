package zenkit

import (
	"errors"

	"github.com/goadesign/goa"
	"github.com/goadesign/goa/middleware"
	"github.com/goadesign/goa/middleware/xray"
)

var (
	ErrNoXRayDaemon = errors.New("no X-Ray daemon address")
)

func UseXRayMiddleware(service *goa.Service, address string, sampleRate int) error {
	if address == "" {
		return ErrNoXRayDaemon
	}
	xraymw, err := xray.New(service.Name, address)
	if err != nil {
		service.LogError("Unable to initialize X-Ray middleware. Retrying.", "err", err)
		return err
	}
	service.Use(middleware.Tracer(sampleRate, xray.NewID, xray.NewTraceID))
	service.Use(xraymw)
	return nil
}
