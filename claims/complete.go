package claims

var (
	// CompleteIssuer is the expected issuer of a token containing an CompleteClaims
	CompleteIssuer = "complete"

	// completeClaimsMap defines the shape and fields for a map containing
	// JWT claims from Auth0
	completeClaimsMap = CompleteClaimsMap{
		"iss": CompleteIssuer,
		"sub": "",
		"aud": []string{},
		"exp": int64(0),
		"nbf": int64(0),
		"iat": int64(0),
		"jti": "",
		"tkn": "",
	}
)

// CompleteClaimsMap implements the Claims interface and
// is a map of claims with validation functions
type CompleteClaimsMap map[string]interface{}

// NewCompleteClaimsMap returns a CompleteClaimsMap with keys for Auth0 claims
func NewCompleteClaimsMap() CompleteClaimsMap {
	m := make(map[string]interface{})
	for k, v := range completeClaimsMap {
		m[k] = v
	}
	return m
}

// CompleteClaimsFromStruct assumes claims is valid and
// creates a CompleteClaimsMap from a CompleteClaims
func CompleteClaimsFromStruct(claims CompleteClaims) CompleteClaimsMap {
	return CompleteClaimsMap{
		"iss": claims.Iss,
		"sub": claims.Sub,
		"aud": claims.Aud,
		"exp": claims.Exp,
		"nbf": claims.Nbf,
		"iat": claims.Iat,
		"jti": claims.Jti,
		"tkn": claims.Tkn,
	}
}

// Issuer is the jwt "iss" claim
func (m CompleteClaimsMap) Issuer() string {
	return getIssuer(m)
}

// Subject is the jwt "sub" claim
func (m CompleteClaimsMap) Subject() string {
	return getSubject(m)
}

// Audience is the jwt "aud" claim
func (m CompleteClaimsMap) Audience() []string {
	return getAudience(m)
}

// ExpiresAt is the jwt "exp" claim
func (m CompleteClaimsMap) ExpiresAt() int64 {
	return getExpiresAt(m)
}

// NotBefore is the jwt "nbf" claim
func (m CompleteClaimsMap) NotBefore() int64 {
	return getNotBefore(m)
}

// IssuedAt is the jwt "iat" claim
func (m CompleteClaimsMap) IssuedAt() int64 {
	return getIssuedAt(m)
}

// ID is the jwt "jti" claim
func (m CompleteClaimsMap) ID() string {
	return getID(m)
}

// Token is the complete "tkn" claim
func (m CompleteClaimsMap) Token() string {
	return getToken(m)
}

// Valid determines if a JWT should be rejected or not and implements jwt-go Claims interface
func (m CompleteClaimsMap) Valid() error {
	if err := Valid(m); err != nil {
		return err
	}
	return ValidateIssuer(m, []string{CompleteIssuer})
}

// Validate checks validity of all fields and verifies
// the claims satisfy the audience
func (m CompleteClaimsMap) Validate(audience string) error {
	if err := m.Valid(); err != nil {
		return err
	}
	return ValidateAudience(m, audience)
}

// CompleteClaims is the expected claims from a token received from Auth0
type CompleteClaims struct {
	StandardClaims
	Tkn string `json:"tkn"`
}

// CompleteClaimsFromMap assumes m is Valid and
// creates a CompleteClaims from a CompleteClaimsMap
func CompleteClaimsFromMap(m CompleteClaimsMap) CompleteClaims {
	return CompleteClaims{
		StandardClaims: StandardClaims{
			Iss: m.Issuer(),
			Sub: m.Subject(),
			Aud: m.Audience(),
			Exp: m.ExpiresAt(),
			Nbf: m.NotBefore(),
			Iat: m.IssuedAt(),
			Jti: m.ID(),
		},
		Tkn: m.Token(),
	}
}

// Token is the complete "tkn" claim
func (claims CompleteClaims) Token() string {
	return claims.Tkn
}

// Valid determines if a JWT should be rejected or not and implements jwt-go Claims interface
func (claims CompleteClaims) Valid() error {
	if err := Valid(claims); err != nil {
		return err
	}
	return ValidateIssuer(claims, []string{CompleteIssuer})
}

// Validate checks validity of all fields and verifies
// the claims satisfy the issuers and audience
func (claims CompleteClaims) Validate(audience string) error {
	if err := claims.Valid(); err != nil {
		return err
	}
	return ValidateAudience(claims, audience)
}

func getToken(m map[string]interface{}) string {
	tkn, ok := m["tkn"]
	if !ok {
		return ""
	}
	tknStr, ok := tkn.(string)
	if !ok {
		return ""
	}
	return tknStr
}
