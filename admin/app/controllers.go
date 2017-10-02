// Code generated by goagen v1.2.0-dirty, DO NOT EDIT.
//
// API "Admin": Application Controllers
//
// Command:
// $ goagen
// --design=github.com/zenoss/zenkit/admin/design
// --out=$(GOPATH)/src/github.com/zenoss/zenkit/admin
// --version=v1.3.0

package app

import (
	"context"
	"github.com/goadesign/goa"
	"net/http"
)

// initService sets up the service encoders, decoders and mux.
func initService(service *goa.Service) {
	// Setup encoders and decoders
	service.Encoder.Register(goa.NewJSONEncoder, "application/json")
	service.Decoder.Register(goa.NewJSONDecoder, "application/json")

	// Setup default encoder and decoder
	service.Encoder.Register(goa.NewJSONEncoder, "*/*")
	service.Decoder.Register(goa.NewJSONDecoder, "*/*")
}

// AdminController is the controller interface for the Admin actions.
type AdminController interface {
	goa.Muxer
	Metrics(*MetricsAdminContext) error
	Ping(*PingAdminContext) error
	Swagger(*SwaggerAdminContext) error
	SwaggerJSON(*SwaggerJSONAdminContext) error
}

// MountAdminController "mounts" a Admin resource controller on the given service.
func MountAdminController(service *goa.Service, ctrl AdminController) {
	initService(service)
	var h goa.Handler

	h = func(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
		// Check if there was an error loading the request
		if err := goa.ContextError(ctx); err != nil {
			return err
		}
		// Build the context
		rctx, err := NewMetricsAdminContext(ctx, req, service)
		if err != nil {
			return err
		}
		return ctrl.Metrics(rctx)
	}
	service.Mux.Handle("GET", "/metrics", ctrl.MuxHandler("metrics", h, nil))
	service.LogInfo("mount", "ctrl", "Admin", "action", "Metrics", "route", "GET /metrics")

	h = func(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
		// Check if there was an error loading the request
		if err := goa.ContextError(ctx); err != nil {
			return err
		}
		// Build the context
		rctx, err := NewPingAdminContext(ctx, req, service)
		if err != nil {
			return err
		}
		return ctrl.Ping(rctx)
	}
	service.Mux.Handle("HEAD", "/ping", ctrl.MuxHandler("ping", h, nil))
	service.LogInfo("mount", "ctrl", "Admin", "action", "Ping", "route", "HEAD /ping")
	service.Mux.Handle("GET", "/ping", ctrl.MuxHandler("ping", h, nil))
	service.LogInfo("mount", "ctrl", "Admin", "action", "Ping", "route", "GET /ping")

	h = func(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
		// Check if there was an error loading the request
		if err := goa.ContextError(ctx); err != nil {
			return err
		}
		// Build the context
		rctx, err := NewSwaggerAdminContext(ctx, req, service)
		if err != nil {
			return err
		}
		return ctrl.Swagger(rctx)
	}
	service.Mux.Handle("GET", "/swagger", ctrl.MuxHandler("swagger", h, nil))
	service.LogInfo("mount", "ctrl", "Admin", "action", "Swagger", "route", "GET /swagger")

	h = func(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
		// Check if there was an error loading the request
		if err := goa.ContextError(ctx); err != nil {
			return err
		}
		// Build the context
		rctx, err := NewSwaggerJSONAdminContext(ctx, req, service)
		if err != nil {
			return err
		}
		return ctrl.SwaggerJSON(rctx)
	}
	service.Mux.Handle("GET", "/swagger.json", ctrl.MuxHandler("swagger.json", h, nil))
	service.LogInfo("mount", "ctrl", "Admin", "action", "SwaggerJSON", "route", "GET /swagger.json")
}
