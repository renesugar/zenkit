package databus_test

import (
	"errors"

	"github.com/Shopify/sarama"
	. "github.com/zenoss/zenkit/databus"
	"github.com/zenoss/zenkit/test"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type mockSyncProducer struct {
	messages []*sarama.ProducerMessage
	closed   bool
	err      error
}

func (p *mockSyncProducer) SendMessages(msgs []*sarama.ProducerMessage) error {
	// Not implemented
	return nil
}

func (p *mockSyncProducer) Close() error {
	if p.err != nil {
		return p.err
	}
	p.closed = true
	return nil
}

func (p *mockSyncProducer) SendMessage(msg *sarama.ProducerMessage) (int32, int64, error) {
	if p.err != nil {
		return 0, 0, p.err
	}
	p.messages = append(p.messages, msg)
	return 0, 0, nil
}

var _ = Describe("Producer", func() {

	var (
		producer    *mockSyncProducer
		msgProducer MessageProducer
		topic       string
		key         []byte
		value       []byte
		msg         Message
	)

	BeforeEach(func() {
		producer = &mockSyncProducer{messages: make([]*sarama.ProducerMessage, 0)}
		msgProducer = NewMessageProducer(producer)
		topic = test.RandString(8)
		key = []byte(test.RandString(8))
		value = []byte(test.RandString(20))
		msg = NewMessage(topic, key, value)
	})

	It("should send a message through a Kafka producer", func() {
		err := msgProducer.Send(msg)
		Ω(err).ShouldNot(HaveOccurred())
		Ω(producer.messages).Should(HaveLen(1))
		m := producer.messages[0]

		Ω(m.Topic).Should(Equal(topic))

		encodedKey, _ := m.Key.Encode()
		Ω(encodedKey).Should(Equal(key))

		encodedVal, _ := m.Value.Encode()
		Ω(encodedVal).Should(Equal(value))
	})

	It("should error if the underlying Kafka producer fails to send the message", func() {
		e := errors.New("NOPE")
		producer.err = e
		err := msgProducer.Send(msg)
		Ω(err).Should(HaveOccurred())
	})

	It("should close the underlying Kafka producer", func() {
		err := msgProducer.Close()
		Ω(err).ShouldNot(HaveOccurred())
		Ω(producer.closed).Should(BeTrue())
	})

	It("should error when closing if the underlying Kafka producer does so", func() {
		e := errors.New("NOPE")
		producer.err = e
		err := msgProducer.Close()
		Ω(err).Should(HaveOccurred())
	})

})
