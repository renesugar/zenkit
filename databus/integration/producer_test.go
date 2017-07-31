package databus_test

import (
	"github.com/zenoss/zenkit/test"

	"github.com/Shopify/sarama"
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
		partConsumer, err := testConsumer.ConsumePartition(topic, 0, sarama.OffsetNewest)
		Ω(err).ShouldNot(HaveOccurred())

		defer partConsumer.Close()

		err = msgProducer.Send(msg)
		Ω(err).ShouldNot(HaveOccurred())

		var saramaMessage *sarama.ConsumerMessage
		Eventually(partConsumer.Messages()).Should(Receive(&saramaMessage))

		Ω(saramaMessage.Key).Should(Equal(key))

		Ω(saramaMessage.Value).Should(Equal(value))
	})

})
