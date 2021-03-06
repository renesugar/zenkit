package claims

import "time"

var (
	// standardClaimsMap defines the shape and fields for a map containing
	// standard JWT claims
	standardClaimsMap = StandardClaimsMap{
		"iss": "",
		"sub": "",
		"aud": []string{},
		"exp": int64(0),
		"iat": int64(0),
		"jti": "",
	}
)

// NOTE: At the time of writing this, 8-29-17, goa jwt security cannot use
// claim structs other than a jwt-go MapClaims

// StandardClaimsMap implements the Claims interface and
// is a map of claims with validation functions
type StandardClaimsMap map[string]interface{}

// NewStandardClaimsMap returns a StandardClaimsMap with keys for standard claims
func NewStandardClaimsMap() StandardClaimsMap {
	m := make(map[string]interface{})
	for k, v := range standardClaimsMap {
		m[k] = v
	}
	return m
}

// StandardClaimsFromStruct assumes claims is valid and
// creates a StandardClaimsMap from a StandardClaims
func StandardClaimsFromStruct(claims StandardClaims) StandardClaimsMap {
	return StandardClaimsMap{
		"iss": claims.Iss,
		"sub": claims.Sub,
		"aud": claims.Aud,
		"exp": claims.Exp,
		"iat": claims.Iat,
	}
}

// Issuer is the jwt "iss" claim
func (m StandardClaimsMap) Issuer() string {
	return GetIssuer(m)
}

// Subject is the jwt "sub" claim
func (m StandardClaimsMap) Subject() string {
	return GetSubject(m)
}

// Audience is the jwt "aud" claim
func (m StandardClaimsMap) Audience() []string {
	return GetAudience(m)
}

// ExpiresAt is the jwt "exp" claim
func (m StandardClaimsMap) ExpiresAt() int64 {
	return GetExpiresAt(m)
}

// IssuedAt is the jwt "iat" claim
func (m StandardClaimsMap) IssuedAt() int64 {
	return GetIssuedAt(m)
}

// Valid verifies that mandatory claims exist and are valid
func (m StandardClaimsMap) Valid() error {
	return ValidateClaims(m)
}

// Validate checks validity of all fields and verifies
// the claims satisfy the issuers and audience
func (m StandardClaimsMap) Validate(issuers []string, audience string) error {
	if err := m.Valid(); err != nil {
		return err
	}
	return ValidateIssuer(m, issuers)
}

// StandardClaims implements registered claim names according to RFC 7519
type StandardClaims struct {
	Iss string   `json:"iss"`
	Sub string   `json:"sub"`
	Aud []string `json:"aud"`
	Exp int64    `json:"exp"`
	Iat int64    `json:"iat"`
}

// NewStandardClaims creates a new StandardClaims
func NewStandardClaims(iss, sub string, aud []string, validDuration time.Duration) StandardClaims {
	now := time.Now()
	return StandardClaims{
		Iss: iss,
		Sub: sub,
		Aud: aud,
		Exp: now.Add(validDuration).Unix(),
		Iat: now.Unix(),
	}
}

// StandardClaimsFromMap assumes m is Valid and
// creates a StandardClaims from a StandardClaimsMap
func StandardClaimsFromMap(m StandardClaimsMap) StandardClaims {
	return StandardClaims{
		Iss: m.Issuer(),
		Sub: m.Subject(),
		Aud: m.Audience(),
		Exp: m.ExpiresAt(),
		Iat: m.IssuedAt(),
	}
}

// Issuer is the jwt "iss" claim
func (claims StandardClaims) Issuer() string {
	return claims.Iss
}

// Subject is the jwt "sub" claim
func (claims StandardClaims) Subject() string {
	return claims.Sub
}

// Audience is the jwt "aud" claim
func (claims StandardClaims) Audience() []string {
	return claims.Aud
}

// ExpiresAt is the jwt "exp" claim
func (claims StandardClaims) ExpiresAt() int64 {
	return claims.Exp
}

// IssuedAt is the jwt "iat" claim
func (claims StandardClaims) IssuedAt() int64 {
	return claims.Iat
}

// Valid determines if a JWT should be rejected or not and implements jwt-go Claims interface
func (claims StandardClaims) Valid() error {
	return ValidateClaims(claims)
}

// Validate checks validity of all fields and verifies
// the claims satisfy the issuers and audience
func (claims StandardClaims) Validate(issuers []string, audience string) error {
	if err := claims.Valid(); err != nil {
		return err
	}
	return ValidateIssuer(claims, issuers)
}
