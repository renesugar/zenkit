package zenkit

import (
	"github.com/goadesign/goa"
	"github.com/goadesign/goa/middleware"
	"github.com/goadesign/goa/middleware/xray"
	"github.com/pkg/errors"
)

var (
	ErrNoXRayDaemon = errors.New("no X-Ray daemon address")
)

func UseXRayMiddleware(service *goa.Service, address string, sampleRate int) error {
	if address == "" {
		return errors.WithStack(ErrNoXRayDaemon)
	}
	xraymw, err := xray.New(service.Name, address)
	if err != nil {
		service.LogError("Unable to initialize X-Ray middleware. Retrying.", "err", err)
		return errors.Wrap(err, "unable to initialize xray middleware")
	}
	service.Use(middleware.Tracer(sampleRate, xray.NewID, xray.NewTraceID))
	service.Use(xraymw)
	return nil
}
