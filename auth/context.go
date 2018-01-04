package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	jwtgo "github.com/dgrijalva/jwt-go"
	"github.com/goadesign/goa"
	"github.com/goadesign/goa/middleware/security/jwt"
	"github.com/pkg/errors"
	"github.com/zenoss/zenkit/claims"
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

func JWTMiddleware(h goa.Handler) goa.Handler {
	return func(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
		val := req.Header.Get("Authorization")
		if val == "" {
			return jwt.ErrJWTError("missing header \"Authorization\"")
		}

		parts := strings.Split(val, ".")
		if len(parts) != 3 {
			return jwt.ErrJWTError("JWT validation failed")
		}

		// parse Claims
		var mapClaims jwtgo.MapClaims
		claimBytes, _ := jwtgo.DecodeSegment(parts[1])
		err := json.Unmarshal(claimBytes, &mapClaims)
		if err != nil {
			return jwt.ErrJWTError("JWT validation failed")
		}

		stdClaims := claims.StandardClaimsMap(mapClaims)
		ident := &tokenIdentity{stdClaims}
		ctx = WithIdentity(ctx, ident)
		ctx = goa.WithLogContext(ctx, "user_id", ident.ID())
		if len(ident.Tenant()) == 0 {
			return errors.New("unable to retrieve tenant from token")
		}

		return h(ctx, rw, req)
	}
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
