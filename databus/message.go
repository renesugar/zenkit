package databus

import (
	"io"

	"github.com/karrick/goavro"
)

type Message interface {
	// Topic is the Kafka topic to which this message will be published
	Topic() string
	// Key is the encoded key of the message
	Key() []byte
	// Value is the encoded value of the message
	Value() []byte
}

type MessageFactory interface {
	// Topic is the Kafka topic to which these messages are destined
	Topic() string
	// KeySchema is the registry subject for the schema of the key
	KeySchema() string
	// ValueSchema is the registry subject for the schema of the value
	ValueSchema() string
	// Codec is the Avro codec that will encode the message
	Codec() *goavro.Codec
	// Message produces an encoded message, ready to publish to Kafka
	Message(key, value io.Reader) (Message, error)
}
