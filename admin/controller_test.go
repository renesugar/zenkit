package admin_test

import (
	"github.com/goadesign/goa"
	. "github.com/zenoss/zenkit/admin"

	. "github.com/onsi/ginkgo"
	//	. "github.com/onsi/gomega"
)

var _ = Describe("Controller", func() {

	var (
		svc = goa.New("controller-test")
	)

	It("should mount all controllers", func() {
		MountAllControllers(svc)
	})
})
