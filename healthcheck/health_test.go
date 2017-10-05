package healthcheck_test

import (
	"errors"
	"time"

	. "github.com/zenoss/zenkit/healthcheck"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Health", func() {

	AfterEach(func() {
		DefaultRegistry = NewRegistry()
	})

	Context("with a StatusUpdater", func() {
		var u Updater

		BeforeEach(func() {
			u = NewStatusUpdater()
		})

		It("should return a nil status if status is not set", func() {
			Ω(u.Check()).Should(BeNil())
		})

		It("should return a non-nil status if the status is set as such", func() {
			err := errors.New("he dead")
			u.Update(err)
			Ω(u.Check()).Should(Equal(err))
		})
	})

	Context("with a ThresholdStatusUpdater", func() {
		var (
			u Updater
		)

		BeforeEach(func() {
			u = NewThresholdStatusUpdater(2)
		})

		It("should not change the state if consecutive errors do not exceed the threshold", func() {
			err := errors.New("he dead")

			By("receiving a success")

			u.Update(nil)
			Ω(u.Check()).Should(BeNil())

			By("receieving an error less than the threshold")

			u.Update(err)
			Ω(u.Check()).Should(BeNil())

			By("receiving a success before the error count meets the threshold")

			u.Update(nil)
			Ω(u.Check()).Should(BeNil())

			By("receiving a number of errors that meet or exceed the threshold")

			u.Update(err)
			Ω(u.Check()).Should(BeNil())

			u.Update(err)
			Ω(u.Check()).Should(Equal(err))
		})
	})

	Context("with a PeriodicChecker", func() {

		var (
			u Updater
		)

		BeforeEach(func() {
			u = NewStatusUpdater()
		})

		It("should only update the status on the tick", func() {
			c := PeriodicChecker(u, time.Second)
			Ω(c.Check()).Should(BeNil())
			u.Update(errors.New("he dead"))
			Ω(c.Check()).Should(BeNil())
			Eventually(c.Check, 2*time.Second).ShouldNot(BeNil())
		})
	})

	Context("with a PeriodicThresholdChecker", func() {

		var (
			u Updater
		)

		BeforeEach(func() {
			u = NewStatusUpdater()
		})

		It("should only update the status on the tick and after it meets the threshold", func() {
			c := PeriodicThresholdChecker(u, time.Second, 2)
			Ω(c.Check()).Should(BeNil())
			u.Update(errors.New("he dead"))
			Ω(c.Check()).Should(BeNil())
			Eventually(c.Check, 3*time.Second).ShouldNot(BeNil())
		})
	})

	Context("calling CheckStatus", func() {
		var u Updater

		BeforeEach(func() {
			u = NewStatusUpdater()
			Register("test", u)
		})

		It("should return an empty map if there are no health check errors", func() {
			Ω(CheckStatus()).Should(BeEmpty())
		})

		It("should return a non-empty map if there are health check errors", func() {
			u.Update(errors.New("he dead"))
			m := CheckStatus()
			Ω(m).Should(HaveLen(1))
			Ω(m).Should(HaveKeyWithValue("test", "he dead"))
		})
	})

	Context("registering a health check", func() {

		It("should register a checker", func() {
			u := NewStatusUpdater()
			Register("test", u)

			By("registering another checker with the same name, it should panic")

			Ω(func() { Register("test", u) }).Should(Panic())
		})

		It("should register a function checker", func() {
			f := func() error { return nil }
			RegisterFunc("test", f)
		})

		It("should register a periodic function checker", func() {
			f := func() error { return nil }
			RegisterPeriodicFunc("test", time.Second, f)
		})

		It("should register a periodic threshold function checker", func() {
			f := func() error { return nil }
			RegisterPeriodicThresholdFunc("test", time.Second, 2, f)
		})
	})
})
