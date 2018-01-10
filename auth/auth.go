package auth

import (
	"context"
	"fmt"
	"net/http"
	"time"

	jwtgo "github.com/dgrijalva/jwt-go"
	"github.com/goadesign/goa"
	"github.com/goadesign/goa/design"
	"github.com/goadesign/goa/design/apidsl"
	"github.com/pkg/errors"
	"github.com/spf13/afero"
)

var (
	FS                  = afero.NewReadOnlyFs(afero.NewOsFs())
	AuthorizationHeader = "Authorization"
	KeyFileTimeout      = 30 * time.Second
	localJWT            *design.SecuritySchemeDefinition
	devToken            string
)

func JWT() *design.SecuritySchemeDefinition {
	if localJWT == nil {
		localJWT = apidsl.JWTSecurity("jwt", func() {
			apidsl.Header(AuthorizationHeader)
		})
	}
	return localJWT
}

// BuildToken builds a token from the given params
func BuildToken(claims jwtgo.Claims, signingMethod jwtgo.SigningMethod, key interface{}) (string, error) {
	token := jwtgo.NewWithClaims(signingMethod, claims)
	signedToken, err := token.SignedString(key)
	if err != nil {
		return "", err
	}
	return signedToken, nil
}

// NewDevJWTMiddleware creates a middleware that inserts a dev token in the request header
func NewDevJWTMiddleware(devClaims jwtgo.Claims, signingMethod jwtgo.SigningMethod, key interface{}) goa.Middleware {
	return func(h goa.Handler) goa.Handler {
		return func(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
			header := req.Header.Get(AuthorizationHeader)
			if header == "" {
				if len(devToken) == 0 {
					signedToken, err := BuildToken(devClaims, signingMethod, key)
					if err != nil {
						return errors.Wrap(err, "Unable to sign token")
					}
					devToken = fmt.Sprintf("Bearer %s", signedToken)
				}
				req.Header.Set("Authorization", devToken)
			}
			return h(ctx, rw, req)
		}
	}
}
