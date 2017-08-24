package claims_test

import (
	"time"

	. "github.com/zenoss/zenkit/claims"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Authoirzation Claims", func() {
	var (
		now   = time.Now().Unix()
		claim = AuthorizationClaims{
			StandardClaims: StandardClaims{
				Issuer:    AuthorizationIssuer,
				Subject:   StringOrURI("abcd"),
				Audience:  []StringOrURI{StringOrURI("tester")},
				ExpiresAt: now + time.Hour.Nanoseconds(),
				NotBefore: now,
				IssuedAt:  now,
				ID:        "0",
			},
			Roles: []string{
				"api:access",
				"api:admin",
			},
		}
		validAudience = StringOrURI("tester")
	)
	Context("when validating fields with criteria", func() {
		Context("when the issuer is not valid", func() {
			BeforeEach(func() {
				claim.Issuer = "keanu"
			})
			It("should return an error", func() {
				err := claim.MoreValid(validAudience)
				Ω(err).Should(Equal(ErrIssuer))
			})
		})
		Context("when the audience is not valid", func() {
			BeforeEach(func() {
				claim.Audience = []StringOrURI{StringOrURI("keanu")}
			})
			It("should return an error", func() {
				err := claim.MoreValid(validAudience)
				Ω(err).Should(Equal(ErrAudience))
			})
		})
		Context("when the roles are not valid", func() {
			BeforeEach(func() {
				claim.Roles = []string{"action hero", "the one"}
			})
			It("should return an error", func() {
				err := claim.Valid()
				Ω(err).Should(Equal(ErrRoles))
			})
		})
		Context("when the claims are just right", func() {
			It("should validate with no errors", func() {
				err := claim.Valid()
				Ω(err).Should(BeNil())
				err = claim.MoreValid(validAudience)
				Ω(err).Should(BeNil())
			})
		})
	})
})
