package databus

import (
	"github.com/Shopify/sarama"
)

type MessageProducer interface {
	Send(Message) error
	Close() error
}

func NewMessageProducer(producer sarama.SyncProducer) MessageProducer {
	return &kafkaProducer{
		prod: producer,
	}
}

type kafkaProducer struct {
	prod sarama.SyncProducer
}

func (p *kafkaProducer) Send(msg Message) error {
	_, _, err := p.prod.SendMessage(&sarama.ProducerMessage{
		Topic: msg.Topic(),
		Key:   sarama.ByteEncoder(msg.Key()),
		Value: sarama.ByteEncoder(msg.Value()),
	})

	return err
}

func (p *kafkaProducer) Close() error {
	return p.prod.Close()
}
