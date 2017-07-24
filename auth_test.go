package zenkit_test

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"time"

	"github.com/davecgh/go-spew/spew"
	jwtpkg "github.com/dgrijalva/jwt-go"
	"github.com/goadesign/goa"
	"github.com/goadesign/goa/design"
	"github.com/goadesign/goa/design/apidsl"
	"github.com/goadesign/goa/dslengine"
	"github.com/goadesign/goa/middleware/security/jwt"
	. "github.com/zenoss/zenkit"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type testIdentity struct {
	id string
}

func (t *testIdentity) ID() string {
	return t.id
}

var _ = Describe("Auth utilities", func() {

	var (
		id      string
		ident   Identity
		file    *os.File
		token   *jwtpkg.Token
		secret  = []byte("secret")
		svc     *goa.Service
		resp    *httptest.ResponseRecorder
		req     *http.Request
		handler = func(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
			token = jwt.ContextJWT(ctx)
			ident = ContextIdentity(ctx)
			return nil
		}
		emptyMiddleware = func(h goa.Handler) goa.Handler {
			return h
		}
		security = &goa.JWTSecurity{
			In:   goa.LocHeader,
			Name: "Authorization",
		}
	)

	BeforeEach(func() {
		file, _ = ioutil.TempFile("", "zenkit-")
		defer file.Close()
		file.Write(secret)
		svc = goa.New(RandStringRunes(8))
		req, _ = http.NewRequest("", "http://example.com/", nil)
		resp = httptest.NewRecorder()
	})

	AfterEach(func() {
		os.Remove(file.Name())
		file = nil
		token = nil
		ident = nil
	})

	getHandler := func(mw goa.Middleware) goa.Handler {
		wrapper, err := JWTMiddleware(svc, file.Name(), mw, security)
		Ω(err).ShouldNot(HaveOccurred())
		return wrapper(handler)
	}

	signedToken := func() string {
		t := jwtpkg.New(jwtpkg.SigningMethodHS256)
		t.Claims = jwtpkg.MapClaims(map[string]interface{}{
			"sub": id,
		})
		signed, _ := t.SignedString(secret)
		return signed
	}

	assertSecurityError := func(err error, msg string) {
		errResp, ok := err.(*goa.ErrorResponse)
		Ω(ok).Should(BeTrue())
		Ω(errResp.Status).Should(Equal(401))
		Ω(errResp.Code).Should(Equal("jwt_security_error"))
		Ω(errResp.Detail).Should(Equal(msg))
	}

	Context("using common JWT security", func() {
		It("should register the authorization header \"Authorization\"", func() {
			dslengine.Reset()
			apidsl.API("test", func() {
				apidsl.Security(JWT, func() {
				})
			})
			apidsl.Resource("test2", func() {
				apidsl.Action("whatever", func() {
				})
			})
			dslengine.Run()
			d := design.Design
			spew.Dump(d)
		})
	})

	Context("context identity functions", func() {

		It("should be able to pass an identity to the context", func() {
			id = RandStringRunes(8)
			ident = &testIdentity{id}

			ctx := WithIdentity(context.Background(), ident)
			received := ContextIdentity(ctx)
			Ω(received).Should(Equal(ident))
		})
	})

	Context("using the dev mode middleware", func() {

		It("should inject an authorization header when none exists", func() {
			h := getHandler(DefaultJWTValidation)
			err := DevModeMiddleware(h)(context.Background(), resp, req)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(ident).ShouldNot(BeNil())
			Ω(ident.ID()).Should(Equal("developer"))
		})

		It("should respect an existing authorization header", func() {
			id = RandStringRunes(8)
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", signedToken()))
			h := getHandler(DefaultJWTValidation)
			err := DevModeMiddleware(h)(context.Background(), resp, req)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(ident).ShouldNot(BeNil())
			Ω(ident.ID()).Should(Equal(id))
		})

	})

	Context("JWT middleware factory", func() {

		It("should fail to create middleware when the file passed does not exist", func() {
			KeyFileTimeout = 1 * time.Millisecond
			security := &goa.JWTSecurity{}
			_, err := JWTMiddleware(goa.New("hi"), "/this/is/notafile", DefaultJWTValidation, security)
			Ω(err).Should(HaveOccurred())
		})

		It("should create middleware that rejects requests without a token", func() {
			err := getHandler(emptyMiddleware)(context.Background(), resp, req)
			Ω(err).Should(HaveOccurred())
			assertSecurityError(err, `missing header "Authorization"`)
			Ω(token).Should(BeNil())
		})

		It("should create middleware that rejects requests with an invalid token", func() {
			req.Header.Set("Authorization", "Bearer badheader")
			err := getHandler(emptyMiddleware)(context.Background(), resp, req)
			Ω(err).Should(HaveOccurred())
			assertSecurityError(err, "JWT validation failed")
			Ω(token).Should(BeNil())
		})

		It("should create middleware that allows requests with a valid token", func() {
			id = RandStringRunes(8)
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", signedToken()))
			err := getHandler(emptyMiddleware)(context.Background(), resp, req)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(token).ShouldNot(BeNil())
		})

		Context("using the default validator middleware", func() {
			It("should pass the identity through the context", func() {
				id = RandStringRunes(8)
				req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", signedToken()))
				ctx := context.Background()
				err := getHandler(DefaultJWTValidation)(ctx, resp, req)
				Ω(err).ShouldNot(HaveOccurred())
				Ω(ident).ShouldNot(BeNil())
				Ω(ident.ID()).Should(Equal(id))
			})
		})

		Context("using custom validator middleware", func() {
			It("should reject requests that fail the custom middleware", func() {
				TestError := errors.New("test error")
				custom := func(ctx context.Context) error {
					return TestError
				}
				id = RandStringRunes(8)
				req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", signedToken()))
				ctx := context.Background()
				err := getHandler(JWTValidatorFunc(custom))(ctx, resp, req)
				Ω(err).Should(HaveOccurred())
				Ω(err).Should(Equal(TestError))
			})
		})

	})

})
