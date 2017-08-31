package claims

var (
	// EdgeIssuer is the expected issuer of a token containing an EdgeClaims
	EdgeIssuer = "edge"

	// edgeClaimsMap defines the shape and fields for a map containing
	// JWT claims from Auth0
	edgeClaimsMap = EdgeClaimsMap{
		"iss": EdgeIssuer,
		"sub": "",
		"aud": []string{},
		"exp": int64(0),
		"nbf": int64(0),
		"iat": int64(0),
		"jti": "",
	}
)

// EdgeClaimsMap implements the Claims interface and
// is a map of claims with validation functions
type EdgeClaimsMap map[string]interface{}

// NewEdgeClaimsMap returns a EdgeClaimsMap with keys for Auth0 claims
func NewEdgeClaimsMap() EdgeClaimsMap {
	m := make(map[string]interface{})
	for k, v := range edgeClaimsMap {
		m[k] = v
	}
	return m
}

// EdgeClaimsFromStruct assumes claims is valid and
// creates a EdgeClaimsMap from a EdgeClaims
func EdgeClaimsFromStruct(claims EdgeClaims) EdgeClaimsMap {
	return EdgeClaimsMap{
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
func (m EdgeClaimsMap) Issuer() string {
	return getIssuer(m)
}

// Subject is the jwt "sub" claim
func (m EdgeClaimsMap) Subject() string {
	return getSubject(m)
}

// Audience is the jwt "aud" claim
func (m EdgeClaimsMap) Audience() []string {
	return getAudience(m)
}

// ExpiresAt is the jwt "exp" claim
func (m EdgeClaimsMap) ExpiresAt() int64 {
	return getExpiresAt(m)
}

// NotBefore is the jwt "nbf" claim
func (m EdgeClaimsMap) NotBefore() int64 {
	return getNotBefore(m)
}

// IssuedAt is the jwt "iat" claim
func (m EdgeClaimsMap) IssuedAt() int64 {
	return getIssuedAt(m)
}

// ID is the jwt "jti" claim
func (m EdgeClaimsMap) ID() string {
	return getID(m)
}

// Valid determines if a JWT should be rejected or not and implements jwt-go Claims interface
func (m EdgeClaimsMap) Valid() error {
	if err := Valid(m); err != nil {
		return err
	}
	return ValidateIssuer(m, []string{EdgeIssuer})
}

// Validate checks validity of all fields and verifies
// the claims satisfy the audience
func (m EdgeClaimsMap) Validate(audience string) error {
	if err := m.Valid(); err != nil {
		return err
	}
	return ValidateAudience(m, audience)
}

// EdgeClaims is the expected claims from a token received from Auth0
type EdgeClaims struct {
	StandardClaims
}

// EdgeClaimsFromMap assumes m is Valid and
// creates a EdgeClaims from a EdgeClaimsMap
func EdgeClaimsFromMap(m EdgeClaimsMap) EdgeClaims {
	return EdgeClaims{
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
func (claims EdgeClaims) Valid() error {
	if err := Valid(claims); err != nil {
		return err
	}
	return ValidateIssuer(claims, []string{EdgeIssuer})
}

// Validate checks validity of all fields and verifies
// the claims satisfy the issuers and audience
func (claims EdgeClaims) Validate(audience string) error {
	if err := claims.Valid(); err != nil {
		return err
	}
	return ValidateAudience(claims, audience)
}
