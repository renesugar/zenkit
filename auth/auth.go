package auth

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/cenkalti/backoff"
	jwtgo "github.com/dgrijalva/jwt-go"
	"github.com/goadesign/goa"
	"github.com/goadesign/goa/client"
	"github.com/goadesign/goa/design"
	"github.com/goadesign/goa/design/apidsl"
	"github.com/goadesign/goa/middleware/security/jwt"
	"github.com/pkg/errors"
	"github.com/spf13/afero"
	"github.com/zenoss/zenkit/claims"
	"github.com/zenoss/zenkit/logging"
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

func JWTSigner(req *http.Request) *client.JWTSigner {
	token := &client.StaticToken{}

	parts := strings.Fields(req.Header.Get("Authorization"))
	switch len(parts) {
	case 0:
		return nil
	case 1:
		token.Value = parts[0]
	default:
		token.Type = parts[0]
		token.Value = strings.Join(parts[1:], " ")
	}
	return &client.JWTSigner{TokenSource: &client.StaticTokenSource{StaticToken: token}}
}

func BuildDevToken(ctx context.Context, signingMethod jwtgo.SigningMethod) string {
	token := jwtgo.NewWithClaims(signingMethod, devClaims)
	signedToken, _ := token.SignedString([]byte(signingKey))
	return signedToken
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

// GetKeysFromFS creates a slice of jwt.Key from keys in files
func GetKeysFromFS(logger logging.ErrorLogger, files []string) ([]jwt.Key, error) {
	var parsedKeys []jwt.Key
	for _, keyFile := range files {
		keyBytes, err := ReadKeyFromFS(logger, keyFile)
		if err != nil {
			return []jwt.Key{}, errors.Wrap(err, "Unable to read key file from fs")
		}
		key := ConvertToKey(keyBytes)
		parsedKeys = append(parsedKeys, key)
	}
	return parsedKeys, nil
}

// ConvertToKey converts a byte slice to a jwt.Key
func ConvertToKey(key []byte) jwt.Key {
	pubkey, err := jwtgo.ParseRSAPublicKeyFromPEM(key)
	if err == nil {
		return pubkey
	}

	return key
}

func ReadKeyFromFS(logger logging.ErrorLogger, filename string) ([]byte, error) {
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
	signingKey = "secret"
)
