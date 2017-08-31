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
	// AuthZeroValidation is middleware that should be used by services receiving a jwt from Auth0
	AuthZeroValidation = JWTValidatorFunc(authZeroJWTValidator)
	// EdgeValidation is middleware that should be used by services receiving a jwt from edge that does not contain roles
	EdgeValidation = JWTValidatorFunc(edgeJWTValidator)
	// AuthorizationValidation is middleware that should be used by services receiving a jwt from authorization
	AuthorizationValidation = JWTValidatorFunc(authorizationJWTValidator)
	// CompleteValidation is middleware that should be used by services receiving a jwt from edge that does contain roles
	CompleteValidation = JWTValidatorFunc(completeJWTValidator)
	KeyFileTimeout     = 30 * time.Second
	localJWT           *design.SecuritySchemeDefinition
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
			ident := &tokenIdentity{claims: stdClaims}
			ctx = WithIdentity(ctx, ident)
			ctx = goa.WithLogContext(ctx, "user_id", ident.ID())
			return h(ctx, rw, req)
		}
	}
}

func authZeroJWTValidator(ctx context.Context) error {
	audience := ctx.Value(serviceKey).(string)
	m := claims.AuthZeroClaimsMap(jwt.ContextJWT(ctx).Claims.(jwtgo.MapClaims))
	err := m.Validate(audience)
	return err
}

func edgeJWTValidator(ctx context.Context) error {
	audience := ctx.Value(serviceKey).(string)
	m := claims.EdgeClaimsMap(jwt.ContextJWT(ctx).Claims.(jwtgo.MapClaims))
	err := m.Validate(audience)
	return err
}

func authorizationJWTValidator(ctx context.Context) error {
	audience := ctx.Value(serviceKey).(string)
	m := claims.AuthorizationClaimsMap(jwt.ContextJWT(ctx).Claims.(jwtgo.MapClaims))
	err := m.Validate(audience)
	return err
}

func completeJWTValidator(ctx context.Context) error {
	audience := ctx.Value(serviceKey).(string)
	m := claims.CompleteClaimsMap(jwt.ContextJWT(ctx).Claims.(jwtgo.MapClaims))
	err := m.Validate(audience)
	return err
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
	claims claims.Claims
}

func (t *tokenIdentity) ID() string {
	return t.claims.Subject()
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

// WithService registers the service name to the context
func WithService(ctx context.Context, service string) context.Context {
	return context.WithValue(ctx, serviceKey, service)
}

// ContextService gets the service name from the context
func ContextService(ctx context.Context) string {
	if v := ctx.Value(serviceKey); v != nil {
		return v.(string)
	}
	return ""
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
	// This token is signed with the secret "secret" and gives the bearer the scopes "api:admin api:access"
	devJWT = `Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJkZXZlbG9wZXIiLCJuYW1lIjoiQS4gRGV2ZWxvcGVyLCBFc3EuIiwiYWRtaW4iOnRydWUsInNjb3BlcyI6ImFwaTphZG1pbiBhcGk6YWNjZXNzIn0.e2YYHulpdvpnBdvdpUJyyJnC2xsm4VMrs6riy9WX4Ug`
)
