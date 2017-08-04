// +build integration

package databus_test

import (
	"fmt"

	schemaregistry "github.com/datamountaineer/schema-registry"
	. "github.com/zenoss/zenkit/databus"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Schema registry integration", func() {

	var (
		client  schemaregistry.Client
		subject = "domain-object"
		schema  = `"string"`
	)

	BeforeEach(func() {
		addr, err := harness.Resolve("kafka-schema-registry", 8081)
		Ω(err).ShouldNot(HaveOccurred())

		client, err = schemaregistry.NewClient(fmt.Sprintf("http://%s/", addr))
		Ω(err).ShouldNot(HaveOccurred())

		_, err = client.RegisterNewSchema(subject, fmt.Sprintf(`{"type": %s}`, schema))
		Ω(err).ShouldNot(HaveOccurred())
	})

	It("should create a codec from a valid schema", func() {
		_, codec, err := GetCodec(client, subject)
		Ω(err).ShouldNot(HaveOccurred())
		Ω(codec.Schema()).Should(Equal(schema))
	})

	It("should fail to create a codec from an unregistered subject", func() {
		_, _, err := GetCodec(client, "not-registered")
		Ω(err).Should(HaveOccurred())
	})
})
