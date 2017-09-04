package databus_test

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"strconv"

	"github.com/Shopify/sarama"
	"github.com/bsm/sarama-cluster"
	"github.com/datamountaineer/schema-registry"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pkg/errors"
	. "github.com/zenoss/zenkit/databus"
	"github.com/zenoss/zenkit/test"
)

func newMockClusterConsumer(chanSize int) *mockClusterConsumer {
	return &mockClusterConsumer{
		messages:      make(chan *sarama.ConsumerMessage, chanSize),
		errors:        make(chan error, chanSize),
		notifications: make(chan *cluster.Notification, chanSize),
		err:           nil,
	}
}

type mockClusterConsumer struct {
	messages      chan *sarama.ConsumerMessage
	errors        chan error
	notifications chan *cluster.Notification

	err error
}

func (c *mockClusterConsumer) Messages() <-chan *sarama.ConsumerMessage {
	return c.messages
}

func (c *mockClusterConsumer) Errors() <-chan error {
	return c.errors
}

func (c *mockClusterConsumer) Notifications() <-chan *cluster.Notification {
	return c.notifications
}

func (c *mockClusterConsumer) Close() error {
	if c.err != nil {
		return c.err
	}
	close(c.messages)
	return nil
}

func (c *mockClusterConsumer) MarkOffset(msg *sarama.ConsumerMessage, metadata string) {
	return
}

type TestMessageType struct {
	TestKey   string `zenkit:"message-key"`
	TestValue int    `zenkit:"message-value"`
}

type TestMessageTypeStruct struct {
	TestKey   KeyTest `zenkit:"message-key"`
	TestValue ValTest `zenkit:"message-value"`
}

type TestMessageTypeStructPtr struct {
	TestKey   *KeyTest `zenkit:"message-key"`
	TestValue *ValTest `zenkit:"message-value"`
}

