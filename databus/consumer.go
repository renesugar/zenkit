package databus

import (
	"context"
	"reflect"

	"github.com/Shopify/sarama"
	cluster "github.com/bsm/sarama-cluster"
	"github.com/datamountaineer/schema-registry"
	"github.com/pkg/errors"
)

const (
	StructTagKey        = "message-key"
	StructTagValue      = "message-value"
	StructTagIdentifier = "zenkit"
)

var (
	// ErrInvalidMessageType is thrown when a key or value that can't be
	// decoded according to the schema is provided
	ErrInvalidMessageType = errors.New("invalid message type")
	// ErrConsumerClosed is thrown when the databus consumer is no longer open
	ErrConsumerClosed = errors.New("consumer is closed")
)

// DatabusConsumer is capable of receiving schema-encoded messages from a databus.
type DatabusConsumer interface {
	// Consume reads a message from the databus and applies it to the object at
	// the pointer provided.
	Consume(context.Context, interface{}) error
	// Close closes the consumer.
	Close() error
}

// NewDatabusConsumer returns the default implementation of a DatabusConsumer,
// which reads Avro-encoded messages from a Kafka consumer.
func NewDatabusConsumer(brokers []string, schemaRegistry, topic, keySubject, valueSubject, groupId string) (DatabusConsumer, error) {
	// Get our schema registry
	schemaRegistryClient, err := schemaregistry.NewClient(schemaRegistry)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create schema registry client")
	}

	// Get our message factory
	messageFactory, err := NewMessageFactory(topic, keySubject, valueSubject, schemaRegistryClient)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create message factory")
	}

	// Get our sarama cluster consumer
	// init (custom) config, disable errors and notifications
	config := cluster.NewConfig()
	config.Consumer.Return.Errors = false
	config.Group.Return.Notifications = false

	consumer, err := cluster.NewConsumer(brokers, groupId, []string{topic}, config)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create cluster consumer")
	}
	return NewSaramaClusterDatabusConsumer(consumer, messageFactory)
}

// NewSaramaClusterDatabusConsumer returns a SaramaDatabusConsumer created from
// an existing cluster.Consumer directly, instead of creating a new one from
// broker addresses.
func NewSaramaClusterDatabusConsumer(consumer SaramaClusterConsumer, messageFactory MessageFactory) (DatabusConsumer, error) {

	c := &saramaClusterDatabusConsumer{
		con:            consumer,
		messageFactory: messageFactory,
	}
	return c, nil
}

// SaramaClusterConsumer is implemented by sarama's cluster.Consumer, and is
// here for the sake of tests, since sarama doesn't provide an interface.
type SaramaClusterConsumer interface {
	Errors() <-chan error
	Notifications() <-chan *cluster.Notification
	Messages() <-chan *sarama.ConsumerMessage
	MarkOffset(*sarama.ConsumerMessage, string)
	Close() error
}

// saramaClusterDatabusConsumer is the default implementation of
// DatabusConsumer, which receives Avro-encoded messages from a Kafka databus.
type saramaClusterDatabusConsumer struct {
	con            SaramaClusterConsumer
	messageFactory MessageFactory
}

func (c *saramaClusterDatabusConsumer) Consume(ctx context.Context, v interface{}) error {
	// Make sure what was passed in is a pointer to messageType
	keyField, valueField, err := c.validateType(v)
	if err != nil {
		return errors.WithStack(err)
	}

	// Get errors and notifications first
	stop := false
	for !stop {
		select {
		case <-c.con.Errors():
		// TODO
		case <-c.con.Notifications():
		// TODO
		case <-ctx.Done():
			stop = true
		default:
			stop = true
		}
	}

	select {
	case msg, more := <-c.con.Messages():
		if more {
			err := c.decodeMessage(msg, v, keyField, valueField)
			if err != nil {
				return errors.Wrap(err, "failed to decode message")
			}

			c.con.MarkOffset(msg, "") // mark message as processed
			return nil
		} else {
			return errors.Wrap(ErrConsumerClosed, "messages channel closed")
		}
	case <-ctx.Done():
		return errors.Wrap(ErrConsumerClosed, "context is cancelled")
	}
}

func (c *saramaClusterDatabusConsumer) Close() error {
	return c.con.Close()
}

// validateType ensures that the message type is valid. It must be a pointer to
// a struct with fields tagged as `zenkit:"message-key"` and
// `zenkit:"message-value"`.
func (c *saramaClusterDatabusConsumer) validateType(message interface{}) (int, int, error) {
	messageType := reflect.TypeOf(message)
	if messageType.Kind() != reflect.Ptr {
		return 0, 0, errors.Wrap(ErrInvalidMessageType, "type is not a pointer")
	}
	mType := messageType.Elem()
	if mType.Kind() != reflect.Struct {
		return 0, 0, errors.Wrap(ErrInvalidMessageType, "type does not point to a struct")
	}
	keyField := -1
	valueField := -1
	for i := 0; i < mType.NumField(); i++ {
		tag := mType.Field(i).Tag.Get(StructTagIdentifier)
		if tag == StructTagKey {
			keyField = i
		} else if tag == StructTagValue {
			valueField = i
		}
	}

	if keyField < 0 || valueField < 0 {
		return 0, 0, errors.Wrap(ErrInvalidMessageType, "missing key or value field")
	}

	if !reflect.ValueOf(message).Elem().Field(keyField).CanSet() {
		return 0, 0, errors.Wrap(ErrInvalidMessageType, "key field is not settable")
	}

	if !reflect.ValueOf(message).Elem().Field(valueField).CanSet() {
		return 0, 0, errors.Wrap(ErrInvalidMessageType, "value field is not settable")
	}

	return keyField, valueField, nil
}

func (c *saramaClusterDatabusConsumer) decodeMessage(rawMsg *sarama.ConsumerMessage, v interface{}, keyField, valueField int) error {
	messageType := reflect.TypeOf(v).Elem()

	// Get a Message object we can pass to MessageFactory.Decode
	wrappedMessage := &SaramaMessage{rawMsg}

	// Decode the Key and Value
	key := reflect.New(messageType.Field(keyField).Type).Interface()
	value := reflect.New(messageType.Field(valueField).Type).Interface()
	err := c.messageFactory.Decode(wrappedMessage, key, value)
	if err != nil {
		return errors.Wrap(err, "failed to decode key or value")
	}

	// Populate v with the key and value
	reflect.ValueOf(v).Elem().Field(keyField).Set(reflect.ValueOf(key).Elem())
	reflect.ValueOf(v).Elem().Field(valueField).Set(reflect.ValueOf(value).Elem())

	return nil
}

// SaramaMessage is a simple wrapper to make a sarama ConsumerMessage implement
// Message
type SaramaMessage struct {
	Message *sarama.ConsumerMessage
}

func (m *SaramaMessage) Topic() string {
	return m.Message.Topic
}

func (m *SaramaMessage) Key() []byte {
	return m.Message.Key
}
func (m *SaramaMessage) Value() []byte {
	return m.Message.Value
}
