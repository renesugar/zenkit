/*
Package databus provides a client for sending and receiving schema-encoded
messages to a databus backend. The default implementations of each
encode/decode according to Avro schemas in a schema registry, and send/receive
from a Kafka topic.

Typically, you would create a producer using `NewDatabusProducer`, passing in the addresses for Kafka and the schema registry, the Kafka topic to publish to, and the names of the key and value schemas.

	producer, _ := NewDatabusProducer([]string{"kafka-host:9200"}, "schema-registry-host:8081", "topic", "message-key-schema", "message-value-schema")
	defer producer.Close()

	key := SomeStruct{Field:"whatever"} // Should match message-key-schema
	value := MoreInterestingStruct{ID: 1234, Field:"whatever"} // Should match message-value-schema

	err := producer.Send(key, value)

Similarly, you would create a consumer using `NewDatabusConsumer`, passing in the addresses for Kafka and the schema registry, the Kafka topic to consume from, the names of the key and value schemas, and the group ID to which this consumer should belong (consumers in the same group consume as a group, rather than each consuming from the topic independently).

	type MyMessage struct {
		Key   SomeStruct            `zenkit:"message-key"`
		Value MoreInterestingStruct `zenkit:"message-value"`
	}


	consumer, _ := NewDatabusConsumer([]{"kafka-host:9200"}, "schema-registry-host:8081", "topic", "message-key-schema", "message-value-schema", "my-cool-group")
	defer consumer.Close()

	var msg MyMessage
	for consumer.Consume(&msg) == nil {
		go Process(msg) // Get a copy here, since the pointer will be reused
	}

*/
package databus
