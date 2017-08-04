package databus

import (
	"github.com/Shopify/sarama"
	"github.com/datamountaineer/schema-registry"
	"github.com/pkg/errors"
)

// DatabusProducer is capable of sending messages with a given key and value to
// a databus.
type DatabusProducer interface {
	Send(key, value interface{}) error
	Close() error
}

// NewDatabusProducer returns the default implementation of DatabusProducer,
// which sends Avro-encoded messages to a Kafka topic.
func NewDatabusProducer(brokers []string, schemaRegistry, topic, keySubject, valueSubject string) (DatabusProducer, error) {
	schemaRegistryClient, err := schemaregistry.NewClient(schemaRegistry)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create schema registry client")
	}

	messageFactory, err := NewMessageFactory(topic, keySubject, valueSubject, schemaRegistryClient)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create message factory")
	}

	producer, err := sarama.NewSyncProducer(brokers, nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create sarama producer")
	}
	return NewSaramaDatabusProducer(producer, messageFactory), nil
}

// NewSaramaDatabusProducer is a way to create a sarama-based DatabusProducer
// using an existing SyncProducer, in distinction to NewDatabusProducer, which
// creates a new SyncProducer from broker addresses.
func NewSaramaDatabusProducer(producer sarama.SyncProducer, factory MessageFactory) DatabusProducer {
	return &saramaDatabusProducer{producer, factory}
}

// saramaDatabusProducer is the default implementation of DatabusProducer. It
// sends Avro-encoded messages to a Kafka-based databus.
type saramaDatabusProducer struct {
	producer sarama.SyncProducer
	factory  MessageFactory
}

func (s *saramaDatabusProducer) Send(key, value interface{}) error {
	message, err := s.factory.Message(key, value)
	if err != nil {
		return errors.Wrap(err, "failed to get message from factory")
	}

	_, _, err = s.producer.SendMessage(&sarama.ProducerMessage{
		Topic: message.Topic(),
		Key:   sarama.ByteEncoder(message.Key()),
		Value: sarama.ByteEncoder(message.Value()),
	})

	return errors.Wrap(err, "failed to send message via sarama producer")
}

func (s *saramaDatabusProducer) Close() error {
	return s.producer.Close()
}
