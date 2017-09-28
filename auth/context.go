package auth

import (
	"context"
	"fmt"
	"net/http"

	jwtgo "github.com/dgrijalva/jwt-go"
	"github.com/goadesign/goa"
	"github.com/goadesign/goa/middleware/security/jwt"
	"github.com/pkg/errors"
	"github.com/zenoss/zenkit/claims"
	"github.com/zenoss/zenkit/logging"
)

type key int

const (
	identityKey key = iota + 1
)

func WithIdentity(ctx context.Context, identity Identity) context.Context {
	return context.WithValue(ctx, identityKey, identity)
}

func ContextIdentity(ctx context.Context) Identity {
	if v := ctx.Value(identityKey); v != nil {
		return v.(Identity)
	}
	return nil
}

func JWTMiddleware(logger logging.ErrorLogger, filename string, validator goa.Middleware, security *goa.JWTSecurity) (goa.Middleware, error) {
	key, err := ReadKeyFromFS(logger, filename)
	if err != nil {
		return nil, errors.Wrap(err, "couldn't read key from filesystem")
	}
	resolver := jwt.NewSimpleResolver([]jwt.Key{key})
	return jwt.New(resolver, validator, security), nil
}

func DevModeMiddleware(h goa.Handler) goa.Handler {
	return func(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
		header := req.Header.Get("Authorization")
		if header == "" {
			if len(devJWT) == 0 {
				signedToken := BuildDevToken(ctx, jwtgo.SigningMethodHS256)
				devJWT = fmt.Sprintf("Bearer %s", signedToken)
			}
			req.Header.Set("Authorization", devJWT)
		}
		return h(ctx, rw, req)
	}
}

func JWTValidatorFunc(m JWTValidator) goa.Middleware {
	return func(h goa.Handler) goa.Handler {
		return func(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
			if err := m(ctx); err != nil {
				return errors.WithStack(err)
			}
			token := jwt.ContextJWT(ctx)
			stdClaims := claims.StandardClaimsMap(token.Claims.(jwtgo.MapClaims))
			ident := &tokenIdentity{stdClaims}
			ctx = WithIdentity(ctx, ident)
			ctx = goa.WithLogContext(ctx, "user_id", ident.ID())
			if len(ident.Tenant()) == 0 {
				message := "Unable to retrieve tenant from token"
				return errors.WithStack(errors.New(message))
			}
			return h(ctx, rw, req)
		}
	}
}
