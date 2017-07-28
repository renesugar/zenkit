package databus

type Message interface {
	// Topic is the Kafka topic to which this message will be published
	Topic() string
	// Key is the encoded key of the message
	Key() []byte
	// Value is the encoded value of the message
	Value() []byte
}

func NewMessage(topic string, key, value []byte) Message {
	return &defaultMessage{topic, key, value}
}

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
