package zenkit

import (
	"context"
	"net/http"
	"time"

	"github.com/cenkalti/backoff"
	jwtgo "github.com/dgrijalva/jwt-go"
	"github.com/goadesign/goa"
	"github.com/goadesign/goa/design/apidsl"
	"github.com/goadesign/goa/middleware/security/jwt"
	"github.com/spf13/afero"
)

type JWTValidator func(ctx context.Context) error

var (
	FS                   = afero.NewReadOnlyFs(afero.NewOsFs())
	JWT                  = apidsl.JWTSecurity("jwt", func() { apidsl.Header("Authorization") })
	DefaultJWTValidation = JWTValidatorFunc(func(_ context.Context) error { return nil })
)

func JWTMiddleware(service *goa.Service, filename string, validator goa.Middleware, security *goa.JWTSecurity) (goa.Middleware, error) {
	key, err := readKeyFromFS(service, filename)
	if err != nil {
		return nil, err
	}
	resolver := jwt.NewSimpleResolver([]jwt.Key{key})
	return jwt.New(resolver, validator, security), nil
}

func JWTValidatorFunc(m JWTValidator) goa.Middleware {
	return func(h goa.Handler) goa.Handler {
		return func(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
			if err := m(ctx); err != nil {
				return err
			}
			token := jwt.ContextJWT(ctx)
			ident := &tokenIdentity{claims: token.Claims.(jwtgo.MapClaims)}
			ctx = WithIdentity(ctx, ident)
			ctx = goa.WithLogContext(ctx, "user_id", ident.ID())
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

func readKeyFromFS(service *goa.Service, filename string) ([]byte, error) {
	// Get the secret key
	var key []byte
	readKey := func() error {
		data, err := afero.ReadFile(FS, filename)
		if err != nil {
			service.LogError("Unable to load auth key. Retrying.", "keyfile", filename, "err", err)
			return err
		}
		key = data
		return nil
	}
	// Docker sometimes doesn't mount the secret right away, so we'll do a short retry
	boff := backoff.NewExponentialBackOff()
	boff.MaxElapsedTime = 30 * time.Second
	if err := backoff.Retry(readKey, boff); err != nil {
		return nil, err
	}
	return key, nil
}

const (
	ScopeAPIAccess = "api:access"
	ScopeAPIAdmin  = "api:admin"
)
