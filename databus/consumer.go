package databus

//import (
//	cluster "github.com/bsm/sarama-cluster"
//	"context"
//	"github.com/Shopify/sarama"
//)
//
//type MessageConsumer interface {
//	Receive(Message) Message
//	Close() error
//}
//
//func NewMessageConsumer(ctx context.Context, consumer cluster.Consumer) MessageConsumer {
//	return &kafkaConsumer{
//		con: consumer,
//		messages: make(chan Message),
//	}
//}
//
//type SaramaMessage struct {
//	Message *sarama.ConsumerMessage
//}
//
//func (m *SaramaMessage) Topic() string {
//	return m.Message.Topic
//}
//
//func (m *SaramaMessage) Key() []byte {
//	return m.Message.Key
//}
//
//func (m *SaramaMessage) Value() []byte {
//	return m.Message.Value
//}
//
//type kafkaConsumer struct {
//	con      cluster.Consumer
//	messages chan Message
//}
//
//func (p *kafkaConsumer) start(ctx context.Context) {
//	for {
//		select {
//		case msg, more := <-p.con.Messages():
//			if more {
//				wrappedMessage := &SaramaMessage{msg}
//				p.messages <- wrappedMessage
//				p.con.MarkOffset(msg, "")        // mark message as processed
//			}
//		case <-p.con.Errors():
//		// TODO
//		case <-p.con.Notifications():
//		// TODO
//		case <-ctx.Done():
//			return
//		}
//	}
//}
//
//func (p *kafkaConsumer) Messages() <-chan Message {
//	return p.messages
//}
//
//func (p *kafkaConsumer) Close() error {
//	return p.con.Close()
//}
