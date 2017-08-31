package claims

import "errors"

var (
	// AuthorizationIssuer is the expected issuer of a token containing
	// an AuthorizationClaims
	AuthorizationIssuer = string("Authorization")
	// ValidRoles are all possible roles that a user may have
	ValidRoles = map[string]struct{}{
		"api:access": {},
		"api:admin":  {},
	}
	// ErrRoles occurs when the rols claim is not valid
	ErrRoles = errors.New("invalid roles")
	// authorizationClaimsMap defines the shape and fields for a map containing
	// JWT claims from authorization
	authorizationClaimsMap = AuthorizationClaimsMap{
		"iss": AuthorizationIssuer,
		"sub": "",
		"aud": []string{},
		"exp": int64(0),
		"nbf": int64(0),
		"iat": int64(0),
		"jti": "",
		"rls": []string{},
	}
)

// AuthorizationClaimsMap implements the Claims interface and
// is a map of claims with validation functions
type AuthorizationClaimsMap map[string]interface{}

// NewAuthorizationClaimsMap returns a AuthorizationClaimsMap with keys for
// authorization claims
func NewAuthorizationClaimsMap() AuthorizationClaimsMap {
	m := make(map[string]interface{})
	for k, v := range authorizationClaimsMap {
		m[k] = v
	}
	return m
}

// AuthorizationClaimsFromStruct assumes claims is valid and
// creates a AuthorizationClaimsMap from a AuthorizationClaims
func AuthorizationClaimsFromStruct(claims AuthorizationClaims) AuthorizationClaimsMap {
	return AuthorizationClaimsMap{
		"iss": claims.Iss,
		"sub": claims.Sub,
		"aud": claims.Aud,
		"exp": claims.Exp,
		"nbf": claims.Nbf,
		"iat": claims.Iat,
		"jti": claims.Jti,
		"rls": claims.Rls,
	}
}

// Issuer is the jwt "iss" claim
func (m AuthorizationClaimsMap) Issuer() string {
	return getIssuer(m)
}

// Subject is the jwt "sub" claim
func (m AuthorizationClaimsMap) Subject() string {
	return getSubject(m)
}

// Audience is the jwt "aud" claim
func (m AuthorizationClaimsMap) Audience() []string {
	return getAudience(m)
}

// ExpiresAt is the jwt "exp" claim
func (m AuthorizationClaimsMap) ExpiresAt() int64 {
	return getExpiresAt(m)
}

// NotBefore is the jwt "nbf" claim
func (m AuthorizationClaimsMap) NotBefore() int64 {
	return getNotBefore(m)
}

// IssuedAt is the jwt "iat" claim
func (m AuthorizationClaimsMap) IssuedAt() int64 {
	return getIssuedAt(m)
}

// ID is the jwt "jti" claim
func (m AuthorizationClaimsMap) ID() string {
	return getID(m)
}

// Roles is the authorization "rls" claim
func (m AuthorizationClaimsMap) Roles() []string {
	return getRoles(m)
}

// Valid determines if a JWT should be rejected or not and implements jwt-go Claims interface
func (m AuthorizationClaimsMap) Valid() error {
	if err := Valid(m); err != nil {
		return err
	}
	return ValidateIssuer(m, []string{AuthorizationIssuer})
}

// Validate checks validity of all fields and verifies
// the claims satisfy the audience
func (m AuthorizationClaimsMap) Validate(audience string) error {
	if err := m.Valid(); err != nil {
		return err
	}
	return ValidateAudience(m, audience)
}

// AuthorizationClaims is the expected claims from a token received from Authorization
type AuthorizationClaims struct {
	StandardClaims
	Rls []string `json:"rls"`
}

// AuthorizationClaimsFromMap assumes m is Valid and
// creates a AuthorizationClaims from a AuthorizationClaimsMap
func AuthorizationClaimsFromMap(m AuthorizationClaimsMap) AuthorizationClaims {
	return AuthorizationClaims{
		StandardClaims: StandardClaims{
			Iss: m.Issuer(),
			Sub: m.Subject(),
			Aud: m.Audience(),
			Exp: m.ExpiresAt(),
			Nbf: m.NotBefore(),
			Iat: m.IssuedAt(),
			Jti: m.ID(),
		},
		Rls: m.Roles(),
	}
}

// Roles is the authorization "rls" claim
func (claims AuthorizationClaims) Roles() []string {
	return claims.Rls
}

// Valid determines if a JWT should be rejected or not and implements jwt-go Claims interface
func (claims AuthorizationClaims) Valid() error {
	if err := Valid(claims); err != nil {
		return err
	} else if err := ValidateIssuer(claims, []string{AuthorizationIssuer}); err != nil {
		return err
	}
	return ValidateRoles(claims, ValidRoles)
}

// Validate checks validity of all fields and verifies
// the claims satisfy the audience
func (claims AuthorizationClaims) Validate(audience string) error {
	if err := claims.Valid(); err != nil {
		return err
	}
	return ValidateAudience(claims, audience)
}

// ValidateRoles checks that the roles in claims exist in roles
func ValidateRoles(claims AuthorizationClaims, roles map[string]struct{}) error {
	if !verifyRoles(claims.Roles(), roles) {
		return ErrRoles
	}
	return nil
}

func verifyRoles(claimed []string, validRoles map[string]struct{}) bool {
	for _, claim := range claimed {
		if _, ok := validRoles[claim]; !ok {
			return false
		}
	}
	return true
}

func getRoles(m map[string]interface{}) []string {
	rls, ok := m["rls"]
	if !ok {
		return []string{}
	}
	switch rlsSlc := rls.(type) {
	case []interface{}:
		rlsStrs := make([]string, len(rlsSlc))
		for n, v := range rlsSlc {
			val, ok := v.(string)
			if !ok {
				return []string{}
			}
			rlsStrs[n] = val
		}
		return rlsStrs
	case []string:
		return rlsSlc
	}
	return []string{}
}
