package claims_test

import (
	"time"

	. "github.com/zenoss/zenkit/claims"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Standard Claims", func() {
	var (
		now   = time.Now().Unix()
		claim = StandardClaims{
			Issuer:    StringOrURI("test"),
			Subject:   StringOrURI("abcd"),
			Audience:  []StringOrURI{StringOrURI("tester")},
			ExpiresAt: now + time.Hour.Nanoseconds(),
			NotBefore: now,
			IssuedAt:  now,
			ID:        "0",
		}
		validIssuers  = []StringOrURI{StringOrURI("test")}
		validAudience = StringOrURI("tester")
	)
	Context("when verifying all fields have a value", func() {
		Context("when issuer is empty", func() {
			BeforeEach(func() {
				claim.Issuer = StringOrURI("")
			})
			It("should return an error", func() {
				err := claim.Valid()
				Ω(err).Should(Equal(ErrIssuer))
			})
		})
		Context("when subject is empty", func() {
			BeforeEach(func() {
				claim.Subject = StringOrURI("")
			})
			It("should return an error", func() {
				err := claim.Valid()
				Ω(err).Should(Equal(ErrSubject))
			})
		})
		Context("when audience is empty", func() {
			BeforeEach(func() {
				claim.Audience = []StringOrURI{}
			})
			It("should return an error", func() {
				err := claim.Valid()
				Ω(err).Should(Equal(ErrAudience))
			})
		})
		Context("when expiration is empty", func() {
			BeforeEach(func() {
				claim.ExpiresAt = int64(0)
			})
			It("should return an error", func() {
				err := claim.Valid()
				Ω(err).Should(Equal(ErrExpiresAt))
			})
		})
		Context("when not before is empty", func() {
			BeforeEach(func() {
				claim.NotBefore = int64(0)
			})
			It("should return an error", func() {
				err := claim.Valid()
				Ω(err).Should(Equal(ErrNotBefore))
			})
		})
		Context("when ID is empty", func() {
			BeforeEach(func() {
				claim.ID = ""
			})
			It("should return an error", func() {
				err := claim.Valid()
				Ω(err).Should(Equal(ErrID))
			})
		})
		Context("when no claimed fields are empty", func() {
			It("should validate with no errors", func() {
				err := claim.Valid()
				Ω(err).Should(BeNil())
			})
		})
	})
	Context("when validating fields with criteria", func() {
		Context("when the issuer is not valid", func() {
			BeforeEach(func() {
				claim.Issuer = "keanu"
			})
			It("should return an error", func() {
				err := claim.MoreValid(validIssuers, validAudience)
				Ω(err).Should(Equal(ErrIssuer))
			})
		})
		Context("when the audience is not valid", func() {
			BeforeEach(func() {
				claim.Audience = []StringOrURI{StringOrURI("keanu")}
			})
			It("should return an error", func() {
				err := claim.MoreValid(validIssuers, validAudience)
				Ω(err).Should(Equal(ErrAudience))
			})
		})
		Context("when the token is expired", func() {
			BeforeEach(func() {
				claim.ExpiresAt = now - 1
			})
			It("should return an error", func() {
				err := claim.Valid()
				Ω(err).Should(Equal(ErrExpiresAt))
			})
		})
		Context("when not before is invalid", func() {
			BeforeEach(func() {
				claim.NotBefore = now + time.Hour.Nanoseconds()
			})
			It("should return an error", func() {
				err := claim.Valid()
				Ω(err).Should(Equal(ErrNotBefore))
			})
		})
		Context("when issued at is in the future", func() {
			BeforeEach(func() {
				claim.IssuedAt = now + time.Hour.Nanoseconds()
			})
			It("should return an error", func() {
				err := claim.Valid()
				Ω(err).Should(Equal(ErrIssuedAt))
			})
		})
		Context("when the claims are just right", func() {
			It("should validate with no errors", func() {
				err := claim.Valid()
				Ω(err).Should(BeNil())
				err = claim.MoreValid(validIssuers, validAudience)
				Ω(err).Should(BeNil())
			})
		})
	})
})
