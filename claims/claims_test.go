package claims_test

import (
	"time"

	. "github.com/zenoss/zenkit/claims"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Standard Claims", func() {
	var (
		now           = time.Now().Unix()
		defaultClaims = StandardClaims{
			Issuer:    "test",
			Subject:   "abcd",
			Audience:  []string{"tester"},
			ExpiresAt: now + int64(time.Hour),
			NotBefore: now,
			IssuedAt:  now,
			ID:        "0",
		}
		claims        = defaultClaims
		validIssuers  = []string{"test"}
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
	Context("when verifying all fields have a value", func() {
		Context("when issuer is empty", func() {
			BeforeEach(func() {
				claims.Issuer = ""
			})
			It("should return an error", func() {
				err := claims.Valid()
				Ω(err).Should(Equal(ErrIssuer))
			})
		})
		Context("when subject is empty", func() {
			BeforeEach(func() {
				claims.Subject = ""
			})
			It("should return an error", func() {
				err := claims.Valid()
				Ω(err).Should(Equal(ErrSubject))
			})
		})
		Context("when audience is empty", func() {
			BeforeEach(func() {
				claims.Audience = []string{}
			})
			It("should return an error", func() {
				err := claims.Valid()
				Ω(err).Should(Equal(ErrAudience))
			})
		})
		Context("when expiration is empty", func() {
			BeforeEach(func() {
				claims.ExpiresAt = int64(0)
			})
			It("should return an error", func() {
				err := claims.Valid()
				Ω(err).Should(Equal(ErrExpiresAt))
			})
		})
		Context("when not before is empty", func() {
			BeforeEach(func() {
				claims.NotBefore = int64(0)
			})
			It("should return an error", func() {
				err := claims.Valid()
				Ω(err).Should(Equal(ErrNotBefore))
			})
		})
		Context("when ID is empty", func() {
			BeforeEach(func() {
				claims.ID = ""
			})
			It("should return an error", func() {
				err := claims.Valid()
				Ω(err).Should(Equal(ErrID))
			})
		})
		Context("when no claimed fields are empty", func() {
			It("should validate with no errors", func() {
				err := claims.Valid()
				Ω(err).Should(BeNil())
			})
		})
	})
	Context("when validating fields with criteria", func() {
		Context("when the issuer is not valid", func() {
			BeforeEach(func() {
				claims.Issuer = "keanu"
			})
			It("should return an error", func() {
				err := claims.MoreValid(validIssuers, validAudience)
				Ω(err).Should(Equal(ErrIssuer))
			})
		})
		Context("when the audience is not valid", func() {
			BeforeEach(func() {
				claims.Audience = []string{"keanu"}
			})
			It("should return an error", func() {
				err := claims.MoreValid(validIssuers, validAudience)
				Ω(err).Should(Equal(ErrAudience))
			})
		})
		Context("when the token is expired", func() {
			BeforeEach(func() {
				claims.ExpiresAt = now - 1
			})
			It("should return an error", func() {
				err := claims.Valid()
				Ω(err).Should(Equal(ErrExpiresAt))
			})
		})
		Context("when not before is invalid", func() {
			BeforeEach(func() {
				claims.NotBefore = now + int64(time.Hour)
			})
			It("should return an error", func() {
				err := claims.Valid()
				Ω(err).Should(Equal(ErrNotBefore))
			})
		})
		Context("when issued at is in the future", func() {
			BeforeEach(func() {
				claims.IssuedAt = now + int64(time.Hour)
			})
			It("should return an error", func() {
				err := claims.Valid()
				Ω(err).Should(Equal(ErrIssuedAt))
			})
		})
		Context("when the claims are just right", func() {
			It("should validate with no errors", func() {
				err := claims.Valid()
				Ω(err).Should(BeNil())
				err = claims.MoreValid(validIssuers, validAudience)
				Ω(err).Should(BeNil())
			})
		})
	})
})
