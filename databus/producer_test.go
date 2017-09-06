package databus_test

import (
	"errors"

	"github.com/Shopify/sarama"
	"github.com/linkedin/goavro"
	. "github.com/zenoss/zenkit/databus"
	"github.com/zenoss/zenkit/test"

	"math/rand"

	"github.com/datamountaineer/schema-registry"
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
		producer             *mockSyncProducer
		databusProducer      DatabusProducer
		topic                string
		key                  string
		value                int
		schemas              map[string]string
		ids                  map[string]int
		keyCodec, valCodec   *goavro.Codec
		schemaRegistryClient schemaregistry.Client
	)

	BeforeEach(func() {
		topic = test.RandString(8)
		key = test.RandString(8)
		value = rand.Intn(100)

		schemas = map[string]string{
			"object-key":   `"string"`,
			"object-value": `"int"`,
		}

		ids = map[string]int{
			"object-key":   1,
			"object-value": 2,
		}

		keyCodec, _ = goavro.NewCodec(schemas["object-key"])
		valCodec, _ = goavro.NewCodec(schemas["object-value"])

		producer = &mockSyncProducer{messages: make([]*sarama.ProducerMessage, 0)}
		schemaRegistryClient = GetSchemaRegistryMockClient(schemas, ids)
		messageFactory, err := NewMessageFactory(topic, "object-key", "object-value", schemaRegistryClient)
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
		result, _, err := keyCodec.NativeFromBinary(stripAvroHeader(encodedKey))
		Ω(err).ShouldNot(HaveOccurred())
		Ω(result).Should(Equal(key))

		encodedVal, _ := m.Value.Encode()
		result, _, err = valCodec.NativeFromBinary(stripAvroHeader(encodedVal))
		Ω(result).Should(BeNumerically("==", value))
	})

	It("should error if the underlying Kafka producer fails to send the message", func() {
		e := errors.New("NOPE")
		producer.err = e
		err := databusProducer.Send(key, value)
		Ω(err).Should(HaveOccurred())
	})

	It("should error if the message factory errors", func() {
		err := databusProducer.Send(1, value)
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
