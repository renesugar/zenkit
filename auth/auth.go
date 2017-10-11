package auth

import (
	"context"
	"time"

	jwtgo "github.com/dgrijalva/jwt-go"
	"github.com/goadesign/goa/design"
	"github.com/goadesign/goa/design/apidsl"
	"github.com/spf13/afero"
)

type JWTValidator func(ctx context.Context) error

var (
	FS                   = afero.NewReadOnlyFs(afero.NewOsFs())
	AuthorizationHeader  = "Authorization"
	DefaultJWTValidation = JWTValidatorFunc(func(_ context.Context) error { return nil })
	KeyFileTimeout       = 30 * time.Second
	localJWT             *design.SecuritySchemeDefinition
	devJWT               string

	// These claims are used to populate the dev token (devJWT) using the secret defined by signingKey
	//  - exp is equivalent to Monday, November 16, 2020 7:29:46 AM GMT
	//  - iat/nbf is equivalent to Thursday, September 14, 2017 9:43:06 PM GMT
	// These claims are a union of claims that could be provided by Auth0 or the edge service:
	//  - "https://zing.zenoss/tnt" and "https://zing.zenoss/src" are used to simulate tokens provided by Auth0
	//  - "aud" and "src" are used to simulate tokens from the edge service.
	devClaims = jwtgo.MapClaims{
		"iss": "Auth0",
		"sub": "1",
		"aud": []string{"anyone"},
		"https://zing.zenoss/tnt": "anyone",
		"exp":    1605511786,
		"nbf":    1505425386,
		"iat":    1505425386,
		"jti":    "1",
		"scopes": "api:admin api:access",
		"src":    "rm1",
		"https://zing.zenoss/src": "rm1",
	}
)

func JWT() *design.SecuritySchemeDefinition {
	if localJWT == nil {
		localJWT = apidsl.JWTSecurity("jwt", func() {
			apidsl.Header(AuthorizationHeader)
		})
	}
	return localJWT
}

func BuildDevToken(ctx context.Context, signingMethod jwtgo.SigningMethod) string {
	token := jwtgo.NewWithClaims(signingMethod, devClaims)
	signedToken, _ := token.SignedString([]byte(signingKey))
	return signedToken
}

const (
	ScopeAPIAccess = "api:access"
	ScopeAPIAdmin  = "api:admin"
)

const (
	signingKey = "secret"
)
