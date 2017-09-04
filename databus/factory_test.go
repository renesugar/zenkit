package databus_test

import (
	"encoding/json"
	"math/rand"

	schemaregistry "github.com/datamountaineer/schema-registry"
	. "github.com/zenoss/zenkit/databus"
	"github.com/zenoss/zenkit/test"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type KeyTest struct {
	SomeString string
	AnInt      int
}

var keyTestSchema = `{
	"type": "record",
	"name": "keyTest",
	"fields": [{
		"type": "string",
		"name": "SomeString"
	},{
		"type": "int",
		"name": "AnInt"
	}]
}`

type ValTest struct {
	TotallyCool string
}

var valTestSchema = `{
	"type": "record",
	"name": "valTest",
	"fields": [{
		"type": "string",
		"name": "TotallyCool"
	}]
}`

func stripAvroHeader(b []byte) []byte {
	_, k, err := AvroDeserialize(b)
	Ω(err).ShouldNot(HaveOccurred())
	return k
}

var _ = Describe("Factory", func() {

	var (
		factory    MessageFactory
		topic      string
		client     schemaregistry.Client
		keySubject = "object-key"
		valSubject = "object-value"
		schemas    = map[string]string{
			"object-key":   `"string"`,
			"object-value": `"int"`,
			"key-test":     keyTestSchema,
			"val-test":     valTestSchema,
		}

		ids = map[string]int{
			"object-key":   1,
			"object-value": 2,
			"key-test":     3,
			"val-test":     4,
		}
	)

	JustBeforeEach(func() {
		var err error
		client = GetSchemaRegistryMockClient(schemas, ids)
		topic = test.RandString(8)
		factory, err = NewMessageFactory(topic, keySubject, valSubject, client)
		Ω(err).ShouldNot(HaveOccurred())
	})

	It("should fail if given an invalid key subject", func() {
		_, err := NewMessageFactory(topic, "nothing", valSubject, client)
		Ω(err).Should(HaveOccurred())
	})

	It("should fail if given an invalid value subject", func() {
		_, err := NewMessageFactory(topic, keySubject, "nothing", client)
		Ω(err).Should(HaveOccurred())
	})

	It("should return the topic it was created with", func() {
		Ω(factory.Topic()).Should(Equal(topic))
	})

	It("should return the key subject it was created with", func() {
		Ω(factory.KeySubject()).Should(Equal(keySubject))
	})

	It("should return the value subject it was created with", func() {
		Ω(factory.ValueSubject()).Should(Equal(valSubject))
	})

	It("should return the codec for the key subject provided", func() {
		Ω(factory.KeyCodec().Schema()).Should(Equal(schemas[keySubject]))
	})

	It("should return the codec for the value subject provided", func() {
		Ω(factory.ValueCodec().Schema()).Should(Equal(schemas[valSubject]))
	})

	Context("encoding a message", func() {
		var (
			key   string
			value int
		)
		BeforeEach(func() {
			key = test.RandString(8)
			value = rand.Intn(1000)
		})

		It("should return a message with the key encoded using the right schema", func() {
			msg, err := factory.Message(key, value)
			Ω(err).ShouldNot(HaveOccurred())
			var result string
			err = json.Unmarshal(stripAvroHeader(msg.Key()), &result)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(result).Should(Equal(key))
		})

		It("should return a message with the value encoded using the right schema", func() {
			msg, err := factory.Message(key, value)
			Ω(err).ShouldNot(HaveOccurred())
			var result int
			err = json.Unmarshal(stripAvroHeader(msg.Value()), &result)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(result).Should(Equal(value))
		})

		It("should fail to encode a message if the key doesn't match the schema", func() {
			_, err := factory.Message(1234, value)
			Ω(err).Should(HaveOccurred())
		})

		It("should fail to encode a message if the value doesn't match the schema", func() {
			_, err := factory.Message(key, "jklfds")
			Ω(err).Should(HaveOccurred())
		})
	})

	Context("with complex message structures", func() {
		var (
			key   KeyTest
			value ValTest
		)

		BeforeEach(func() {
			key = KeyTest{
				SomeString: test.RandString(8),
				AnInt:      rand.Intn(100),
			}
			value = ValTest{
				TotallyCool: test.RandString(8),
			}
			keySubject = "key-test"
			valSubject = "val-test"
		})

		It("should be able to encode and decode messages", func() {

			By("encoding a message")
			msg, err := factory.Message(key, value)
			Ω(err).ShouldNot(HaveOccurred())

			var (
				k KeyTest
				v ValTest
			)

			By("decoding a message")
			err = factory.Decode(msg, &k, &v)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(k).Should(BeEquivalentTo(key))
			Ω(v).Should(BeEquivalentTo(value))
		})

		It("should fail to encode unserializable keys", func() {
			_, err := factory.Message(make(chan struct{}), value)
			Ω(err).Should(HaveOccurred())
		})

		It("should fail to encode unserializable values", func() {
			_, err := factory.Message(key, make(chan struct{}))
			Ω(err).Should(HaveOccurred())
		})

		It("should fail to decode bad messages", func() {
			var (
				badk = []byte(test.RandString(8))
				badv = []byte(test.RandString(8))
				k    KeyTest
				v    ValTest
			)
			msg, err := factory.Message(key, value)
			Ω(err).ShouldNot(HaveOccurred())

			By("failing to decode bad keys")
			m := NewMessage(topic, badk, stripAvroHeader(msg.Value()))
			err = factory.Decode(m, &k, &v)
			Ω(err).Should(HaveOccurred())

			By("failing to decode bad values")
			m = NewMessage(topic, stripAvroHeader(msg.Key()), badv)
			err = factory.Decode(m, &k, &v)
			Ω(err).Should(HaveOccurred())
		})

	})

})
