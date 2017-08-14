// +build integration

package databus_test

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"sync"

	schemaregistry "github.com/datamountaineer/schema-registry"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pkg/errors"
	. "github.com/zenoss/zenkit/databus"
	"github.com/zenoss/zenkit/test"
)

type messageExpectation struct {
	Key   interface{}
	Value interface{}
}

var _ = Describe("Consumer", func() {

	var (
		groupId              string
		brokers              []string
		producer             DatabusProducer
		databusConsumer      DatabusConsumer
		databusConsumer2     DatabusConsumer
		topic                string
		producerTopic        string
		schemaRegistryAddr   string
		schemaRegistryClient schemaregistry.Client
		err                  error
		expectation          []*messageExpectation
		keySchema            string
		valueSchema          string
		schemas              = map[string]string{
			"object-key":   `"string"`,
			"object-value": `"int"`,
			"key-test":     keyTestSchema,
			"val-test":     valTestSchema,
		}
	)

	BeforeEach(func() {
		topic = test.RandString(8)
		producerTopic = topic
		groupId = test.RandString(8)

		brokersstring, err := harness.Resolve("kafka", 9092)
		Ω(err).ShouldNot(HaveOccurred())
		brokers = []string{brokersstring}

		addr, err := harness.Resolve("kafka-schema-registry", 8081)
		Ω(err).ShouldNot(HaveOccurred())

		schemaRegistryAddr = fmt.Sprintf("http://%s/", addr)
		schemaRegistryClient, err = schemaregistry.NewClient(schemaRegistryAddr)
		Ω(err).ShouldNot(HaveOccurred())

		for subject, schema := range schemas {
			_, err = schemaRegistryClient.RegisterNewSchema(subject, schema)
			Ω(err).ShouldNot(HaveOccurred())
		}

		expectation = []*messageExpectation{}
		keySchema = "object-key"
		valueSchema = "object-value"
		err = nil
	})

	JustBeforeEach(func() {

		// We have to create the consumer before messages are sent by the producer
		databusConsumer, err = NewDatabusConsumer(brokers, schemaRegistryAddr, topic, keySchema, valueSchema, groupId)
		// Error checked in tests

		// Sleep for 1 second after creating the consumer, to make sure no messages are sent until after the
		//  consumer is up
		time.Sleep(time.Second)

	})

	Context("with an invalid broker", func() {
		BeforeEach(func() {
			brokers = []string{"bad"}
		})
		It("should fail to create the consumer", func() {
			Ω(err).Should(HaveOccurred())
		})
	})

	Context("with an invalid key subject", func() {
		BeforeEach(func() {
			keySchema = "bad"
		})
		It("should fail to create the consumer", func() {
			Ω(err).Should(HaveOccurred())
		})
	})

	Context("with an invalid value subject", func() {
		BeforeEach(func() {
			valueSchema = "bad"
		})
		It("should fail to create the consumer", func() {
			Ω(err).Should(HaveOccurred())
		})
	})

	Context("with an invalid schema registry", func() {
		BeforeEach(func() {
			schemaRegistryAddr = ":bad"
		})
		It("should fail to create the consumer", func() {
			Ω(err).Should(HaveOccurred())
		})
	})

	Context("with valid inputs", func() {
		JustBeforeEach(func() {
			var localErr error

			producer, localErr = NewDatabusProducer(brokers, schemaRegistryAddr, producerTopic, keySchema, valueSchema)
			Ω(localErr).ShouldNot(HaveOccurred())

			// Send messages via a producer
			for _, exp := range expectation {
				localErr = producer.Send(exp.Key, exp.Value)
				Ω(localErr).ShouldNot(HaveOccurred())
			}
		})

		Context("with primitive keys and values", func() {
			BeforeEach(func() {
				for i := 0; i < 5; i++ {
					key := test.RandString(8)
					value := rand.Intn(100)
					expectation = append(expectation, &messageExpectation{key, value})
				}
			})

			It("should receive messages through a sarama-cluster consumer", func() {
				Ω(err).ShouldNot(HaveOccurred())

				var message TestMessageType
				done := make(chan interface{})
				go func() {
					for i := 0; i < 5; i++ {
						err = databusConsumer.Consume(context.Background(), &message)
						Ω(err).ShouldNot(HaveOccurred())
						Ω(message.TestKey).Should(Equal(expectation[i].Key.(string)))
						Ω(message.TestValue).Should(Equal(expectation[i].Value.(int)))
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
				Ω(errors.Cause(consumeErr)).Should(Equal(ErrConsumerClosed))
			})

			Context("with a different topic", func() {
				BeforeEach(func() {
					topic = "bad-topic"
				})

				It("should not receive any messages", func() {
					Ω(err).ShouldNot(HaveOccurred())

					var message TestMessageType
					done := make(chan interface{})
					var consumeErr error
					go func() {
						consumeErr = databusConsumer.Consume(context.Background(), &message)
						close(done)
					}()

					Consistently(done).ShouldNot(BeClosed())
					err = databusConsumer.Close()
					Ω(err).ShouldNot(HaveOccurred())
					Eventually(done).Should(BeClosed())
				})
			})

			Context("with a second consumer using a different group ID", func() {
				BeforeEach(func() {
					// We have to create the consumer before messages are sent by the producer
					var err2 error
					databusConsumer2, err2 = NewDatabusConsumer(brokers, schemaRegistryAddr, topic, keySchema, valueSchema, "different-group")
					Ω(err2).ShouldNot(HaveOccurred())

					// Sleep for 1 second after creating the consumer, to make sure no messages are sent until after the
					//  consumer is up
					time.Sleep(time.Second)
				})

				It("should still receive all messages", func() {
					Ω(err).ShouldNot(HaveOccurred())

					var message TestMessageType
					done := make(chan interface{})
					go func() {
						for i := 0; i < 5; i++ {
							err = databusConsumer2.Consume(context.Background(), &message)
							Ω(err).ShouldNot(HaveOccurred())
							Ω(message.TestKey).Should(Equal(expectation[i].Key.(string)))
							Ω(message.TestValue).Should(Equal(expectation[i].Value.(int)))
						}
						databusConsumer2.Close()
						close(done)
					}()

					Eventually(done).Should(BeClosed())

					done = make(chan interface{})
					go func() {
						for i := 0; i < 5; i++ {
							err = databusConsumer.Consume(context.Background(), &message)
							Ω(err).ShouldNot(HaveOccurred())
							Ω(message.TestKey).Should(Equal(expectation[i].Key.(string)))
							Ω(message.TestValue).Should(Equal(expectation[i].Value.(int)))
						}
						databusConsumer.Close()
						close(done)
					}()

					Eventually(done).Should(BeClosed())
				})
			})

			Context("with a second consumer using the same group ID", func() {
				BeforeEach(func() {
					// We have to create the consumer before messages are sent by the producer
					var err2 error
					databusConsumer2, err2 = NewDatabusConsumer(brokers, schemaRegistryAddr, topic, keySchema, valueSchema, groupId)
					Ω(err2).ShouldNot(HaveOccurred())

					// Sleep for 1 second after creating the consumer, to make sure no messages are sent until after the
					//  consumer is up
					time.Sleep(time.Second)
				})

				It("should still receive all messages between the two", func() {
					Ω(err).ShouldNot(HaveOccurred())

					var resultLock sync.Mutex
					resultCount := 0
					countResult := func() {
						resultLock.Lock()
						defer resultLock.Unlock()
						resultCount++
					}

					// Start both consumers consuming, count the number of messages
					done := make(chan interface{})
					go func() {
						var message TestMessageType
						var localErr error
						for localErr == nil {
							localErr = databusConsumer.Consume(context.Background(), &message)
							if localErr == nil {
								countResult()
							}
						}
						close(done)
					}()

					done2 := make(chan interface{})
					go func() {
						var message TestMessageType
						var localErr error
						for localErr == nil {
							localErr = databusConsumer2.Consume(context.Background(), &message)
							if localErr == nil {
								countResult()
							}
						}
						close(done2)
					}()

					// Make sure we never read more than 5 results total
					Consistently(func() bool {
						resultLock.Lock()
						defer resultLock.Unlock()
						return resultCount <= len(expectation)
					}()).Should(BeTrue())

					// Stop the consumers
					databusConsumer.Close()
					databusConsumer2.Close()
					Eventually(done).Should(BeClosed(), "consumer 1 did not exit")
					Eventually(done2).Should(BeClosed(), "consumer 2 did not exit")

					// Make sure we read exactly 5 results
					Ω(resultCount).Should(Equal(len(expectation)))
				})
			})

		})

		Context("with struct keys and values", func() {
			BeforeEach(func() {
				keySchema = "key-test"
				valueSchema = "val-test"
				for i := 0; i < 5; i++ {
					key := KeyTest{test.RandString(8), rand.Intn(100)}
					value := ValTest{test.RandString(8)}

					expectation = append(expectation, &messageExpectation{key, value})
				}
			})

			It("should receive messages through a sarama-cluster consumer", func() {
				Ω(err).ShouldNot(HaveOccurred())

				message := TestMessageTypeStruct{}
				done := make(chan interface{})
				go func() {
					for i := 0; i < 5; i++ {
						err = databusConsumer.Consume(context.Background(), &message)
						Ω(err).ShouldNot(HaveOccurred())
						Ω(message.TestKey).Should(Equal(expectation[i].Key.(KeyTest)))
						Ω(message.TestValue).Should(Equal(expectation[i].Value.(ValTest)))
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

					expectation = append(expectation, &messageExpectation{key, value})
				}
			})

			It("should receive messages through a sarama-cluster consumer", func() {
				Ω(err).ShouldNot(HaveOccurred())

				message := TestMessageTypeStructPtr{}
				done := make(chan interface{})
				go func() {
					for i := 0; i < 5; i++ {
						err = databusConsumer.Consume(context.Background(), &message)
						Ω(err).ShouldNot(HaveOccurred())
						Ω(message.TestKey).Should(Equal(expectation[i].Key.(*KeyTest)))
						Ω(message.TestValue).Should(Equal(expectation[i].Value.(*ValTest)))
					}

					close(done)
				}()
				Eventually(done, 5*time.Second).Should(BeClosed())
			})
		})
	})

})
