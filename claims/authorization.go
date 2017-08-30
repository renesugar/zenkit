package claims

import "errors"

var (
	AuthorizationIssuer = string("Authorization")
	ValidRoles          = map[string]struct{}{
		"api:access": {},
		"api:admin":  {},
	}

	ErrRoles               = errors.New("invalid roles")
	authorizationClaimsMap = AuthorizationClaimsMap{
		"Issuer":    "",
		"Subject":   "",
		"Audience":  []string{},
		"ExpiresAt": int64(0),
		"NotBefore": int64(0),
		"IssuedAt":  int64(0),
		"ID":        "",
		"Roles":     []string{},
	}
)

type AuthorizationClaimsMap map[string]interface{}

func NewAuthorizationClaimsMap() AuthorizationClaimsMap {
	m := make(map[string]interface{})
	for k, v := range authorizationClaimsMap {
		m[k] = v
	}
	return m
}

// AuthorizationClaims is the expected claims from a token received from Authorization
type AuthorizationClaims struct {
	StandardClaims
	Roles []string `json:"rls"`
}

// Valid determines if a JWT should be rejected or not and implements jwt-go Claims interface
func (claims AuthorizationClaims) Valid() error {
	err := claims.StandardClaims.Valid()
	if err != nil {
		return err
	}
	if !verifyIssuer(claims.Issuer(), []string{AuthorizationIssuer}) {
		return ErrIssuer
	} else if !verifyRoles(claims.Roles, ValidRoles) {
		return ErrRoles
	}
	return nil
}

// MoreValid determines if a JWT should be rejected or not
func (claims AuthorizationClaims) MoreValid(audience string) error {
	if !verifyAudience(claims.Audience(), audience) {
		return ErrAudience
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
