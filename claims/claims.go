package claims

import (
	"errors"
	"time"
)

var (
	// ErrIssuer occurs when the issuer claimed is not valid
	ErrIssuer = errors.New("issuer is invalid")
	// ErrSubject occurs when the subject claimed is not valid
	ErrSubject = errors.New("subject is invalid")
	// ErrAudience occurs when the audience claimed is not valid
	ErrAudience = errors.New("audience is invalid")
	// ErrExpiresAt occurs when the exp claimed is not valid
	ErrExpiresAt = errors.New("token is expired")
	// ErrNotBefore occurs when the nbf claimed is not valid
	ErrNotBefore = errors.New("used before eligible")
	// ErrIssuedAt occurs when the iat claimed is not valid
	ErrIssuedAt = errors.New("used before issued")
	// ErrID occurs when the issuer jti is not valid
	ErrID = errors.New("ID is invalid")
)

// Claims exposes a collection of mandatory claims from a JWT
type Claims interface {
	// Issuer is the jwt "iss" claim
	Issuer() string
	// Subject is the jwt "sub" claim
	Subject() string
	// Audience is the jwt "aud" claim
	Audience() []string
	// ExpiresAt is the jwt "exp" claim
	ExpiresAt() int64
	// NotBefore is the jwt "nbf" claim
	NotBefore() int64
	// IssuedAt is the jwt "iat" claim
	IssuedAt() int64
	// ID is the jwt "jti" claim
	ID() string
}

// Valid determines if a JWT should be rejected or not based on
// standard, mandatory jwt claims
func Valid(claims Claims) error {
	now := time.Now().Unix()
	if !verifyIssuerExists(claims.Issuer()) {
		return ErrIssuer
	} else if !verifySubject(claims.Subject()) {
		return ErrSubject
	} else if !verifyAudienceExists(claims.Audience()) {
		return ErrAudience
	} else if !verifyExpiresAt(claims.ExpiresAt(), now) {
		return ErrExpiresAt
	} else if !verifyNotBefore(claims.NotBefore(), now) {
		return ErrNotBefore
	} else if !verifyIssuedAt(claims.IssuedAt(), now) {
		return ErrIssuedAt
	} else if !verifyID(claims.ID()) {
		return ErrID
	}
	return nil
}

// ValidateIssuer verifies that the issuer in claims is in issuers, a slice of
// valid issuers
func ValidateIssuer(claims Claims, issuers []string) error {
	if !verifyIssuer(claims.Issuer(), issuers) {
		return ErrIssuer
	}
	return nil
}

// ValidateAudience verifies that the audience in claims contains audience
func ValidateAudience(claims Claims, audience string) error {
	if !verifyAudience(claims.Audience(), audience) {
		return ErrAudience
	}
	return nil
}

func verifyIssuerExists(claimed string) bool {
	return claimed != ""
}

func verifyIssuer(claimed string, issuers []string) bool {
	for _, iss := range issuers {
		if iss == claimed {
			return true
		}
	}
	return false
}

func verifySubject(claimed string) bool {
	return claimed != ""
}

func verifyAudienceExists(claimed []string) bool {
	return claimed != nil && len(claimed) > 0
}

func verifyAudience(claimed []string, audience string) bool {
	for _, aud := range claimed {
		if aud == audience {
			return true
		}
	}
	return false
}

func verifyExpiresAt(claimed int64, now int64) bool {
	if claimed != 0 && claimed > now {
		return true
	}
	return false
}

func verifyNotBefore(claimed int64, now int64) bool {
	if claimed != 0 && claimed <= now {
		return true
	}
	return false
}

func verifyIssuedAt(claimed int64, now int64) bool {
	if claimed != 0 && claimed <= now {
		return true
	}
	return false
}

func verifyID(claimed string) bool {
	return claimed != ""
}
