package auth

import "github.com/zenoss/zenkit/claims"

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
