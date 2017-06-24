package workers

import (
	"log"
	"time"

	"gopkg.in/Shopify/sarama.v1"
)

// Producer ...
type Producer struct {
	kafkaProducer sarama.AsyncProducer
}

// NewProducer creates a new Kafka producer
// brokers is a list of Kakfa brokers for eg. "localhost:9092,localhost:9093"
func NewProducer(brokers []string) Producer {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForLocal       // Only wait for the leader to ack
	config.Producer.Compression = sarama.CompressionSnappy   // Compress messages
	config.Producer.Flush.Frequency = 500 * time.Millisecond // Flush batches every 500ms

	producer, err := sarama.NewAsyncProducer(brokers, config)
	if err != nil {
		log.Fatalln("Failed to start Sarama producer:", err)
	}
	return Producer{kafkaProducer: producer}
}

//Produce produces a test message
func (p Producer) Produce(value []byte, topic string) {
	p.kafkaProducer.Input() <- &sarama.ProducerMessage{
		Topic: topic,
		Key:   sarama.StringEncoder("my_key"),
		Value: sarama.ByteEncoder(value),
	}
}
