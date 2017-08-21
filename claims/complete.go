package claims

// CompleteClaims is the expected claims from a token received from Edge that has been authorized
type CompleteClaims struct {
	EdgeClaims
	Token string `json:"tkn"`
}

// Valid determines if a JWT should be rejected or not and implements jwt-go Claims interface
func (claims *CompleteClaims) Valid() error {
	err := claims.EdgeClaims.Valid()
	if err != nil {
		return err
	}
	return nil
}
