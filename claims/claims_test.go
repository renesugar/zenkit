package claims_test

import (
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/zenoss/zenkit/claims"
)

type multiTenantClaims struct {
	StandardClaims
	Tnt string
}

func (c *multiTenantClaims) Tenant() string {
	return c.Tnt
}

var _ = Describe("MultiTenant Claims", func() {
	var (
		claims *multiTenantClaims
	)
	BeforeEach(func() {
		claims = &multiTenantClaims{
			NewStandardClaims("someone", "abcd", []string{"tester"}, time.Hour),
			"example",
		}
	})

	Context("when validating multi-tenant claims", func() {
		Context("with a tenant", func() {
			It("should validate successfully", func() {
				err := ValidateMultiTenantClaims(claims)
				Ω(err).Should(BeNil())
			})
		})
		Context("with no subject", func() {
			BeforeEach(func() {
				claims.Sub = ""
			})
			It("should fail to validate", func() {
				err := ValidateMultiTenantClaims(claims)
				Ω(err).Should(Equal(ErrSubject))
			})
		})
		Context("with no tenant", func() {
			BeforeEach(func() {
				claims.Tnt = ""
			})
			It("should fail to validate", func() {
				err := ValidateMultiTenantClaims(claims)
				Ω(err).Should(Equal(ErrTenant))
			})
		})
	})
})
