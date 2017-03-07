package workers

import (
	"log"
	"strings"
	"time"

	"gopkg.in/Shopify/sarama.v1"
)

//Produce produces a test message
func Produce(value string) {
	brokerList := strings.Split("localhost:9092", ",")
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForLocal       // Only wait for the leader to ack
	config.Producer.Compression = sarama.CompressionSnappy   // Compress messages
	config.Producer.Flush.Frequency = 500 * time.Millisecond // Flush batches every 500ms

	producer, err := sarama.NewAsyncProducer(brokerList, config)
	if err != nil {
		log.Fatalln("Failed to start Sarama producer:", err)
	}

	producer.Input() <- &sarama.ProducerMessage{
		Topic: "access_log",
		Key:   sarama.StringEncoder("CNA00001"),
		Value: sarama.StringEncoder(value),
	}
}