var _ = Describe("Consumer", func() {

	var (
		clusterConsumer      *mockClusterConsumer
		databusConsumer      DatabusConsumer
		topic                string
		schemaRegistryClient schemaregistry.Client
		err                  error
		messageFactory       MessageFactory
		messagesFromKafka    []*sarama.ConsumerMessage
		expectation          []interface{}
		keySchema            string
		valueSchema          string
		schemas              = map[string]string{
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

	BeforeEach(func() {
		topic = test.RandString(8)
		messagesFromKafka = []*sarama.ConsumerMessage{}
		expectation = []interface{}{}
		keySchema = "object-key"
		valueSchema = "object-value"
		schemaRegistryClient = GetSchemaRegistryMockClient(schemas, ids)
		err = nil
	})

	JustBeforeEach(func() {
		messageFactory, err = NewMessageFactory(topic, keySchema, valueSchema, schemaRegistryClient)
		Ω(err).ShouldNot(HaveOccurred())

		clusterConsumer = newMockClusterConsumer(len(messagesFromKafka))
		done := make(chan interface{})
		go func() {
			for _, msg := range messagesFromKafka {
				clusterConsumer.messages <- msg
			}
			close(done)
		}()
		Eventually(done).Should(BeClosed())

		databusConsumer, err = NewSaramaClusterDatabusConsumer(clusterConsumer, messageFactory)
	})

	Context("with primitive keys and values", func() {
		BeforeEach(func() {
			for i := 0; i < 5; i++ {
				key := test.RandString(8)
				value := rand.Intn(100)
				keybin, err := AvroSerialize([]byte(fmt.Sprintf(`"%s"`, key)), 1)
				Ω(err).ShouldNot(HaveOccurred())
				valbin, err := AvroSerialize([]byte(strconv.Itoa(value)), 2)
				Ω(err).ShouldNot(HaveOccurred())

				saramaMessage := &sarama.ConsumerMessage{
					Key:   keybin,
					Value: valbin,
				}
				messagesFromKafka = append(messagesFromKafka, saramaMessage)
				expectation = append(expectation, TestMessageType{key, value})
			}
		})

		It("should receive messages through a sarama-cluster consumer", func() {
			Ω(err).ShouldNot(HaveOccurred())

			var message TestMessageType
			done := make(chan interface{})
			go func() {
				defer GinkgoRecover()
				for i := 0; i < 5; i++ {
					err = databusConsumer.Consume(context.Background(), &message)
					Ω(err).ShouldNot(HaveOccurred())
					Ω(message.TestKey).Should(Equal(expectation[i].(TestMessageType).TestKey))
					Ω(message.TestValue).Should(Equal(expectation[i].(TestMessageType).TestValue))
				}

				close(done)
			}()
			Eventually(done).Should(BeClosed())

			// Attempting to consume again should block
			done = make(chan interface{})
			var consumeErr error
			go func() {
				consumeErr = databusConsumer.Consume(context.Background(), &message)
				close(done)
			}()

			Consistently(done).ShouldNot(BeClosed())

			// Closing the databusConsumer will cause it to return an error
			err = databusConsumer.Close()
			Ω(err).ShouldNot(HaveOccurred())
			Eventually(done).Should(BeClosed())
			Eventually(clusterConsumer.Messages()).Should(BeClosed())
			Ω(errors.Cause(consumeErr)).Should(Equal(ErrConsumerClosed))
		})

		It("should fail if we pass a value that isn't a pointer", func() {
			Ω(err).ShouldNot(HaveOccurred())
			var message int
			err = databusConsumer.Consume(context.Background(), message)
			Ω(err).Should(HaveOccurred())
		})

		It("should fail if we pass a pointer to something that isn't a struct", func() {
			Ω(err).ShouldNot(HaveOccurred())
			var message int
			err = databusConsumer.Consume(context.Background(), &message)
			Ω(err).Should(HaveOccurred())
		})

		It("should fail if our struct doesn't have a 'message-key' field", func() {
			Ω(err).ShouldNot(HaveOccurred())
			var message struct {
				TestKey   string
				TestValue int `zenkit:"message-value"`
			}
			err = databusConsumer.Consume(context.Background(), &message)
			Ω(err).Should(HaveOccurred())
		})

		It("should fail if our struct doesn't have a 'message-value' field", func() {
			Ω(err).ShouldNot(HaveOccurred())
			var message struct {
				TestKey   string `zenkit:"message-key"`
				TestValue int
			}
			err = databusConsumer.Consume(context.Background(), &message)
			Ω(err).Should(HaveOccurred())
		})

		It("should fail if our key field is not settable", func() {
			Ω(err).ShouldNot(HaveOccurred())
			var message struct {
				testKey   string `zenkit:"message-key"`
				TestValue int    `zenkit:"message-value"`
			}
			err = databusConsumer.Consume(context.Background(), &message)
			Ω(err).Should(HaveOccurred())
		})

		It("should fail if our value field is not settable", func() {
			Ω(err).ShouldNot(HaveOccurred())
			var message struct {
				TestKey   string `zenkit:"message-key"`
				testValue int    `zenkit:"message-value"`
			}
			err = databusConsumer.Consume(context.Background(), &message)
			Ω(err).Should(HaveOccurred())
		})

		It("should fail if our field types don't match the incoming message", func() {
			Ω(err).ShouldNot(HaveOccurred())
			var message struct {
				TestKey   int    `zenkit:"message-key"`
				TestValue string `zenkit:"message-value"`
			}
			err = databusConsumer.Consume(context.Background(), &message)
			Ω(err).Should(HaveOccurred())
		})

		It("should consume all errors from the error channel", func() {
			Ω(err).ShouldNot(HaveOccurred())

			clusterConsumer.errors <- errors.New("A")
			clusterConsumer.errors <- errors.New("B")

			var message TestMessageType
			done := make(chan interface{})
			go func() {
				defer close(done)
				for err == nil {
					err = databusConsumer.Consume(context.Background(), &message)
				}
			}()

			Consistently(done).ShouldNot(BeClosed())
			Consistently(clusterConsumer.errors).ShouldNot(Receive())

			databusConsumer.Close()

			Eventually(done).Should(BeClosed())
		})

		It("should consume all notifications from the notification channel", func() {
			Ω(err).ShouldNot(HaveOccurred())

			clusterConsumer.notifications <- &cluster.Notification{}
			clusterConsumer.notifications <- &cluster.Notification{}

			var message TestMessageType
			done := make(chan interface{})
			go func() {
				defer close(done)
				for err == nil {
					err = databusConsumer.Consume(context.Background(), &message)
				}
			}()

			Consistently(done).ShouldNot(BeClosed())
			Consistently(clusterConsumer.notifications).ShouldNot(Receive())

			databusConsumer.Close()

			Eventually(done).Should(BeClosed())
		})

		It("should error on close if the underlying consumer errors on close", func() {
			Ω(err).ShouldNot(HaveOccurred())

			clusterConsumer.err = errors.New("Error!")
			err = databusConsumer.Close()
			Ω(err).Should(HaveOccurred())
		})

		It("should stop consuming when the context is cancelled", func() {
			Ω(err).ShouldNot(HaveOccurred())

			var message TestMessageType
			done := make(chan interface{})
			myctx, cancel := context.WithCancel(context.Background())
			go func() {
				defer close(done)
				for err == nil {
					err = databusConsumer.Consume(myctx, &message)
				}
			}()

			Consistently(done).ShouldNot(BeClosed())

			cancel()
			Eventually(done).Should(BeClosed())
			Ω(errors.Cause(err)).Should(Equal(ErrConsumerClosed))
		})

		It("should return immediately if the context is cancelled first", func() {
			Ω(err).ShouldNot(HaveOccurred())

			var message TestMessageType
			done := make(chan interface{})
			myctx, cancel := context.WithCancel(context.Background())
			cancel()
			go func() {
				defer GinkgoRecover()
				defer close(done)
				for err == nil {
					err = databusConsumer.Consume(myctx, &message)
				}
			}()
			Eventually(done).Should(BeClosed())
			Ω(errors.Cause(err)).Should(Equal(ErrConsumerClosed))
		})
	})

	Context("with struct keys and values", func() {
		BeforeEach(func() {
			keySchema = "key-test"
			valueSchema = "val-test"
			for i := 0; i < 5; i++ {
				key := KeyTest{test.RandString(8), rand.Intn(100)}
				value := ValTest{test.RandString(8)}

				keyjson, err := json.Marshal(key)
				Ω(err).ShouldNot(HaveOccurred())
				valjson, err := json.Marshal(value)
				Ω(err).ShouldNot(HaveOccurred())

				keybin, err := AvroSerialize(keyjson, 3)
				Ω(err).ShouldNot(HaveOccurred())
				valbin, _ := AvroSerialize(valjson, 4)
				Ω(err).ShouldNot(HaveOccurred())

				saramaMessage := &sarama.ConsumerMessage{
					Key:   keybin,
					Value: valbin,
				}
				messagesFromKafka = append(messagesFromKafka, saramaMessage)
				expectation = append(expectation, TestMessageTypeStruct{key, value})
			}
		})

		It("should receive messages through a sarama-cluster consumer", func() {
			Ω(err).ShouldNot(HaveOccurred())

			message := TestMessageTypeStruct{}
			done := make(chan interface{})
			go func() {
				defer GinkgoRecover()
				for i := 0; i < 5; i++ {
					err = databusConsumer.Consume(context.Background(), &message)
					Ω(err).ShouldNot(HaveOccurred())
					Ω(message.TestKey).Should(Equal(expectation[i].(TestMessageTypeStruct).TestKey))
					Ω(message.TestValue).Should(Equal(expectation[i].(TestMessageTypeStruct).TestValue))
				}

				close(done)
			}()
			Eventually(done).Should(BeClosed())
		})
	})

	Context("with struct pointer keys and values", func() {
		BeforeEach(func() {
			keySchema = "key-test"
			valueSchema = "val-test"
			for i := 0; i < 5; i++ {
				key := &KeyTest{test.RandString(8), rand.Intn(100)}
				value := &ValTest{test.RandString(8)}

				keyjson, err := json.Marshal(key)
				Ω(err).ShouldNot(HaveOccurred())
				valjson, err := json.Marshal(value)
				Ω(err).ShouldNot(HaveOccurred())

				keybin, err := AvroSerialize(keyjson, 3)
				Ω(err).ShouldNot(HaveOccurred())
				valbin, err := AvroSerialize(valjson, 4)
				Ω(err).ShouldNot(HaveOccurred())

				saramaMessage := &sarama.ConsumerMessage{
					Key:   keybin,
					Value: valbin,
				}
				messagesFromKafka = append(messagesFromKafka, saramaMessage)
				expectation = append(expectation, TestMessageTypeStructPtr{key, value})
			}
		})

		It("should receive messages through a sarama-cluster consumer", func() {
			Ω(err).ShouldNot(HaveOccurred())

			message := TestMessageTypeStructPtr{}
			done := make(chan interface{})
			go func() {
				defer GinkgoRecover()
				for i := 0; i < 5; i++ {
					err = databusConsumer.Consume(context.Background(), &message)
					Ω(err).ShouldNot(HaveOccurred())
					Ω(message.TestKey).Should(Equal(expectation[i].(TestMessageTypeStructPtr).TestKey))
					Ω(message.TestValue).Should(Equal(expectation[i].(TestMessageTypeStructPtr).TestValue))
				}

				close(done)
			}()
			Eventually(done).Should(BeClosed())
		})
	})
})

var _ = Describe("SaramaMessage", func() {
	var (
		topic         string
		key           []byte
		value         []byte
		saramaMessage *sarama.ConsumerMessage
	)

	BeforeEach(func() {
		topic = test.RandString(8)
		key = []byte(test.RandString(8))
		value = []byte(test.RandString(8))

		saramaMessage = &sarama.ConsumerMessage{
			Topic: topic,
			Key:   key,
			Value: value,
		}
	})

	It("should have the same topic, key, and value as the sarama message", func() {
		msg := SaramaMessage{saramaMessage}
		Ω(msg.Key()).Should(Equal(key))
		Ω(msg.Value()).Should(Equal(value))
		Ω(msg.Topic()).Should(Equal(topic))
	})
})
