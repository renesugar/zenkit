package databus

import (
	schemaregistry "github.com/datamountaineer/schema-registry"
	"github.com/karrick/goavro"
)

type MessageFactory interface {
	// Topic is the Kafka topic to which these messages are destined
	Topic() string
	// KeySubject is the registry subject for the schema of the key
	KeySubject() string
	// ValueSubject is the registry subject for the schema of the value
	ValueSubject() string
	// KeyCodec is the Avro codec that will encode the message key
	KeyCodec() *goavro.Codec
	// ValueCodec is the Avro codec that will encode the message value
	ValueCodec() *goavro.Codec
	// Message produces an encoded message, ready to publish to Kafka
	Message(key, value interface{}) (Message, error)
}

func NewMessageFactory(topic, keySubject, valueSubject string, client schemaregistry.Client) (MessageFactory, error) {
	_, keyCodec, err := GetCodec(client, keySubject)
	if err != nil {
		return nil, err
	}
	_, valCodec, err := GetCodec(client, valueSubject)
	if err != nil {
		return nil, err
	}
	return &avroMessageFactory{
		topic:      topic,
		keySubject: keySubject,
		valSubject: valueSubject,
		keyCodec:   keyCodec,
		valCodec:   valCodec,
	}, nil
}

type avroMessageFactory struct {
	topic      string
	keySubject string
	valSubject string
	keyCodec   *goavro.Codec
	valCodec   *goavro.Codec
}

func (f *avroMessageFactory) Topic() string {
	return f.topic
}

func (f *avroMessageFactory) KeyCodec() *goavro.Codec {
	return f.keyCodec

}
func (f *avroMessageFactory) ValueCodec() *goavro.Codec {
	return f.valCodec
}

func (f *avroMessageFactory) KeySubject() string {
	return f.keySubject
}

func (f *avroMessageFactory) ValueSubject() string {
	return f.valSubject
}

func (f *avroMessageFactory) Message(key, value interface{}) (Message, error) {
	encodedKey, err := f.keyCodec.TextualFromNative([]byte{}, key)
	if err != nil {
		return nil, err
	}
	encodedValue, err := f.valCodec.TextualFromNative([]byte{}, value)
	if err != nil {
		return nil, err
	}
	return NewMessage(f.topic, encodedKey, encodedValue), nil
}
