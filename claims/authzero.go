package claims

// AuthZeroIssuer is the expected issuer of a token that parses to an AuthZeroClaims
var (
	AuthZeroIssuer = string("Auth0")

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

type AuthZeroClaimsMap map[string]interface{}

func NewAuthZeroClaimsMap() AuthZeroClaimsMap {
	m := make(map[string]interface{})
	for k, v := range authZeroClaimsMap {
		m[k] = v
	}
	return m
}

// Valid determines if a JWT should be rejected or not and implements jwt-go Claims interface
func (claims AuthZeroClaimsMap) Valid() error {
	return Valid(claims)
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

// Validate checks that the claim is valid and issuers and audience are satisfied
func (claims AuthZeroClaimsMap) Validate(audience string) error {
	if err := Valid(claims); err != nil {
		return err
	} else if err := ValidateIssuer(claims, []string{AuthZeroIssuer}); err != nil {
		return err
	}
	return ValidateAudience(claims, audience)
}

func (m AuthZeroClaimsMap) ToAuthZeroClaims() (AuthZeroClaims, error) {
	var claims AuthZeroClaims
	err := m.Valid()
	if err != nil {
		return claims, err
	}
	claims = AuthZeroClaims{
		StandardClaims: StandardClaims{
			Iss: m["iss"].(string),
		},
	}
	return claims, nil
}

// AuthZeroClaims is the expected claims from a token received from Auth0
type AuthZeroClaims struct {
	StandardClaims
}

// Valid determines if a JWT should be rejected or not and implements jwt-go Claims interface
func (claims AuthZeroClaims) Valid() error {
	return claims.StandardClaims.Valid()
}

// Validate checks that the claim is valid and issuers and audience are satisfied
func (claims AuthZeroClaims) Validate(audience string) error {
	if err := Valid(claims); err != nil {
		return err
	} else if err := ValidateIssuer(claims, []string{AuthZeroIssuer}); err != nil {
		return err
	}
	return ValidateAudience(claims, audience)
}
