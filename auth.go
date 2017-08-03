package zenkit

import (
	"context"
	"net/http"
	"time"

	"github.com/cenkalti/backoff"
	jwtgo "github.com/dgrijalva/jwt-go"
	"github.com/goadesign/goa"
	"github.com/goadesign/goa/design"
	"github.com/goadesign/goa/design/apidsl"
	"github.com/goadesign/goa/middleware/security/jwt"
	"github.com/pkg/errors"
	"github.com/spf13/afero"
)

type JWTValidator func(ctx context.Context) error

var (
	FS                   = afero.NewReadOnlyFs(afero.NewOsFs())
	AuthorizationHeader  = "Authorization"
	DefaultJWTValidation = JWTValidatorFunc(func(_ context.Context) error { return nil })
	KeyFileTimeout       = 30 * time.Second
	localJWT             *design.SecuritySchemeDefinition
)

func JWT() *design.SecuritySchemeDefinition {
	if localJWT == nil {
		localJWT = apidsl.JWTSecurity("jwt", func() {
			apidsl.Header(AuthorizationHeader)
		})
	}
	return localJWT
}

func JWTMiddleware(service *goa.Service, filename string, validator goa.Middleware, security *goa.JWTSecurity) (goa.Middleware, error) {
	key, err := readKeyFromFS(service, filename)
	if err != nil {
		return nil, errors.Wrap(err, "couldn't read key from filesystem")
	}
	resolver := jwt.NewSimpleResolver([]jwt.Key{key})
	return jwt.New(resolver, validator, security), nil
}

func JWTValidatorFunc(m JWTValidator) goa.Middleware {
	return func(h goa.Handler) goa.Handler {
		return func(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
			if err := m(ctx); err != nil {
				return errors.WithStack(err)
			}
			token := jwt.ContextJWT(ctx)
			ident := &tokenIdentity{claims: token.Claims.(jwtgo.MapClaims)}
			ctx = WithIdentity(ctx, ident)
			ctx = goa.WithLogContext(ctx, "user_id", ident.ID())
			return h(ctx, rw, req)
		}
	}
}

func DevModeMiddleware(h goa.Handler) goa.Handler {
	return func(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
		header := req.Header.Get("Authorization")
		if header == "" {
			req.Header.Set("Authorization", devJWT)
		}
		return h(ctx, rw, req)
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
			return errors.Wrap(err, "unable to load auth key")
		}
		key = data
		return nil
	}
	// Docker sometimes doesn't mount the secret right away, so we'll do a short retry
	boff := backoff.NewExponentialBackOff()
	boff.MaxElapsedTime = KeyFileTimeout
	if err := backoff.Retry(readKey, boff); err != nil {
		return nil, errors.Wrap(err, "unable to load auth key within the timeout")
	}
	return key, nil
}

const (
	ScopeAPIAccess = "api:access"
	ScopeAPIAdmin  = "api:admin"
)

const (
	// This token is signed with the secret "secret" and gives the bearer the scopes "api:admin api:access"
	devJWT = `Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJkZXZlbG9wZXIiLCJuYW1lIjoiQS4gRGV2ZWxvcGVyLCBFc3EuIiwiYWRtaW4iOnRydWUsInNjb3BlcyI6ImFwaTphZG1pbiBhcGk6YWNjZXNzIn0.e2YYHulpdvpnBdvdpUJyyJnC2xsm4VMrs6riy9WX4Ug`
)
