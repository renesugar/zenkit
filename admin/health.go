package admin

import (
	"errors"
	"github.com/goadesign/goa"
	"github.com/zenoss/zenkit/admin/app"
	"github.com/zenoss/zenkit/healthcheck"
)

// HealthController implements the health resource.
type HealthController struct {
	*goa.Controller
}

// NewHealthController creates a health controller.
func NewHealthController(service *goa.Service) *HealthController {
	return &HealthController{Controller: service.NewController("HealthController")}
}

// Down runs the down action.
func (c *HealthController) Down(ctx *app.DownHealthContext) error {
	// HealthController_Down: start_implement

	ContextLogger(ctx).WithField("reason", ctx.Payload.Reason).Info("Manual HTTP Status DOWN")
	updater.Update(errors.New(ctx.Payload.Reason))

	// HealthController_Down: end_implement
	return nil
}

// Health runs the health action.
func (c *HealthController) Health(ctx *app.HealthHealthContext) error {
	// HealthController_Health: start_implement

	output := healthcheck.CheckStatus()
	if len(output) > 0 {
		return ctx.ServiceUnavailable(output)
	}

	// HealthController_Health: end_implement
	return nil
}

// Up runs the up action.
func (c *HealthController) Up(ctx *app.UpHealthContext) error {
	// HealthController_Up: start_implement

	ContextLogger(ctx).Info("Manual HTTP Status UP")
	updater.Update(nil)

	// HealthController_Up: end_implement
	return nil
}
