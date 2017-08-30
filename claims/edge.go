package claims

var (
	// EdgeIssuer is the expected issuer of a token that parses to an EdgeClaims
	EdgeIssuer = string("Edge")

	edgeClaimsMap = EdgeClaimsMap{
		"Issuer":    EdgeIssuer,
		"Subject":   "",
		"Audience":  []string{},
		"ExpiresAt": int64(0),
		"NotBefore": int64(0),
		"IssuedAt":  int64(0),
		"ID":        "",
	}
)

func NewEdgeClaimsMap() EdgeClaimsMap {
	m := make(map[string]interface{})
	for k, v := range edgeClaimsMap {
		m[k] = v
	}
	return m
}

type EdgeClaimsMap map[string]interface{}

// EdgeClaims is the expected claims from a token received from Edge that has not yet been authorized
type EdgeClaims struct {
	StandardClaims
}

// Valid determines if a JWT should be rejected or not and implements jwt-go Claims interface
func (claims EdgeClaims) Valid() error {
	err := claims.StandardClaims.Valid()
	if err != nil {
		return err
	}
	if !verifyIssuer(claims.Issuer(), []string{EdgeIssuer}) {
		return ErrIssuer
	}
	return nil
}

// MoreValid determines if a JWT should be rejected or not
func (claims EdgeClaims) MoreValid(audience string) error {
	if !verifyAudience(claims.Audience(), audience) {
		return ErrAudience
	}
	return nil
}
