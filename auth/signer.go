package auth

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/goadesign/goa/client"
	"github.com/pkg/errors"

	jwtgo "github.com/dgrijalva/jwt-go"
)

// ClaimsFunc is used by DynamicSigner to get claims at signing time
type ClaimsFunc func() (jwtgo.Claims, error)

// DynamicSigner is a goa client.Signer that gets claims at signing time,
// which allows a single signer to be used for multiple identities
type DynamicSigner struct {
	ClaimsFunc ClaimsFunc
	Method     jwtgo.SigningMethod
	Secret     []byte
	KeyName    string
	KeyFormat  string
}

// NewSigner creates a DynamicSigner for the claimsFunc, signing method, and secret
func NewSigner(claimsFunc ClaimsFunc, signingMethod jwtgo.SigningMethod, secret []byte) *DynamicSigner {
	return &DynamicSigner{
		ClaimsFunc: claimsFunc,
		Method:     signingMethod,
		Secret:     secret,
		KeyName:    "Authorization",
		KeyFormat:  "Bearer %s",
	}
}

// Sign signs the request header and satisfies the goa client Signer interface
func (signer *DynamicSigner) Sign(r *http.Request) error {
	claims, err := signer.ClaimsFunc()
	if err != nil {
		return errors.Wrap(err, "unable to get claims")
	}

	token := jwtgo.NewWithClaims(signer.Method, claims)
	key, err := token.SignedString(signer.Secret)
	if err != nil {
		return errors.Wrap(err, "unable to sign with provided secret")
	}

	r.Header.Set(signer.KeyName, fmt.Sprintf(signer.KeyFormat, key))
	return nil
}

// JWTSigner returns a signer that signs with the jwt on req
func JWTSigner(req *http.Request) *client.JWTSigner {
	token := &client.StaticToken{}

	parts := strings.Fields(req.Header.Get("Authorization"))
	switch len(parts) {
	case 0:
		return nil
	case 1:
		token.Value = parts[0]
	default:
		token.Type = parts[0]
		token.Value = strings.Join(parts[1:], " ")
	}
	return &client.JWTSigner{TokenSource: &client.StaticTokenSource{StaticToken: token}}
}
