package admin

import (
	"github.com/goadesign/goa"
	"github.com/zenoss/zenkit/admin/app"
)

// MountAllControllers mounts all of the controllers in this package
func MountAllControllers(service *goa.Service) {

	// Mount the Admin Controller
	c1 := NewAdminController(service)
	app.MountAdminController(service, c1)

	// Mount the Health Controller
	c2 := NewHealthController(service)
	app.MountHealthController(service, c2)
}
