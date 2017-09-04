// +build integration

package databus_test

import (
	"github.com/zenoss/zenkit/test"

	"encoding/json"
	"fmt"
	"math/rand"

	"github.com/Shopify/sarama"
	"github.com/datamountaineer/schema-registry"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/zenoss/zenkit/databus"
)

var _ = Describe("Producer", func() {

	var (
		dbProducer         DatabusProducer
		topic              string
		key                string
		value              int
		keySubject         string
		valueSubject       string
		brokers            []string
		keySchema          = `"string"`
		valueSchema        = `"int"`
		schemaRegistryAddr string
		err                error
	)

	BeforeEach(func() {

		topic = test.RandString(8)
		key = test.RandString(8)
		keySubject = test.RandString(8)
		value = rand.Intn(100)
		valueSubject = test.RandString(8)

		brokersstring, err := harness.Resolve("kafka", 9092)
		Ω(err).ShouldNot(HaveOccurred())
		brokers = []string{brokersstring}

		addr, err := harness.Resolve("kafka-schema-registry", 8081)
		Ω(err).ShouldNot(HaveOccurred())

		schemaRegistryAddr = fmt.Sprintf("http://%s/", addr)
		client, err := schemaregistry.NewClient(schemaRegistryAddr)
		Ω(err).ShouldNot(HaveOccurred())

		_, err = client.RegisterNewSchema(keySubject, fmt.Sprintf(`{"type": %s}`, keySchema))
		Ω(err).ShouldNot(HaveOccurred())

		_, err = client.RegisterNewSchema(valueSubject, fmt.Sprintf(`{"type": %s}`, valueSchema))
		Ω(err).ShouldNot(HaveOccurred())

		err = nil

	})

	JustBeforeEach(func() {
		dbProducer, err = NewDatabusProducer(brokers, schemaRegistryAddr, topic, keySubject, valueSubject)
	})

	It("should send a message through a Kafka producer", func() {
		Ω(err).ShouldNot(HaveOccurred())

		partConsumer, err := testConsumer.ConsumePartition(topic, 0, sarama.OffsetNewest)
		Ω(err).ShouldNot(HaveOccurred())

		defer partConsumer.Close()

		err = dbProducer.Send(key, value)
		Ω(err).ShouldNot(HaveOccurred())

		var saramaMessage *sarama.ConsumerMessage
		Eventually(partConsumer.Messages()).Should(Receive(&saramaMessage))

		expected, _ := json.Marshal(key)
		Ω(saramaMessage.Key).Should(Equal(expected))

		expectedValue, _ := json.Marshal(value)
		Ω(saramaMessage.Value).Should(Equal(expectedValue))
	})

	Context("with an invalid broker", func() {
		BeforeEach(func() {
			brokers = []string{"bad"}
		})
		It("should fail to create the producer", func() {
			Ω(err).Should(HaveOccurred())
		})
	})

	Context("with an invalid key subject", func() {
		BeforeEach(func() {
			keySubject = "bad"
		})
		It("should fail to create the producer", func() {
			Ω(err).Should(HaveOccurred())
		})
	})

	Context("with an invalid value subject", func() {
		BeforeEach(func() {
			valueSubject = "bad"
		})
		It("should fail to create the producer", func() {
			Ω(err).Should(HaveOccurred())
		})
	})

	Context("with an invalid schema registry", func() {
		BeforeEach(func() {
			schemaRegistryAddr = ":bad"
		})
		It("should fail to create the producer", func() {
			Ω(err).Should(HaveOccurred())
		})
	})

})
