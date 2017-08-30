package claims

var (
	completeClaimsMap = CompleteClaimsMap{
		"Issuer":    "",
		"Subject":   "",
		"Audience":  []string{},
		"ExpiresAt": int64(0),
		"NotBefore": int64(0),
		"IssuedAt":  int64(0),
		"ID":        "",
		"Token":     "",
	}
)

type CompleteClaimsMap map[string]interface{}

func NewCompleteClaimsMap() CompleteClaimsMap {
	m := make(map[string]interface{})
	for k, v := range completeClaimsMap {
		m[k] = v
	}
	return m
}

// CompleteClaims is the expected claims from a token received from Edge that has been authorized
type CompleteClaims struct {
	EdgeClaims
	Token string `json:"tkn"`
}

// Valid determines if a JWT should be rejected or not and implements jwt-go Claims interface
func (claims CompleteClaims) Valid() error {
	err := claims.EdgeClaims.Valid()
	if err != nil {
		return err
	}
	return nil
}
