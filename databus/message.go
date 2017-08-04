package databus

// Message is a low-level message that is sent to or received from the databus.
// If a schema is being used, the key and value are already encoded.
type Message interface {
	// Topic is the databus topic to which this message will be published
	Topic() string
	// Key is the encoded key of the message
	Key() []byte
	// Value is the encoded value of the message
	Value() []byte
}

// NewMessage creates a new Message with the values provided.
func NewMessage(topic string, key, value []byte) Message {
	return &defaultMessage{topic, key, value}
}

// defaultMessage is the default implementation of Message, which simply passes
// through values without any further modification.
type defaultMessage struct {
	topic string
	key   []byte
	value []byte
}

func (m *defaultMessage) Topic() string {
	return m.topic
}

func (m *defaultMessage) Key() []byte {
	return m.key
}

func (m *defaultMessage) Value() []byte {
	return m.value
}
