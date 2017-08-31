package claims_test

import (
	"time"

	. "github.com/zenoss/zenkit/claims"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func newAuthZeroClaims() AuthZeroClaims {
	now := time.Now()
	return AuthZeroClaims{
		StandardClaims: StandardClaims{
			Iss: AuthZeroIssuer,
			Sub: "abcd",
			Aud: []string{"tester"},
			Exp: now.Add(ValidDuration).Unix(),
			Nbf: now.Unix(),
			Iat: now.Unix(),
			Jti: "0",
		},
	}
}

func newAuthZeroClaimsMap() AuthZeroClaimsMap {
	now := time.Now()
	return AuthZeroClaimsMap{
		"iss": AuthZeroIssuer,
		"sub": "abcd",
		"aud": []string{"tester"},
		"exp": now.Add(ValidDuration).Unix(),
		"nbf": now.Unix(),
		"iat": now.Unix(),
		"jti": "0",
	}
}

var _ = Describe("AuthZero Claims", func() {
	var (
		claims        = newAuthZeroClaims()
		claimsMap     = newAuthZeroClaimsMap()
		validAudience = "tester"
	)
	BeforeEach(func() {
		claims = newAuthZeroClaims()
		claimsMap = newAuthZeroClaimsMap()
	})
	Context("when creating an AuthZeroClaimsMap", func() {
		Context("from new", func() {
			It("should contain all authzero fields", func() {
				m := NewAuthZeroClaimsMap()
				_, ok := m["iss"]
				Ω(ok).Should(BeTrue())
				_, ok = m["sub"]
				Ω(ok).Should(BeTrue())
				_, ok = m["aud"]
				Ω(ok).Should(BeTrue())
				_, ok = m["exp"]
				Ω(ok).Should(BeTrue())
				_, ok = m["nbf"]
				Ω(ok).Should(BeTrue())
				_, ok = m["iat"]
				Ω(ok).Should(BeTrue())
				_, ok = m["jti"]
				Ω(ok).Should(BeTrue())
			})
		})
		Context("from a valid AuthZeroClaims struct", func() {
			It("should be valid", func() {
				err := claims.Validate(validAudience)
				Ω(err).Should(BeNil())
				m := AuthZeroClaimsFromStruct(claims)
				err = m.Validate(validAudience)
				Ω(err).Should(BeNil())
			})
		})
	})
	Context("when creating a AuthZeroClaims struct", func() {
		Context("from a valid AuthZeroClaimsMap", func() {
			It("should be valid", func() {
				err := claimsMap.Validate(validAudience)
				Ω(err).Should(BeNil())
				m := AuthZeroClaimsFromMap(claimsMap)
				err = m.Validate(validAudience)
				Ω(err).Should(BeNil())
			})
		})
	})
	Context("when standard claims do not validate", func() {
		Context("as a struct", func() {
			BeforeEach(func() {
				claims.Sub = ""
			})
			It("should return an error", func() {
				err := claims.Validate(validAudience)
				Ω(err).ShouldNot(BeNil())
			})
		})
		Context("as a map", func() {
			BeforeEach(func() {
				claimsMap["sub"] = ""
			})
			It("should return an error", func() {
				err := claimsMap.Validate(validAudience)
				Ω(err).ShouldNot(BeNil())
			})
		})
	})
	Context("when validating fields with criteria", func() {
		Context("when using AuthZeroClaimsMap", func() {
			Context("when the issuer is not valid", func() {
				BeforeEach(func() {
					claimsMap["iss"] = "keanu"
				})
				It("should return an error", func() {
					err := claimsMap.Validate(validAudience)
					Ω(err).Should(Equal(ErrIssuer))
				})
			})
			Context("when the subject does not exist", func() {
				BeforeEach(func() {
					claimsMap["sub"] = ""
				})
				It("should return an error", func() {
					err := claimsMap.Valid()
					Ω(err).Should(Equal(ErrSubject))
				})
			})
			Context("when the claims are just right", func() {
				It("should validate with no errors", func() {
					err := claimsMap.Validate(validAudience)
					Ω(err).Should(BeNil())
				})
			})
		})
		Context("when using AuthZeroClaims", func() {
			Context("when the issuer is not valid", func() {
				BeforeEach(func() {
					claims.Iss = "keanu"
				})
				It("should return an error", func() {
					err := claims.Validate(validAudience)
					Ω(err).Should(Equal(ErrIssuer))
				})
			})
			Context("when the audience is not valid", func() {
				BeforeEach(func() {
					claims.Aud = []string{"keanu"}
				})
				It("should return an error", func() {
					err := claims.Validate(validAudience)
					Ω(err).Should(Equal(ErrAudience))
				})
			})
			Context("when the claims are just right", func() {
				It("should validate with no errors", func() {
					err := claims.Validate(validAudience)
					Ω(err).Should(BeNil())
				})
			})
		})
	})
})
