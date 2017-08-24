package claims_test

import (
	"time"

	. "github.com/zenoss/zenkit/claims"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Edge Claims", func() {
	var (
		now           = time.Now().Unix()
		defaultClaims = EdgeClaims{
			StandardClaims: StandardClaims{
				Issuer:    EdgeIssuer,
				Subject:   StringOrURI("abcd"),
				Audience:  []StringOrURI{StringOrURI("tester")},
				ExpiresAt: now + time.Hour.Nanoseconds(),
				NotBefore: now,
				IssuedAt:  now,
				ID:        "0",
			},
		}
		claims        = defaultClaims
		validAudience = StringOrURI("tester")
	)
	BeforeEach(func() {
		claims = defaultClaims
	})
	Context("when standard claims do not validate", func() {
		BeforeEach(func() {
			claims.Subject = StringOrURI("")
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
				claims.Audience = []StringOrURI{StringOrURI("keanu")}
			})
			It("should return an error", func() {
				err := claims.MoreValid(validAudience)
				Ω(err).Should(Equal(ErrAudience))
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
