package databus_test

import (
	"github.com/zenoss/zenkit/test"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/zenoss/zenkit/databus"
)

var _ = Describe("Producer", func() {

	var (
		msgProducer MessageProducer
		topic       string
		key         []byte
		value       []byte
		msg         Message
	)

	BeforeEach(func() {
		msgProducer = NewMessageProducer(testProducer)
		topic = test.RandString(8)
		key = []byte(test.RandString(8))
		value = []byte(test.RandString(20))
		msg = NewMessage(topic, key, value)
	})

	It("should send a message through a Kafka producer", func() {
		err := msgProducer.Send(msg)
		Ω(err).ShouldNot(HaveOccurred())
	})

})
