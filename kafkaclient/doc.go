// Package kafkaclient provides factory functions for creating Kafka readers
// and writers using github.com/segmentio/kafka-go.
//
// # Reader
//
//	reader := kafkaclient.NewReader(kafkaclient.ReaderConfig{
//	    Brokers: []string{"localhost:9092"},
//	    Topic:   "example-topic",
//	    GroupID: "example-group",
//	})
//	defer reader.Close()
//
// # Writer
//
//	writer := kafkaclient.NewWriter(kafkaclient.WriterConfig{
//	    Brokers: []string{"localhost:9092"},
//	    Topic:   "example-topic",
//	})
//	defer writer.Close()
package kafkaclient
