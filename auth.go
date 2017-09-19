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
	"github.com/zenoss/zenkit/claims"
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

func JWTMiddleware(logger ErrorLogger, filename string, validator goa.Middleware, security *goa.JWTSecurity) (goa.Middleware, error) {
	key, err := ReadKeyFromFS(logger, filename)
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
			stdClaims := claims.StandardClaimsMap(token.Claims.(jwtgo.MapClaims))
			ident := &tokenIdentity{stdClaims}
			ctx = WithIdentity(ctx, ident)
			ctx = goa.WithLogContext(ctx, "user_id", ident.ID())
			if len(ident.Tenant()) == 0 {
				message := "Unable to retrieve tenant from token"
				logger := ContextLogger(ctx)
				if logger != nil {
					logger.Debug(message)
				}
				return errors.WithStack(errors.New(message))
			}
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
	Tenant() string
}

type tokenIdentity struct {
	claims.Claims
}

func (t *tokenIdentity) ID() string {
	return t.Claims.Subject()
}

func (t *tokenIdentity) Tenant() string {
	aud := t.Claims.Audience()

	// The tenant should be the only thing in the audience claim.
	if len(aud) != 1 {
		return ""
	}

	return aud[0]
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

func ReadKeyFromFS(logger ErrorLogger, filename string) ([]byte, error) {
	// Get the secret key
	var key []byte
	readKey := func() error {
		data, err := afero.ReadFile(FS, filename)
		if err != nil {
			logger.LogError("Unable to load auth key. Retrying.", "keyfile", filename, "err", err)
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
	/*
		This token is signed with the secret "secret" and gives the bearer the scopes "api:admin api:access"

		HEADER:

		{
		  "alg": "HS256",
		  "typ": "JWT"
		}

		PAYLOAD:
		{
			"iss": "Auth0",
			"sub": "1",
			"aud": [
				"anyone"
			],
			"exp": 1605511786,
			"nbf": 1505425386,
			"iat": 1505425386,
			"jti": "1",
			"scopes": "api:admin api:access"
		}
	*/
	devJWT = `Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJBdXRoMCIsInN1YiI6IjEiLCJhdWQiOlsiYW55b25lIl0sImV4cCI6MTYwNTUxMTc4NiwibmJmIjoxNTA1NDI1Mzg2LCJpYXQiOjE1MDU0MjUzODYsImp0aSI6IjEiLCJzY29wZXMiOiJhcGk6YWRtaW4gYXBpOmFjY2VzcyJ9.cmw2-w7efRhtNMDXypaI84UYHFZ_CmG0m9Jb6SnzX1Q`
)
