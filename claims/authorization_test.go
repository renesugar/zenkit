package claims_test

import (
	"time"

	. "github.com/zenoss/zenkit/claims"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func newAuthorizationClaims() AuthorizationClaims {
	now := time.Now()
	return AuthorizationClaims{
		StandardClaims: StandardClaims{
			Iss: AuthorizationIssuer,
			Sub: "abcd",
			Aud: []string{"tester"},
			Exp: now.Add(ValidDuration).Unix(),
			Nbf: now.Unix(),
			Iat: now.Unix(),
			Jti: "0",
		},
		Rls: []string{"api:access"},
	}
}

func newAuthorizationClaimsMap() AuthorizationClaimsMap {
	now := time.Now()
	return AuthorizationClaimsMap{
		"iss": AuthorizationIssuer,
		"sub": "abcd",
		"aud": []string{"tester"},
		"exp": now.Add(ValidDuration).Unix(),
		"nbf": now.Unix(),
		"iat": now.Unix(),
		"jti": "0",
		"rls": []string{"api:access"},
	}
}

var _ = Describe("Authorization Claims", func() {
	var (
		claims        = newAuthorizationClaims()
		claimsMap     = newAuthorizationClaimsMap()
		validAudience = "tester"
	)
	BeforeEach(func() {
		claims = newAuthorizationClaims()
		claimsMap = newAuthorizationClaimsMap()
	})
	Context("when creating an AuthorizationClaimsMap", func() {
		Context("from new", func() {
			It("should contain all authorization fields", func() {
				m := NewAuthorizationClaimsMap()
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
				_, ok = m["rls"]
				Ω(ok).Should(BeTrue())
			})
		})
		Context("from a valid AuthorizationClaims struct", func() {
			It("should be valid", func() {
				err := claims.Validate(validAudience)
				Ω(err).Should(BeNil())
				m := AuthorizationClaimsFromStruct(claims)
				err = m.Validate(validAudience)
				Ω(err).Should(BeNil())
			})
		})
	})
	Context("when creating a AuthorizationClaims struct", func() {
		Context("from a valid AuthorizationClaimsMap", func() {
			It("should be valid", func() {
				err := claimsMap.Validate(validAudience)
				Ω(err).Should(BeNil())
				m := AuthorizationClaimsFromMap(claimsMap)
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
		Context("when using AuthorizationClaimsMap", func() {
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
		Context("when using AuthorizationClaims", func() {
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
			Context("when the roles are not vlaid", func() {
				BeforeEach(func() {
					claims.Rls = []string{"cheesehoarder"}
				})
				It("should return an error", func() {
					err := claims.Validate(validAudience)
					Ω(err).Should(Equal(ErrRoles))
				})
			})
		})
	})
	Context("when validating roles in a map", func() {
		Context("when the rls field is empty", func() {
			m := AuthorizationClaimsMap{}
			It("should return an empty string slice", func() {
				rls := m.Roles()
				Ω(rls).Should(Equal([]string{}))
			})
		})
		Context("when the aud field is not a slice", func() {
			m := AuthorizationClaimsMap{
				"rls": claims,
			}
			It("should return an empty string slice", func() {
				rls := m.Roles()
				Ω(rls).Should(Equal([]string{}))
			})
		})
		Context("when the aud field is a slice of non-strings", func() {
			Context("when the non-strings cannot be asserted to string", func() {
				m := AuthorizationClaimsMap{
					"rls": []interface{}{claims},
				}
				It("should return an empty string slice", func() {
					rls := m.Roles()
					Ω(rls).Should(Equal([]string{}))
				})
			})
			Context("when the non-strings can be asserted to string", func() {
				m := AuthorizationClaimsMap{
					"rls": []interface{}{"wee"},
				}
				It("should return the roles", func() {
					rls := m.Roles()
					Ω(rls).Should(Equal([]string{"wee"}))
				})
			})
		})
	})
})
