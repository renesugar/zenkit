package auth_test

import (
	"context"
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"time"

	jwtgo "github.com/dgrijalva/jwt-go"
	"github.com/goadesign/goa"
	"github.com/goadesign/goa/client"
	"github.com/goadesign/goa/design"
	"github.com/goadesign/goa/design/apidsl"
	"github.com/goadesign/goa/dslengine"
	"github.com/pkg/errors"
	. "github.com/zenoss/zenkit/auth"
	"github.com/zenoss/zenkit/claims"
	"github.com/zenoss/zenkit/logging"
	"github.com/zenoss/zenkit/test"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type testIdentity struct {
	id     string
	tenant string
}

func (t *testIdentity) ID() string {
	return t.id
}

func (t *testIdentity) Tenant() string {
	return t.tenant
}

type errorLogger struct {
	Buf string
}

func (logger *errorLogger) LogError(msg string, keys ...interface{}) {
	logger.Buf += fmt.Sprint(msg)
}

func newClaims(id string, audience []string) claims.StandardClaims {
	now := time.Now()
	return claims.StandardClaims{
		Iss: "test",
		Sub: id,
		Aud: audience,
		Exp: now.Add(time.Hour).Unix(),
		Iat: now.Unix(),
	}
}

func newClaimsMap(id string, audience []string) claims.StandardClaimsMap {
	return claims.StandardClaimsFromStruct(newClaims(id, audience))
}

type mockSigningMethod struct {
	signature string
	err       error
}

func (mock *mockSigningMethod) Verify(signingString, signature string, key interface{}) error {
	// unused
	return nil
}

func (mock *mockSigningMethod) Sign(signingString string, key interface{}) (string, error) {
	return mock.signature, mock.err
}

func (mock *mockSigningMethod) Alg() string {
	// unused
	return "mock"
}

var _ = Describe("Auth utilities", func() {

	var (
		id           string
		audience     []string
		ident        TenantIdentity
		file         *os.File
		secret       = []byte("secret")
		rsaPublicKey = []byte(`-----BEGIN CERTIFICATE-----
MIIC+zCCAeOgAwIBAgIJVuc8LVr4E501MA0GCSqGSIb3DQEBCwUAMBsxGTAXBgNV
BAMTEGptYXRvcy5hdXRoMC5jb20wHhcNMTcwNzE2MTcwNzEwWhcNMzEwMzI1MTcw
NzEwWjAbMRkwFwYDVQQDExBqbWF0b3MuYXV0aDAuY29tMIIBIjANBgkqhkiG9w0B
AQEFAAOCAQ8AMIIBCgKCAQEA9v9VOd3ZgbtUnvaNTRvrkMbq+bkrH5dVr/tfKSD6
FSEfCKIo9Mr0hh3maXzXdKvrx+qHuISr+yivqTkOtQUaEsWPK8v2tp+qsqnPCXsi
1kB4LOkq6MIzcCZ5d3b8/Z8wHRpuDWYlhvRLpyMmHzxCyX6KAERiDsxSmyRllY+O
//4Z0ieA8F9ixVtLEKcPimLMk4eX3Xv7eVIe6WgMcDe56JQEFCHGdIBL7h5zARKl
JdCinivfYmUcUfKnJ5b+lYvqr5zMP4XxZ0wzz073Yy0QsNkuJzWjtBcTwVFrkyzG
Dmdq0AUYcSAE3Ez5cLqEBbbfOTdzAyjzWRpNmEG3uwiCGwIDAQABo0IwQDAPBgNV
HRMBAf8EBTADAQH/MB0GA1UdDgQWBBRBp5DJ126Mi8ZdAIM8FQ4Z+woJyTAOBgNV
HQ8BAf8EBAMCAoQwDQYJKoZIhvcNAQELBQADggEBALZlvTtIW4MzZV84Bp+lZ91J
FsaduYohBjeTxuIz38uWHFYPTpJoKHwMS9yaCm4psOt3nQN8ipil2OblUHb4Pi9X
F+b5j4TfxD9Uc6vOnzVYk1GnLFny/Sl41QDUqg78cNE81Li47pw/RfWjSkdXDa1q
PZ7f7nhaGd0pr6KF1z/GJUA2IpgsZ/pzJmAO3BZMAFfzp3u2kpBRry+BUXf5xg+3
xhcmeuiFwygRmLe2q0SQ1n6ekrw+RcIHfsWxqq6A028/N8GqGdbcJ5qL5ITEKJZT
BMUjCjMj7krg2mdNb3PmGN97AtEelKgC8RRdlswCdPQkFVQq2tBfPXrckdMHO18=
-----END CERTIFICATE-----`)
		rsaFile   *os.File
		svc       *goa.Service
		resp      *httptest.ResponseRecorder
		req       *http.Request
		claimsCtx context.Context
		// handler   = func(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
		// 	ident = ContextTenantIdentity(ctx)
		// 	return nil
		// }
		// emptyMiddleware = func(h goa.Handler) goa.Handler {
		// 	return h
		// }
		// DevModeMiddleware = NewDevJWTMiddleware(
		// 	newClaims("abcd", []string{"tests"}),
		// 	jwtgo.SigningMethodNone,
		// 	string(jwtgo.UnsafeAllowNoneSignatureType),
		// )
	)

	BeforeEach(func() {
		svc = goa.New(test.RandString(8))
		svc.WithLogger(logging.ServiceLogger())
		req, _ = http.NewRequest("", "http://example.com/", nil)
		resp = httptest.NewRecorder()
		audience = []string{"tenant"}
		claimsCtx = context.Background()
	})

	JustBeforeEach(func() {
		file, _ = ioutil.TempFile("", "zenkit-")
		defer file.Close()
		file.Write(secret)

		rsaFile, _ = ioutil.TempFile("", "zenkit-")
		defer rsaFile.Close()
		rsaFile.Write(rsaPublicKey)
	})

	AfterEach(func() {
		os.Remove(file.Name())
		os.Remove(rsaFile.Name())
		file = nil
		rsaFile = nil
		ident = nil
	})

	// getHandler := func(mw goa.Middleware) goa.Handler {
	// 	wrapper, err := goa.NewMiddleware(JWTMiddleware)
	// 	Ω(err).ShouldNot(HaveOccurred())
	// 	return wrapper(handler)
	// }

	signedToken := func() string {
		stdClaims := newClaims(id, audience)
		t := jwtgo.NewWithClaims(jwtgo.SigningMethodHS256, stdClaims)
		signed, _ := t.SignedString(secret)
		return signed
	}

	// assertSecurityError := func(err error, msg string) {
	// 	err = errors.Cause(err)
	// 	errResp, ok := err.(*goa.ErrorResponse)
	// 	Ω(ok).Should(BeTrue())
	// 	Ω(errResp.Status).Should(Equal(401))
	// 	Ω(errResp.Code).Should(Equal("jwt_security_error"))
	// 	Ω(errResp.Detail).Should(Equal(msg))
	// }

	Context("when reading a key from the fs", func() {
		Context("when the file does not exist", func() {
			It("should log retries and return an error", func() {
				KeyFileTimeout = 1 * time.Millisecond
				logger := &errorLogger{}
				_, err := ReadKeyFromFS(logger, "/does/not/exist")
				Ω(err).Should(HaveOccurred())
				Ω(logger.Buf).ShouldNot(BeNil())
			})
		})
		Context("when the file does exist", func() {
			It("should get the key written in the file", func() {
				KeyFileTimeout = 1 * time.Millisecond
				logger := &errorLogger{}
				key, err := ReadKeyFromFS(logger, file.Name())
				Ω(err).Should(BeNil())
				Ω(key).Should(Equal(secret))
			})
		})
	})

	Context("when getting multiple keys from the FS", func() {
		logger := &errorLogger{}
		Context("when one of the keys cannot be read", func() {
			It("should return an error", func() {
				_, err := GetKeysFromFS(logger, []string{"/not/here"})
				Ω(err).Should(HaveOccurred())
			})
		})
		Context("when the keys can be read", func() {
			It("should return some keys", func() {
				keys, err := GetKeysFromFS(logger, []string{file.Name(), rsaFile.Name()})
				Ω(err).ShouldNot(HaveOccurred())
				Ω(keys).Should(HaveLen(2))
			})
		})
	})

	Context("when trying to convert bytes to a key", func() {
		Context("with a public rsa key", func() {
			It("should return an RSA key in a jwt.Key", func() {
				var pubKeyType *rsa.PublicKey
				key := ConvertToKey(rsaPublicKey)
				Ω(key).Should(BeAssignableToTypeOf(pubKeyType))
			})
		})
		Context("with a secret", func() {
			It("should return some bytes", func() {
				var bytes []byte
				key := ConvertToKey(secret)
				Ω(key).Should(BeAssignableToTypeOf(bytes))
			})
		})
	})

	Context("using common JWT security", func() {
		It(fmt.Sprintf("should register the authorization header \"%\"", AuthorizationHeader), func() {
			dslengine.Reset()
			apidsl.API("test", func() {
				apidsl.Security(JWT(), func() {})
			})
			dslengine.Run()
			d := design.Design
			Ω(d.SecuritySchemes).Should(HaveLen(1))
			Ω(d.SecuritySchemes[0].Name).Should(Equal(AuthorizationHeader))
		})
	})

	Context("when using TenantIdentity", func() {
		It("should be able to pass an identity to the context", func() {
			id = test.RandString(8)
			ident = &testIdentity{id, "tenant"}

			ctx := WithTenantIdentity(context.Background(), ident)
			received := ContextTenantIdentity(ctx)
			Ω(received).Should(Equal(ident))
		})

		It("should return nil if there is no identity on the context", func() {
			Ω(ContextTenantIdentity(context.Background())).Should(BeNil())
		})
	})

	Context("when creating an Auth0TenantIdentity", func() {
		It("should parse the subject to get the user id", func() {

		})
	})

	Context("when building a dev token", func() {
		It("should populate claims", func() {
			signedToken, err := BuildToken(
				newClaims("1", []string{"anyone"}),
				jwtgo.SigningMethodNone,
				jwtgo.UnsafeAllowNoneSignatureType)
			Ω(err).Should(BeNil())
			Ω(signedToken).ShouldNot(Equal(""))

			actualClaims := jwtgo.MapClaims{}
			data, err := jwtgo.DecodeSegment(strings.Split(signedToken, ".")[1])
			Ω(err).Should(BeNil())
			err = json.Unmarshal(data, &actualClaims)
			Ω(err).Should(BeNil())

			Ω(actualClaims["sub"]).Should(Equal("1"))

			audience := actualClaims["aud"].([]interface{})
			Ω(len(audience)).Should(Equal(1))
			Ω(audience[0].(string)).Should(Equal("anyone"))
		})
	})

	// Context("using the dev mode middleware", func() {
	//
	// 	It("should inject an authorization header when none exists", func() {
	// 		h := getHandler(JWTMiddleware)
	// 		err := DevModeMiddleware(h)(context.Background(), resp, req)
	// 		Ω(err).ShouldNot(HaveOccurred())
	// 		Ω(ident).ShouldNot(BeNil())
	// 		Ω(ident.ID()).Should(Equal("1"))
	// 	})
	//
	// 	It("should respect an existing authorization header", func() {
	// 		id = test.RandString(8)
	// 		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", signedToken()))
	// 		h := getHandler(JWTMiddleware)
	// 		err := DevModeMiddleware(h)(svc.Context, resp, req)
	// 		Ω(err).ShouldNot(HaveOccurred())
	// 		Ω(ident).ShouldNot(BeNil())
	// 		Ω(ident.ID()).Should(Equal(id))
	// 	})
	//
	// })

	Context("signing with a dynamic signer", func() {
		claimsFunc := func() (jwtgo.Claims, error) {
			ctxIdent := ContextTenantIdentity(claimsCtx)
			return newClaims(ctxIdent.ID(), []string{ctxIdent.Tenant()}), nil
		}
		signer := &DynamicSigner{}
		BeforeEach(func() {
			claimsCtx = WithTenantIdentity(context.Background(), &testIdentity{
				id,
				audience[0],
			})
			signer = NewSigner(claimsFunc, jwtgo.SigningMethodHS256, secret)
			req.Header.Set("Authorization", "abc123")
		})
		Context("when the claimsFunc returns a claims", func() {
			BeforeEach(func() {
				signer.ClaimsFunc = claimsFunc
			})
			Context("with an invalid secret", func() {
				BeforeEach(func() {
					signer.Method = jwtgo.SigningMethodRS256
					signer.Secret = []byte("secret")
				})
				It("signing should return an error", func() {
					err := signer.Sign(req)
					Ω(err).Should(HaveOccurred())
				})
			})
			Context("with a valid secret", func() {
				BeforeEach(func() {
					signer.Secret = secret
					id = "abcd"
					audience = []string{"tenant"}
					claimsCtx = WithTenantIdentity(context.Background(), &testIdentity{
						id,
						audience[0],
					})
				})
				It("should return a Signer, that puts a good token in request headers", func() {
					err := signer.Sign(req)
					Ω(err).ShouldNot(HaveOccurred())
					Ω(req.Header["Authorization"][0]).Should(Equal("Bearer " + signedToken()))
				})
			})
		})
		Context("when the claimsFunc returns an error", func() {
			badClaimsFunc := func() (jwtgo.Claims, error) {
				return nil, errors.New("some error")
			}
			BeforeEach(func() {
				signer.ClaimsFunc = badClaimsFunc
			})
			Context("with a valid secret", func() {
				BeforeEach(func() {
					signer.Secret = []byte("secret")
				})
				It("should return an error", func() {
					err := signer.Sign(req)
					Ω(err).Should(HaveOccurred())
				})
			})
		})
	})

	Context("using the jwt signer", func() {
		It("should return nil if there is no auth on the header", func() {
			signer := JWTSigner(req)
			Ω(signer).Should(BeNil())
		})

		It("should not set the token type if none is specified on the header", func() {
			req.Header.Set("Authorization", "abc123")
			signer := JWTSigner(req)
			Ω(signer).ShouldNot(BeNil())

			token, err := signer.TokenSource.Token()
			Ω(err).ShouldNot(HaveOccurred())

			staticToken, ok := token.(*client.StaticToken)
			Ω(ok).Should(BeTrue())
			Ω(staticToken.Value).To(Equal("abc123"))
		})

		It("should set the token and the type if it is specified on the header", func() {
			req.Header.Set("Authorization", "Bearer abc123")
			signer := JWTSigner(req)
			Ω(signer).ShouldNot(BeNil())

			token, err := signer.TokenSource.Token()
			Ω(err).ShouldNot(HaveOccurred())

			staticToken, ok := token.(*client.StaticToken)
			Ω(ok).Should(BeTrue())
			Ω(staticToken.Value).To(Equal("abc123"))
			Ω(staticToken.Type).To(Equal("Bearer"))
		})
	})

})
