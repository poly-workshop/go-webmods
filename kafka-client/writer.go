package kafka_client

import (
	"time"

	"github.com/segmentio/kafka-go"
)

type WriterConfig struct {
	Brokers          []string
	Topic            string
	Async            bool
	RequiredAcks     kafka.RequiredAcks
	Balancer         kafka.Balancer
	BatchSize        int
	BatchBytes       int
	BatchTimeout     time.Duration
	CompressionCodec kafka.CompressionCodec
	MaxAttempts      int
}

func NewWriter(cfg WriterConfig) *kafka.Writer {
	if len(cfg.Brokers) == 0 {
		panic("kafka_client: no kafka brokers configured")
	}
	if cfg.Topic == "" {
		panic("kafka_client: topic is required")
	}

	writerCfg := kafka.WriterConfig{
		Brokers:          cfg.Brokers,
		Topic:            cfg.Topic,
		Async:            cfg.Async,
		RequiredAcks:     int(cfg.RequiredAcks),
		Balancer:         cfg.Balancer,
		BatchSize:        cfg.BatchSize,
		BatchBytes:       cfg.BatchBytes,
		BatchTimeout:     cfg.BatchTimeout,
		CompressionCodec: cfg.CompressionCodec,
		MaxAttempts:      cfg.MaxAttempts,
	}

	return kafka.NewWriter(writerCfg)
}
