package workers

import (
	"log"
	"strings"
	"time"

	"gopkg.in/Shopify/sarama.v1"
)

// Producer ...
type Producer struct {
	kafkaProducer sarama.AsyncProducer
}

//NewProducer ...
func NewProducer() Producer {
	brokerList := strings.Split("localhost:9092", ",")
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForLocal       // Only wait for the leader to ack
	config.Producer.Compression = sarama.CompressionSnappy   // Compress messages
	config.Producer.Flush.Frequency = 500 * time.Millisecond // Flush batches every 500ms

	producer, err := sarama.NewAsyncProducer(brokerList, config)
	if err != nil {
		log.Fatalln("Failed to start Sarama producer:", err)
	}
	return Producer{kafkaProducer: producer}
}

//Produce produces a test message
func (p Producer) Produce(value string) {
	p.kafkaProducer.Input() <- &sarama.ProducerMessage{
		Topic: "access_log",
		Key:   sarama.StringEncoder("my_key"),
		Value: sarama.StringEncoder(value),
	}
}
