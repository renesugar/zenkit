package claims

// EdgeIssuer is the expected issuer of a token that parses to an EdgeClaims
var EdgeIssuer = StringOrURI("Edge")

// EdgeClaims is the expected claims from a token received from Edge that has not yet been authorized
type EdgeClaims struct {
	StandardClaims
}

// Valid determines if a JWT should be rejected or not and implements jwt-go Claims interface
func (claims *EdgeClaims) Valid() error {
	err := claims.StandardClaims.Valid()
	if err != nil {
		return err
	}
	if !verifyIssuer(claims.Issuer, []StringOrURI{EdgeIssuer}) {
		return ErrIssuer
	}
	return nil
}

// MoreValid determines if a JWT should be rejected or not
func (claims *EdgeClaims) MoreValid(audience StringOrURI) error {
	if !verifyAudience(claims.Audience, audience) {
		return ErrAudience
	}
	return nil
}
