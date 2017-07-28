package databus_test

import (
	"errors"

	schemaregistry "github.com/datamountaineer/schema-registry"
	. "github.com/zenoss/zenkit/databus"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Schema", func() {

	It("should create a codec from a valid schema", func() {
		client := &schemaregistry.MockClient{
			GetLatestSchemaFn: func(subject string) (schemaregistry.Schema, error) {
				return schemaregistry.Schema{
					Id:      1,
					Schema:  `"string"`,
					Subject: "domain-object",
					Version: 1,
				}, nil
			},
		}
		id, codec, err := GetCodec(client, "domain-object")
		Ω(err).ShouldNot(HaveOccurred())
		Ω(id).Should(Equal(1))
		Ω(codec.Schema()).Should(Equal(`"string"`))
	})

	It("should fail to create a codec from an invalid schema", func() {
		client := &schemaregistry.MockClient{
			GetLatestSchemaFn: func(subject string) (schemaregistry.Schema, error) {
				return schemaregistry.Schema{
					Id:      1,
					Schema:  `"invalid"`,
					Subject: "domain-object",
					Version: 1,
				}, nil
			},
		}
		_, _, err := GetCodec(client, "domain-object")
		Ω(err).Should(HaveOccurred())
	})

	It("should fail to create a codec from an unregistered subject", func() {
		client := &schemaregistry.MockClient{
			GetLatestSchemaFn: func(subject string) (schemaregistry.Schema, error) {
				return schemaregistry.Schema{}, errors.New("no schema")
			},
		}
		_, _, err := GetCodec(client, "domain-object")
		Ω(err).Should(HaveOccurred())
	})

})
