package health_test

import (
	"errors"

	. "github.com/zenoss/zenkit/health"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type mockStatusChecker struct {
	status Status
	err    error
}

func (sc *mockStatusChecker) CheckStatus() (Status, error) {
	return sc.status, sc.err
}

var _ = Describe("Health", func() {

	AfterEach(func() {
		Reset()
	})

	It("should return an empty list if no health checks are registered", func() {
		result := Execute()
		Ω(result).Should(BeEmpty())
	})

	It("should return the status of a registered health check", func() {
		name := "test-status"
		checker := &mockStatusChecker{
			status: DEGRADED,
			err:    errors.New("it broke"),
		}
		Register(name, checker)
		result := Execute()
		Ω(result).Should(HaveLen(1))
		Ω(result[0].Name).Should(Equal(name))
		Ω(result[0].Status).Should(Equal(checker.status))
		Ω(result[0].Err).Should(Equal(checker.err))
	})
})
