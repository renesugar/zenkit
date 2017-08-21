package claims

// AuthZeroIssuer is the expected issuer of a token that parses to an AuthZeroClaims
var AuthZeroIssuer = StringOrURI("Auth0")

// AuthZeroClaims is the expected claims from a token received from Auth0
type AuthZeroClaims struct {
	StandardClaims
}

// Valid determines if a JWT should be rejected or not and implements jwt-go Claims interface
func (claims *AuthZeroClaims) Valid() error {
	err := claims.StandardClaims.Valid()
	if err != nil {
		return err
	}
	if !verifyIssuer(claims.Issuer, []StringOrURI{AuthZeroIssuer}) {
		return ErrIssuer
	}
	return nil
}

// MoreValid determines if a JWT should be rejected or not
func (claims *AuthZeroClaims) MoreValid(audience StringOrURI) error {
	if !verifyAudience(claims.Audience, audience) {
		return ErrAudience
	}
	return nil
}
