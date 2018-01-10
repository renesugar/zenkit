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
	// ErrTenant occurs when a MultiTenantClaims does not have a tenant
	ErrTenant = errors.New("Tenant is empty")
)

// Claims exposes a collection of mandatory claims from a JWT
type Claims interface {
	Issuer() string
	Subject() string
	Audience() []string
	ExpiresAt() int64
	IssuedAt() int64
	Valid() error
}

// MultiTenantClaims is a set of claims fit for consumption by a multitenant app
type MultiTenantClaims interface {
	Issuer() string
	Subject() string
	Audience() []string
	ExpiresAt() int64
	IssuedAt() int64
	Tenant() string
	Valid() error
}

// ValidateClaims determines if a JWT should be rejected
func ValidateClaims(claims Claims) error {
	now := time.Now().Unix()
	if !verifyIssuerExists(claims.Issuer()) {
		return ErrIssuer
	} else if !verifySubject(claims.Subject()) {
		return ErrSubject
	} else if !verifyAudienceExists(claims.Audience()) {
		return ErrAudience
	} else if !verifyExpiresAt(claims.ExpiresAt(), now) {
		return ErrExpiresAt
	} else if !verifyIssuedAt(claims.IssuedAt(), now) {
		return ErrIssuedAt
	}
	return nil
}

// ValidateMultiTenantClaims determines if a JWT should be rejected
func ValidateMultiTenantClaims(claims MultiTenantClaims) error {
	if err := ValidateClaims(claims); err != nil {
		return err
	} else if !verifyTenant(claims.Tenant()) {
		return ErrTenant
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

func verifyExpiresAt(claimed int64, now int64) bool {
	if claimed != 0 && claimed > now {
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

func verifyTenant(claimed string) bool {
	return claimed != ""
}
