package zenkit_test

import (
	"time"

	. "github.com/zenoss/zenkit"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Time", func() {

	It("should default to now", func() {
		Ω(TimeFunc()).Should(BeTemporally("~", time.Now()))
	})

	It("should default to UTC", func() {
		name, offset := TimeFunc().Zone()
		Ω(name).Should(Equal("UTC"))
		Ω(offset).Should(Equal(0))
	})

})
