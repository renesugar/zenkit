package admin

import (
	"fmt"

	"github.com/goadesign/goa"
	"github.com/zenoss/zenkit/admin/app"
	"github.com/zenoss/zenkit/admin/swagger"
)

// SwaggerController implements the swagger resource.
type SwaggerController struct {
	*goa.Controller
}

// NewSwaggerController creates a swagger controller.
func NewSwaggerController(service *goa.Service) *SwaggerController {
	return &SwaggerController{Controller: service.NewController("SwaggerController")}
}

// JSON runs the json action.
func (c *SwaggerController) JSON(ctx *app.JSONSwaggerContext) error {
	// SwaggerController_JSON: start_implement

	data, err := swagger.Asset(SwaggerJSONAsset)
	if err != nil {
		return ctx.InternalServerError(err)
	}
	return ctx.OK(data)

	// SwaggerController_JSON: end_implement
	return nil
}

// Swagger runs the swagger action.
func (c *SwaggerController) Swagger(ctx *app.SwaggerSwaggerContext) error {
	// SwaggerController_Swagger: start_implement

	s := ContextParentService(ctx)
	htmlText := fmt.Sprintf(`<!DOCTYPE html
<html>
       <head>
               <title>%s API</title>
               <meta name="viewport" content="width=device-width, initial-scale=1">
               <style>
                       body {
                               margin: 0;
                               padding: 0;
                       }
               </style>
       </head>
       <body>
               <redoc spec-url='/swagger.json'></redoc>
               <script src="https://rebilly.github.io/ReDoc/releases/latest/redoc.min.js"></script>
       </body>
</html>`, s.Name)

	return ctx.OK([]byte(htmlText))

	// SwaggerController_Swagger: end_implement
	return nil
}
