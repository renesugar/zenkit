package databus

import (
	"encoding/json"

	schemaregistry "github.com/datamountaineer/schema-registry"
	"github.com/pkg/errors"
)

// MessageFactory wraps a topic, key subject and value subject. It handles
// creating messages from Go structures according to the schemas provided, and
// decoding messages according to the schema.
type MessageFactory interface {
	// Topic is the Kafka topic to which these messages are destined
	Topic() string
	// KeySubject is the registry subject for the schema of the key
	KeySubject() string
	// ValueSubject is the registry subject for the schema of the value
	ValueSubject() string
	// KeyCodec is the Avro codec that will encode the message key
	KeyCodec() SchemaCodec
	// ValueCodec is the Avro codec that will encode the message value
	ValueCodec() SchemaCodec
	// Message produces an encoded message, ready to publish to Kafka
	Message(key, value interface{}) (Message, error)
	// Decode decodes a message received from Kafka into the types provided
	Decode(msg Message, key, value interface{}) error
}

// NewMessageFactory creates a MessageFactory using the topic, key and value
// schemas provided.
func NewMessageFactory(topic, keySubject, valueSubject string, client schemaregistry.Client) (MessageFactory, error) {
	_, keyCodec, err := GetCodec(client, keySubject)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get codec for key subject: %s", keySubject)
	}
	_, valCodec, err := GetCodec(client, valueSubject)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get codec for value subject: %s", valueSubject)
	}
	return &avroMessageFactory{
		topic:      topic,
		keySubject: keySubject,
		valSubject: valueSubject,
		keyCodec:   keyCodec,
		valCodec:   valCodec,
	}, nil
}

// avroMessageFactory is the default implementation of MessageFactory. It uses
// goavro.Codecs to do the encoding/decoding.
type avroMessageFactory struct {
	topic      string
	keySubject string
	valSubject string
	keyCodec   SchemaCodec
	valCodec   SchemaCodec
}

func (f *avroMessageFactory) Topic() string {
	return f.topic
}

func (f *avroMessageFactory) KeyCodec() SchemaCodec {
	return f.keyCodec

}
func (f *avroMessageFactory) ValueCodec() SchemaCodec {
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
		return nil, errors.Wrap(err, "failed to encode key")
	}
	encodedValue, err := encode(f.valCodec, value)
	if err != nil {
		return nil, errors.Wrap(err, "failed to encode value")
	}
	return NewMessage(f.topic, encodedKey, encodedValue), nil
}

func (f *avroMessageFactory) Decode(msg Message, key, value interface{}) error {
	if err := decode(f.keyCodec, msg.Key(), key); err != nil {
		return errors.Wrap(err, "failed to decode key")
	}
	if err := decode(f.valCodec, msg.Value(), value); err != nil {
		return errors.Wrap(err, "failed to decode value")
	}
	return nil
}

// encode massages the data specified into Go native types via JSON
// marshal/unmarshal, then encodes it using the SchemaEncoder provided.
func encode(codec SchemaEncoder, data interface{}) ([]byte, error) {
	marshalled, err := json.Marshal(data)
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal data to json")
	}
	var v interface{}
	json.Unmarshal(marshalled, &v)
	return codec.TextualFromNative([]byte{}, v)
}

// decode decodes data using the SchemaDecoder provided, then applies it to the
// pointer provided via JSON marshal/unmarshal.
func decode(codec SchemaDecoder, data []byte, ptr interface{}) error {
	native, _, err := codec.NativeFromTextual(data)
	if err != nil {
		return errors.Wrap(err, "failed to get native value from text")
	}
	marshalled, _ := json.Marshal(native)
	return json.Unmarshal(marshalled, ptr)
}
