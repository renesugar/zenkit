package zenkit_test

import (
	"math/rand"

	. "github.com/zenoss/zenkit"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func level1(level int) string {
	return FuncName(level)
}

func level2(level int) string {
	return level1(level)
}

func level3(level int) string {
	return level2(level)
}

var _ = Describe("Function name utility", func() {

	Context("cache", func() {

		var (
			cache   FnNameCache
			seen    bool
			ptr     uintptr
			val     string
			factory = func() string {
				seen = true
				return val
			}
		)

		BeforeEach(func() {
			cache = NewFnNameCache()
			ptr = uintptr(rand.Intn(1000))
			val = RandStringRunes(8)
		})

		It("should use the factory the first time", func() {
			data, found := cache.Get(ptr, factory)
			Ω(data).Should(Equal(val))
			Ω(found).Should(BeFalse())
			Ω(seen).Should(BeTrue())
		})

		It("should pull from cache the second time", func() {
			data, found := cache.Get(ptr, factory)
			// Reset the seen
			seen = false
			data, found = cache.Get(ptr, factory)
			Ω(data).Should(Equal(val))
			Ω(found).Should(BeTrue())
			Ω(seen).Should(BeFalse())
		})
	})

	Context("name getter", func() {

		It("should return the immediate func at level 1", func() {
			Ω(level1(1)).Should(Equal("level1"))
		})

		It("should return the appropriate functions at level 2", func() {
			Ω(level2(1)).Should(Equal("level1"))
			Ω(level2(2)).Should(Equal("level2"))
		})

		It("should return the appropriate functions at level 3", func() {
			Ω(level3(1)).Should(Equal("level1"))
			Ω(level3(2)).Should(Equal("level2"))
			Ω(level3(3)).Should(Equal("level3"))
		})

	})
})
