package zenkit

import (
	"github.com/goadesign/goa"
	"github.com/goadesign/goa/design/apidsl"
	"github.com/goadesign/goa/middleware"
)

var JWT = apidsl.JWTSecurity("jwt", func() {
	apidsl.Header("Authorization")
})

func NewService(name string, key []byte) *goa.Service {

	svc := goa.New(name)
	svc.Use(middleware.RequestID())
	svc.Use(middleware.LogRequest(false))
	svc.Use(MetricsMiddleware())
	svc.Use(middleware.ErrorHandler(svc, true))
	svc.Use(middleware.Recover())

	svc.WithLogger(ServiceLogger())

	//// Set up security
	//resolver := jwt.NewSimpleResolver([]jwt.Key{key})
	//validationFunc := func(next goa.Handler) goa.Handler {
	//return next
	//}
	//svc.Use(jwt.New(resolver, validationFunc, JWT))

	return svc
}
