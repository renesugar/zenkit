package design

import (
	. "github.com/goadesign/goa/design/apidsl"
)

var _ = API("Admin", func() {
	Title("Admin Service")
	Description("Utilities provided by the admin endpoint")
	Scheme("http")
	Consumes("application/json")
	Produces("application/json")
})
