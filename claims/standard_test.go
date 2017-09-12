package claims_test

import (
	"time"

	. "github.com/zenoss/zenkit/claims"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func newStandardClaimsMap() StandardClaimsMap {
	now := time.Now()
	return StandardClaimsMap{
		"iss": "someone",
		"sub": "abcd",
		"aud": []string{"tester"},
		"exp": now.Add(ValidDuration).Unix(),
		"nbf": now.Unix(),
		"iat": now.Unix(),
		"jti": "0",
	}
}

var _ = Describe("Standard Claims", func() {
	var (
		now           = time.Now().Unix()
		claims        = NewStandardClaims("someone", "abcd", []string{"tester"})
		claimsMap     = newStandardClaimsMap()
		validIssuers  = []string{"someone"}
		validAudience = "tester"
	)
	BeforeEach(func() {
		claims = NewStandardClaims("someone", "abcd", []string{"tester"})
		claimsMap = newStandardClaimsMap()
	})
	Context("when creating a StandardClaimsMap", func() {
		Context("from new", func() {
			It("should contain all standard fields", func() {
				m := NewStandardClaimsMap()
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
		Context("from a valid StandardClaims struct", func() {
			It("should be valid", func() {
				err := claims.Validate(validIssuers, validAudience)
				Ω(err).Should(BeNil())
				m := StandardClaimsFromStruct(claims)
				err = m.Validate(validIssuers, validAudience)
				Ω(err).Should(BeNil())
			})
		})
	})
	Context("when creating a StandardClaims struct", func() {
		Context("from a valid StandardClaimsMap", func() {
			It("should be valid", func() {
				err := claimsMap.Validate(validIssuers, validAudience)
				Ω(err).Should(BeNil())
				m := StandardClaimsFromMap(claimsMap)
				err = m.Validate(validIssuers, validAudience)
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
				err := claims.Validate(validIssuers, validAudience)
				Ω(err).ShouldNot(BeNil())
			})
		})
		Context("as a map", func() {
			BeforeEach(func() {
				claimsMap["sub"] = ""
			})
			It("should return an error", func() {
				err := claimsMap.Validate(validIssuers, validAudience)
				Ω(err).ShouldNot(BeNil())
			})
		})
	})
	Context("when verifying all fields have a value", func() {
		Context("when issuer is empty", func() {
			BeforeEach(func() {
				claims.Iss = ""
			})
			It("should return an error", func() {
				err := claims.Valid()
				Ω(err).Should(Equal(ErrIssuer))
			})
		})
		Context("when subject is empty", func() {
			BeforeEach(func() {
				claims.Sub = ""
			})
			It("should return an error", func() {
				err := claims.Valid()
				Ω(err).Should(Equal(ErrSubject))
			})
		})
		Context("when audience is empty", func() {
			BeforeEach(func() {
				claims.Aud = []string{}
			})
			It("should return an error", func() {
				err := claims.Valid()
				Ω(err).Should(Equal(ErrAudience))
			})
		})
		Context("when expiration is empty", func() {
			BeforeEach(func() {
				claims.Exp = int64(0)
			})
			It("should return an error", func() {
				err := claims.Valid()
				Ω(err).Should(Equal(ErrExpiresAt))
			})
		})
		Context("when not before is empty", func() {
			BeforeEach(func() {
				claims.Nbf = int64(0)
			})
			It("should return an error", func() {
				err := claims.Valid()
				Ω(err).Should(Equal(ErrNotBefore))
			})
		})
		Context("when ID is empty", func() {
			BeforeEach(func() {
				claims.Jti = ""
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
		Context("when using StandardClaimsMap", func() {
			Context("when the issuer is not valid", func() {
				BeforeEach(func() {
					claimsMap["iss"] = "keanu"
				})
				It("should return an error", func() {
					err := claimsMap.Validate(validIssuers, validAudience)
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
					err := claimsMap.Validate(validIssuers, validAudience)
					Ω(err).Should(BeNil())
				})
			})
		})
		Context("when using StandardClaims", func() {
			Context("when the issuer is not valid", func() {
				BeforeEach(func() {
					claims.Iss = "keanu"
				})
				It("should return an error", func() {
					err := claims.Validate(validIssuers, validAudience)
					Ω(err).Should(Equal(ErrIssuer))
				})
			})
			Context("when the token is expired", func() {
				BeforeEach(func() {
					claims.Exp = now - 1
				})
				It("should return an error", func() {
					err := claims.Valid()
					Ω(err).Should(Equal(ErrExpiresAt))
				})
			})
			Context("when not before is invalid", func() {
				BeforeEach(func() {
					claims.Nbf = now + int64(time.Hour)
				})
				It("should return an error", func() {
					err := claims.Valid()
					Ω(err).Should(Equal(ErrNotBefore))
				})
			})
			Context("when issued at is in the future", func() {
				BeforeEach(func() {
					claims.Iat = now + int64(time.Hour)
				})
				It("should return an error", func() {
					err := claims.Valid()
					Ω(err).Should(Equal(ErrIssuedAt))
				})
			})
			Context("when the claims are just right", func() {
				It("should validate with no errors", func() {
					err := claims.Validate(validIssuers, validAudience)
					Ω(err).Should(BeNil())
				})
			})
		})
	})
	Context("when using an unorthodox StandardClaimsMap", func() {
		Context("when checking the issuer", func() {
			Context("when the iss field is empty", func() {
				m := StandardClaimsMap{}
				It("should return an empty string", func() {
					iss := m.Issuer()
					Ω(iss).Should(Equal(""))
				})
			})
			Context("when the iss field is not a string", func() {
				m := StandardClaimsMap{
					"iss": claims,
				}
				It("should return an empty string", func() {
					iss := m.Issuer()
					Ω(iss).Should(Equal(""))
				})
			})
		})
		Context("when checking the subject", func() {
			Context("when the sub field is empty", func() {
				m := StandardClaimsMap{}
				It("should return an empty string", func() {
					sub := m.Subject()
					Ω(sub).Should(Equal(""))
				})
			})
			Context("when the sub field is not a string", func() {
				m := StandardClaimsMap{
					"sub": claims,
				}
				It("should return an empty string", func() {
					sub := m.Subject()
					Ω(sub).Should(Equal(""))
				})
			})
		})
		Context("when checking the audience", func() {
			Context("when the aud field is empty", func() {
				m := StandardClaimsMap{}
				It("should return an empty string slice", func() {
					aud := m.Audience()
					Ω(aud).Should(Equal([]string{}))
				})
			})
			Context("when the aud field is not a slice", func() {
				m := StandardClaimsMap{
					"aud": claims,
				}
				It("should return an empty string slice", func() {
					aud := m.Audience()
					Ω(aud).Should(Equal([]string{}))
				})
			})
			Context("when the aud field is a slice of non-strings", func() {
				Context("when the non-strings cannot be asserted to string", func() {
					m := StandardClaimsMap{
						"aud": []interface{}{claims},
					}
					It("should return an empty string slice", func() {
						aud := m.Audience()
						Ω(aud).Should(Equal([]string{}))
					})
				})
				Context("when the non-strings can be asserted to string", func() {
					m := StandardClaimsMap{
						"aud": []interface{}{"wee"},
					}
					It("should return the audience", func() {
						aud := m.Audience()
						Ω(aud).Should(Equal([]string{"wee"}))
					})
				})
			})
		})
		Context("when checking the expiration", func() {
			Context("when the exp field is empty", func() {
				m := StandardClaimsMap{}
				It("should return 0", func() {
					exp := m.ExpiresAt()
					Ω(exp).Should(Equal(int64(0)))
				})
			})
			Context("when the exp field is not an int64", func() {
				Context("when it can be converted to a float64", func() {
					m := StandardClaimsMap{
						"exp": 2e03,
					}
					It("should return that number as an int64", func() {
						exp := m.ExpiresAt()
						Ω(exp).Should(Equal(int64(2000)))
					})
				})
				Context("when it cannot be converted to a float64", func() {
					m := StandardClaimsMap{
						"exp": claims,
					}
					It("should return 0", func() {
						exp := m.ExpiresAt()
						Ω(exp).Should(Equal(int64(0)))
					})
				})
			})
		})
		Context("when checking the not before", func() {
			Context("when the nbf field is empty", func() {
				m := StandardClaimsMap{}
				It("should return 0", func() {
					nbf := m.NotBefore()
					Ω(nbf).Should(Equal(int64(0)))
				})
			})
			Context("when the nbf field is not an int64", func() {
				Context("when it can be converted to a float64", func() {
					m := StandardClaimsMap{
						"nbf": 2e03,
					}
					It("should return that number as an int64", func() {
						nbf := m.NotBefore()
						Ω(nbf).Should(Equal(int64(2000)))
					})
				})
				Context("when it cannot be converted to a float64", func() {
					m := StandardClaimsMap{
						"nbf": claims,
					}
					It("should return 0", func() {
						nbf := m.NotBefore()
						Ω(nbf).Should(Equal(int64(0)))
					})
				})
			})
		})
		Context("when checking the issued at", func() {
			Context("when the iat field is empty", func() {
				m := StandardClaimsMap{}
				It("should return 0", func() {
					iat := m.IssuedAt()
					Ω(iat).Should(Equal(int64(0)))
				})
			})
			Context("when the exp field is not an int64", func() {
				Context("when it can be converted to a float64", func() {
					m := StandardClaimsMap{
						"iat": 2e03,
					}
					It("should return that number as an int64", func() {
						iat := m.IssuedAt()
						Ω(iat).Should(Equal(int64(2000)))
					})
				})
				Context("when it cannot be converted to a float64", func() {
					m := StandardClaimsMap{
						"iat": claims,
					}
					It("should return 0", func() {
						iat := m.IssuedAt()
						Ω(iat).Should(Equal(int64(0)))
					})
				})
			})
		})
		Context("when checking the jwt ID", func() {
			Context("when the jti field is empty", func() {
				m := StandardClaimsMap{}
				It("should return an empty string", func() {
					jti := m.ID()
					Ω(jti).Should(Equal(""))
				})
			})
			Context("when the jti field is not a string", func() {
				m := StandardClaimsMap{
					"jti": claims,
				}
				It("should return an empty string", func() {
					jti := m.ID()
					Ω(jti).Should(Equal(""))
				})
			})
		})
	})
})
