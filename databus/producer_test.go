package databus_test

import (
	"errors"

	"github.com/Shopify/sarama"
	. "github.com/zenoss/zenkit/databus"
	"github.com/zenoss/zenkit/test"

	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"math/rand"
	"strconv"
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
		producer        *mockSyncProducer
		databusProducer DatabusProducer
		topic           string
		key             string
		value           int

		schemas = map[string]string{
			"object-key":   `"string"`,
			"object-value": `"int"`,
		}

		ids = map[string]int{
			"object-key":   1,
			"object-value": 2}
	)

	BeforeEach(func() {
		topic = test.RandString(8)
		key = test.RandString(8)
		value = rand.Intn(100)

		producer = &mockSyncProducer{messages: make([]*sarama.ProducerMessage, 0)}
		client := GetSchemaRegistryMockClient(schemas, ids)
		messageFactory, err := NewMessageFactory(topic, "object-key", "object-value", client)
		Ω(err).ShouldNot(HaveOccurred())
		databusProducer = NewSaramaDatabusProducer(producer, messageFactory)
	})

	It("should send a message through a Kafka producer", func() {
		err := databusProducer.Send(key, value)
		Ω(err).ShouldNot(HaveOccurred())
		Ω(producer.messages).Should(HaveLen(1))
		m := producer.messages[0]

		Ω(m.Topic).Should(Equal(topic))

		encodedKey, _ := m.Key.Encode()
		Ω(encodedKey).Should(Equal([]byte(fmt.Sprintf(`"%s"`, key))))

		encodedVal, _ := m.Value.Encode()
		Ω(encodedVal).Should(Equal([]byte(strconv.Itoa(value))))
	})

	It("should error if the underlying Kafka producer fails to send the message", func() {
		e := errors.New("NOPE")
		producer.err = e
		err := databusProducer.Send(key, value)
		Ω(err).Should(HaveOccurred())
	})

	It("should close the underlying Kafka producer", func() {
		err := databusProducer.Close()
		Ω(err).ShouldNot(HaveOccurred())
		Ω(producer.closed).Should(BeTrue())
	})

	It("should error when closing if the underlying Kafka producer does so", func() {
		e := errors.New("NOPE")
		producer.err = e
		err := databusProducer.Close()
		Ω(err).Should(HaveOccurred())
	})

})
