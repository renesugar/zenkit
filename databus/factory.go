package databus

import (
	"encoding/json"

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
	// Decode decodes a message received from Kafka into the types provided
	Decode(msg Message, key, value interface{}) error
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
	encodedKey, err := encode(f.keyCodec, key)
	if err != nil {
		return nil, err
	}
	encodedValue, err := encode(f.valCodec, value)
	if err != nil {
		return nil, err
	}
	return NewMessage(f.topic, encodedKey, encodedValue), nil
}

func (f *avroMessageFactory) Decode(msg Message, key, value interface{}) error {
	if err := decode(f.keyCodec, msg.Key(), key); err != nil {
		return err
	}
	return decode(f.valCodec, msg.Value(), value)
}

func encode(codec *goavro.Codec, data interface{}) ([]byte, error) {
	marshalled, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	var v interface{}
	json.Unmarshal(marshalled, &v)
	return codec.TextualFromNative([]byte{}, v)
}

func decode(codec *goavro.Codec, data []byte, ptr interface{}) error {
	native, _, err := codec.NativeFromTextual(data)
	if err != nil {
		return err
	}
	marshalled, _ := json.Marshal(native)
	return json.Unmarshal(marshalled, ptr)
}
