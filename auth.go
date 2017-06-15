package zenkit

import (
	"context"
	"net/http"

	jwtgo "github.com/dgrijalva/jwt-go"
	"github.com/goadesign/goa"
	"github.com/goadesign/goa/design/apidsl"
	"github.com/goadesign/goa/middleware/security/jwt"
)

type JWTValidator func(ctx context.Context) error

var (
	JWT                  = apidsl.JWTSecurity("jwt", func() { apidsl.Header("Authorization") })
	DefaultJWTValidation = JWTValidatorFunc(func(_ context.Context) error { return nil })
)

func JWTMiddleware(key []byte, validator goa.Middleware, security *goa.JWTSecurity) goa.Middleware {
	resolver := jwt.NewSimpleResolver([]jwt.Key{key})
	return jwt.New(resolver, validator, security)
}

func JWTValidatorFunc(m JWTValidator) goa.Middleware {
	return func(h goa.Handler) goa.Handler {
		return func(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
			if err := m(ctx); err != nil {
				return err
			}
			token := jwt.ContextJWT(ctx)
			ctx = WithIdentity(ctx, &tokenIdentity{claims: token.Claims.(jwtgo.MapClaims)})
			return h(ctx, rw, req)
		}
	}
}

type Identity interface {
	ID() string
}

type tokenIdentity struct {
	claims jwtgo.MapClaims
}

func (t *tokenIdentity) ID() string {
	return t.claims["sub"].(string)
}

func WithIdentity(ctx context.Context, identity Identity) context.Context {
	return context.WithValue(ctx, identityKey, identity)
}

func ContextIdentity(ctx context.Context) Identity {
	if v := ctx.Value(identityKey); v != nil {
		return v.(Identity)
	}
	return nil
}

const (
	ScopeAPIAccess = "api:access"
	ScopeAPIAdmin  = "api:admin"
)
