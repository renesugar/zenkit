package auth

import (
	"context"
	"strings"

	"github.com/zenoss/zenkit/claims"
)

type key int

const (
	identityKey key = iota + 1
)

func WithTenantIdentity(ctx context.Context, identity TenantIdentity) context.Context {
	return context.WithValue(ctx, identityKey, identity)
}

func ContextTenantIdentity(ctx context.Context) TenantIdentity {
	if v := ctx.Value(identityKey); v != nil {
		return v.(TenantIdentity)
	}
	return nil
}

// TenantIdentity is an identity in a multi-tenant application
type TenantIdentity interface {
	ID() string
	Tenant() string
}

// Auth0TenantIdentity is an identity from Auth0 for a multi-tenant application
type Auth0TenantIdentity struct {
	id     string
	tenant string
}

// NewAuth0TenantIdentity creates an Auth0TenantIdentity for the tokenClaims
func NewAuth0TenantIdentity(tokenClaims claims.MultiTenantClaims) *Auth0TenantIdentity {
	// subject comes in from Auth0 as "{idp}|{connection}|{userid}"
	// eg: "sub": "ad|acmeco|thedude",
	// connection may be missing in the case that the idp is the connection
	// eg: "sub": "google-apps|foo@bar.com"
	identParts := strings.Split(tokenClaims.Subject(), "|")
	return &Auth0TenantIdentity{
		id:     identParts[len(identParts)-1],
		tenant: tokenClaims.Tenant(),
	}
}

// ID gets the user id for the identity
func (ti *Auth0TenantIdentity) ID() string {
	return ti.id
}

// Tenant gets the tenant the identity belogs to
func (ti *Auth0TenantIdentity) Tenant() string {
	return ti.tenant
}
