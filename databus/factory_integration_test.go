// +build integration

package databus_test

import (
	"fmt"

	schemaregistry "github.com/datamountaineer/schema-registry"
	. "github.com/zenoss/zenkit/databus"
	"github.com/zenoss/zenkit/test"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Factory integration", func() {
	var (
		client     schemaregistry.Client
		keySubject = "object-key"
		valSubject = "object-value"
		keySchema  = `"string"`
		valSchema  = `"int"`

		factory MessageFactory
		topic   string
	)

	BeforeEach(func() {
		addr, err := harness.Resolve("kafka-schema-registry", 8081)
		Ω(err).ShouldNot(HaveOccurred())

		client, err = schemaregistry.NewClient(fmt.Sprintf("http://%s/", addr))
		Ω(err).ShouldNot(HaveOccurred())

		_, err = client.RegisterNewSchema(keySubject, fmt.Sprintf(`{"type": %s}`, keySchema))
		Ω(err).ShouldNot(HaveOccurred())

		_, err = client.RegisterNewSchema(valSubject, fmt.Sprintf(`{"type": %s}`, valSchema))
		Ω(err).ShouldNot(HaveOccurred())

		topic = test.RandString(8)
		factory, err = NewMessageFactory(topic, keySubject, valSubject, client)
	})

	It("should fail if given an invalid key subject", func() {
		_, err := NewMessageFactory(topic, "nothing", valSubject, client)
		Ω(err).Should(HaveOccurred())
	})

	It("should fail if given an invalid value subject", func() {
		_, err := NewMessageFactory(topic, keySubject, "nothing", client)
		Ω(err).Should(HaveOccurred())
	})

	It("should return the codec for the key subject provided", func() {
		Ω(factory.KeyCodec().Schema()).Should(Equal(keySchema))
	})

	It("should return the codec for the value subject provided", func() {
		Ω(factory.ValueCodec().Schema()).Should(Equal(valSchema))
	})

})
