package claims_test

import (
	"time"

	. "github.com/zenoss/zenkit/claims"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func newCompleteClaims() CompleteClaims {
	now := time.Now()
	return CompleteClaims{
		StandardClaims: StandardClaims{
			Iss: CompleteIssuer,
			Sub: "abcd",
			Aud: []string{"tester"},
			Exp: now.Add(ValidDuration).Unix(),
			Nbf: now.Unix(),
			Iat: now.Unix(),
			Jti: "0",
		},
		Tkn: "x",
	}
}

func newCompleteClaimsMap() CompleteClaimsMap {
	now := time.Now()
	return CompleteClaimsMap{
		"iss": CompleteIssuer,
		"sub": "abcd",
		"aud": []string{"tester"},
		"exp": now.Add(ValidDuration).Unix(),
		"nbf": now.Unix(),
		"iat": now.Unix(),
		"jti": "0",
		"tkn": "x",
	}
}

var _ = Describe("Complete Claims", func() {
	var (
		claims        = newCompleteClaims()
		claimsMap     = newCompleteClaimsMap()
		validAudience = "tester"
	)
	BeforeEach(func() {
		claims = newCompleteClaims()
		claimsMap = newCompleteClaimsMap()
	})
	Context("when creating an CompleteClaimsMap", func() {
		Context("from new", func() {
			It("should contain all complete fields", func() {
				m := NewCompleteClaimsMap()
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
				_, ok = m["tkn"]
				Ω(ok).Should(BeTrue())
			})
		})
		Context("from a valid CompleteClaims struct", func() {
			It("should be valid", func() {
				err := claims.Validate(validAudience)
				Ω(err).Should(BeNil())
				m := CompleteClaimsFromStruct(claims)
				err = m.Validate(validAudience)
				Ω(err).Should(BeNil())
			})
		})
	})
	Context("when creating a CompleteClaims struct", func() {
		Context("from a valid CompleteClaimsMap", func() {
			It("should be valid", func() {
				err := claimsMap.Validate(validAudience)
				Ω(err).Should(BeNil())
				m := CompleteClaimsFromMap(claimsMap)
				err = m.Validate(validAudience)
				Ω(err).Should(BeNil())
				Ω(m.Token()).Should(Equal("x"))
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
		Context("when using CompleteClaimsMap", func() {
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
		Context("when using CompleteClaims", func() {
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
	Context("when validating token in a map", func() {
		Context("when the tkn field is empty", func() {
			m := CompleteClaimsMap{}
			It("should return an empty string", func() {
				tkn := m.Token()
				Ω(tkn).Should(Equal(""))
			})
		})
		Context("when the tkn field is not a string", func() {
			m := CompleteClaimsMap{
				"tkn": claims,
			}
			It("should return an empty string", func() {
				tkn := m.Token()
				Ω(tkn).Should(Equal(""))
			})
		})
	})
})
