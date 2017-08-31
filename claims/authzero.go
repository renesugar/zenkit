package claims

var (
	// AuthZeroIssuer is the expected issuer of a token containing an AuthZeroClaims
	AuthZeroIssuer = "Auth0"

	// authZeroClaimsMap defines the shape and fields for a map containing
	// JWT claims from Auth0
	authZeroClaimsMap = AuthZeroClaimsMap{
		"iss": AuthZeroIssuer,
		"sub": "",
		"aud": []string{},
		"exp": int64(0),
		"nbf": int64(0),
		"iat": int64(0),
		"jti": "",
	}
)

// AuthZeroClaimsMap implements the Claims interface and
// is a map of claims with validation functions
type AuthZeroClaimsMap map[string]interface{}

// NewAuthZeroClaimsMap returns a AuthZeroClaimsMap with keys for Auth0 claims
func NewAuthZeroClaimsMap() AuthZeroClaimsMap {
	m := make(map[string]interface{})
	for k, v := range authZeroClaimsMap {
		m[k] = v
	}
	return m
}

// AuthZeroClaimsFromStruct assumes claims is valid and
// creates a AuthZeroClaimsMap from a AuthZeroClaims
func AuthZeroClaimsFromStruct(claims AuthZeroClaims) AuthZeroClaimsMap {
	return AuthZeroClaimsMap{
		"iss": claims.Iss,
		"sub": claims.Sub,
		"aud": claims.Aud,
		"exp": claims.Exp,
		"nbf": claims.Nbf,
		"iat": claims.Iat,
		"jti": claims.Jti,
	}
}

// Issuer is the jwt "iss" claim
func (m AuthZeroClaimsMap) Issuer() string {
	return getIssuer(m)
}

// Subject is the jwt "sub" claim
func (m AuthZeroClaimsMap) Subject() string {
	return getSubject(m)
}

// Audience is the jwt "aud" claim
func (m AuthZeroClaimsMap) Audience() []string {
	return getAudience(m)
}

// ExpiresAt is the jwt "exp" claim
func (m AuthZeroClaimsMap) ExpiresAt() int64 {
	return getExpiresAt(m)
}

// NotBefore is the jwt "nbf" claim
func (m AuthZeroClaimsMap) NotBefore() int64 {
	return getNotBefore(m)
}

// IssuedAt is the jwt "iat" claim
func (m AuthZeroClaimsMap) IssuedAt() int64 {
	return getIssuedAt(m)
}

// ID is the jwt "jti" claim
func (m AuthZeroClaimsMap) ID() string {
	return getID(m)
}

// Valid determines if a JWT should be rejected or not and implements jwt-go Claims interface
func (m AuthZeroClaimsMap) Valid() error {
	if err := Valid(m); err != nil {
		return err
	}
	return ValidateIssuer(m, []string{AuthZeroIssuer})
}

// Validate checks validity of all fields and verifies
// the claims satisfy the audience
func (m AuthZeroClaimsMap) Validate(audience string) error {
	if err := m.Valid(); err != nil {
		return err
	}
	return ValidateAudience(m, audience)
}

// AuthZeroClaims is the expected claims from a token received from Auth0
type AuthZeroClaims struct {
	StandardClaims
}

// AuthZeroClaimsFromMap assumes m is Valid and
// creates a AuthZeroClaims from a AuthZeroClaimsMap
func AuthZeroClaimsFromMap(m AuthZeroClaimsMap) AuthZeroClaims {
	return AuthZeroClaims{
		StandardClaims{
			Iss: m.Issuer(),
			Sub: m.Subject(),
			Aud: m.Audience(),
			Exp: m.ExpiresAt(),
			Nbf: m.NotBefore(),
			Iat: m.IssuedAt(),
			Jti: m.ID(),
		},
	}
}

// Valid determines if a JWT should be rejected or not and implements jwt-go Claims interface
func (claims AuthZeroClaims) Valid() error {
	if err := Valid(claims); err != nil {
		return err
	}
	return ValidateIssuer(claims, []string{AuthZeroIssuer})
}

// Validate checks validity of all fields and verifies
// the claims satisfy the issuers and audience
func (claims AuthZeroClaims) Validate(audience string) error {
	if err := claims.Valid(); err != nil {
		return err
	}
	return ValidateAudience(claims, audience)
}
