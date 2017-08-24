package claims

import (
	"errors"
	"time"
)

var (
	ErrIssuer    = errors.New("issuer is invalid")
	ErrSubject   = errors.New("subject is invalid")
	ErrAudience  = errors.New("audience is invalid")
	ErrExpiresAt = errors.New("token is expired")
	ErrNotBefore = errors.New("used before eligible")
	ErrIssuedAt  = errors.New("used before issued")
	ErrID        = errors.New("ID is invalid")
)

// StringOrURI is a case-sensitive string or uri which can be parsed into a url
type StringOrURI string

// StandardClaims implements registered claim names according to RFC 7519
type StandardClaims struct {
	Issuer    StringOrURI   `json:"iss"`
	Subject   StringOrURI   `json:"sub"`
	Audience  []StringOrURI `json:"aud"`
	ExpiresAt int64         `json:"exp"`
	NotBefore int64         `json:"nbf"`
	IssuedAt  int64         `json:"iat"`
	ID        string        `json:"jti"`
}

// Valid determines if a JWT should be rejected or not and implements jwt-go Claims interface
func (claims *StandardClaims) Valid() error {
	now := time.Now().Unix()

	if !verifyIssuerExists(claims.Issuer) {
		return ErrIssuer
	} else if !verifySubject(claims.Subject) {
		return ErrSubject
	} else if !verifyAudienceExists(claims.Audience) {
		return ErrAudience
	} else if !verifyExpiresAt(claims.ExpiresAt, now) {
		return ErrExpiresAt
	} else if !verifyNotBefore(claims.NotBefore, now) {
		return ErrNotBefore
	} else if !verifyIssuedAt(claims.IssuedAt, now) {
		return ErrIssuedAt
	} else if !verifyID(claims.ID) {
		return ErrID
	}
	return nil
}

func (claims *StandardClaims) MoreValid(issuers []StringOrURI, audience StringOrURI) error {
	if !verifyIssuer(claims.Issuer, issuers) {
		return ErrIssuer
	} else if !verifyAudience(claims.Audience, audience) {
		return ErrAudience
	}
	return nil
}

func verifyIssuerExists(claimed StringOrURI) bool {
	return claimed != ""
}

func verifyIssuer(claimed StringOrURI, validIssuers []StringOrURI) bool {
	for _, iss := range validIssuers {
		if iss == claimed {
			return true
		}
	}
	return false
}

func verifySubject(claimed StringOrURI) bool {
	return claimed != ""
}

func verifyAudienceExists(claimed []StringOrURI) bool {
	return claimed != nil && len(claimed) > 0
}

func verifyAudience(claimed []StringOrURI, validAud StringOrURI) bool {
	for _, aud := range claimed {
		if aud == validAud {
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
