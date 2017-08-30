package claims_test

import (
	"time"

	. "github.com/zenoss/zenkit/claims"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Authoirzation Claims", func() {
	var (
		now           = time.Now().Unix()
		defaultClaims = AuthorizationClaims{
			StandardClaims: StandardClaims{
				Issuer:    AuthorizationIssuer,
				Subject:   "abcd",
				Audience:  []string{"tester"},
				ExpiresAt: now + int64(time.Hour),
				NotBefore: now,
				IssuedAt:  now,
				ID:        "0",
			},
			Roles: []string{
				"api:access",
				"api:admin",
			},
		}
		claims        = defaultClaims
		validAudience = "tester"
	)
	BeforeEach(func() {
		claims = defaultClaims
	})
	Context("when standard claims do not validate", func() {
		BeforeEach(func() {
			claims.Subject = ""
		})
		It("should return an error", func() {
			err := claims.Valid()
			Ω(err).ShouldNot(BeNil())
		})
	})
	Context("when validating fields with criteria", func() {
		Context("when the issuer is not valid", func() {
			BeforeEach(func() {
				claims.Issuer = "keanu"
			})
			It("should return an error", func() {
				err := claims.Valid()
				Ω(err).Should(Equal(ErrIssuer))
			})
		})
		Context("when the audience is not valid", func() {
			BeforeEach(func() {
				claims.Audience = []string{"keanu"}
			})
			It("should return an error", func() {
				err := claims.MoreValid(validAudience)
				Ω(err).Should(Equal(ErrAudience))
			})
		})
		Context("when the roles are not valid", func() {
			BeforeEach(func() {
				claims.Roles = []string{"action hero", "the one"}
			})
			It("should return an error", func() {
				err := claims.Valid()
				Ω(err).Should(Equal(ErrRoles))
			})
		})
		Context("when the claims are just right", func() {
			It("should validate with no errors", func() {
				err := claims.Valid()
				Ω(err).Should(BeNil())
				err = claims.MoreValid(validAudience)
				Ω(err).Should(BeNil())
			})
		})
	})
})
