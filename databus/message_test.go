package databus_test

import (
	. "github.com/zenoss/zenkit/databus"
	"github.com/zenoss/zenkit/test"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Message", func() {

	var (
		message Message
		topic   string
		key     []byte
		value   []byte
	)

	BeforeEach(func() {
		topic = test.RandString(8)
		key = []byte(test.RandString(8))
		value = []byte(test.RandString(8))
		message = NewMessage(topic, key, value)
	})

	It("should return the topic it was created with", func() {
		Ω(message.Topic()).Should(Equal(topic))
	})

	It("should return the key it was created with", func() {
		Ω(message.Key()).Should(Equal(key))
	})

	It("should return the value it was created with", func() {
		Ω(message.Value()).Should(Equal(value))
	})

})
