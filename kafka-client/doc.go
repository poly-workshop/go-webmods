// Package kafka_client provides factory functions for creating Kafka readers
// and writers using github.com/segmentio/kafka-go.
//
// # Reader
//
//	reader := kafka_client.NewReader(kafka_client.ReaderConfig{
//	    Brokers: []string{"localhost:9092"},
//	    Topic:   "example-topic",
//	    GroupID: "example-group",
//	})
//	defer reader.Close()
//
// # Writer
//
//	writer := kafka_client.NewWriter(kafka_client.WriterConfig{
//	    Brokers: []string{"localhost:9092"},
//	    Topic:   "example-topic",
//	})
//	defer writer.Close()
package kafka_client
