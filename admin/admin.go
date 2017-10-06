package admin

import (
	"encoding/json"
	"github.com/goadesign/goa"
	"github.com/zenoss/zenkit/admin/app"
)

// AdminController implements the admin resource.
type AdminController struct {
	*goa.Controller
}

// NewAdminController creates a admin controller.
func NewAdminController(service *goa.Service) *AdminController {
	return &AdminController{Controller: service.NewController("AdminController")}
}

// Metrics runs the metrics action.
func (c *AdminController) Metrics(ctx *app.MetricsAdminContext) error {
	// AdminController_Metrics: start_implement

	registry := ContextMetrics(ctx)
	if registry == nil {
		// No registry was registered; must not be using metrics middleware.
		return ctx.OK([]byte("{}"))
	}
	encoder := json.NewEncoder(ctx.ResponseData)
	if ctx.Pretty {
		encoder.SetIndent("", "    ")
	}
	if err := encoder.Encode(registry); err != nil {
		return ctx.InternalServerError(err)
	}

	// AdminController_Metrics: end_implement
	return nil
}

// Ping runs the ping action.
func (c *AdminController) Ping(ctx *app.PingAdminContext) error {
	// AdminController_Ping: start_implement

	return ctx.OK([]byte(`PONG`))

	// AdminController_Ping: end_implement
	return nil
}
