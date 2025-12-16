package kafka_client

import (
	"github.com/segmentio/kafka-go"
)

type ReaderConfig struct {
	Brokers     []string
	Topic       string
	GroupID     string
	Partition   int
	MinBytes    int
	MaxBytes    int
	StartOffset int64
}

func NewReader(cfg ReaderConfig) *kafka.Reader {
	if len(cfg.Brokers) == 0 {
		panic("kafka_client: no Kafka brokers configured")
	}
	if cfg.Topic == "" {
		panic("kafka_client: topic is required")
	}

	readerCfg := kafka.ReaderConfig{
		Brokers:     cfg.Brokers,
		Topic:       cfg.Topic,
		GroupID:     cfg.GroupID,
		Partition:   cfg.Partition,
		MinBytes:    cfg.MinBytes,
		MaxBytes:    cfg.MaxBytes,
		StartOffset: cfg.StartOffset,
	}

	return kafka.NewReader(readerCfg)
}
