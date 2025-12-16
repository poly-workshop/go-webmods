package kafka_client

import (
	"testing"

	"github.com/segmentio/kafka-go"
)

func TestNewReaderRequiresBrokers(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic when brokers are missing")
		}
	}()

	NewReader(ReaderConfig{
		Topic: "test-topic",
	})
}

func TestNewReaderRequiresTopic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic when topic is missing")
		}
	}()

	NewReader(ReaderConfig{
		Brokers: []string{"localhost:9092"},
	})
}

func TestNewReaderAppliesConfig(t *testing.T) {
	reader := NewReader(ReaderConfig{
		Brokers:     []string{"localhost:9092"},
		Topic:       "test-topic",
		GroupID:     "group-1",
		MinBytes:    1,
		MaxBytes:    10,
		StartOffset: kafka.LastOffset,
	})
	t.Cleanup(func() {
		if err := reader.Close(); err != nil {
			t.Logf("failed to close reader: %v", err)
		}
	})

	cfg := reader.Config()
	if cfg.Topic != "test-topic" {
		t.Fatalf("expected topic %q, got %q", "test-topic", cfg.Topic)
	}
	if cfg.GroupID != "group-1" {
		t.Fatalf("expected group %q, got %q", "group-1", cfg.GroupID)
	}
	if cfg.MinBytes != 1 {
		t.Fatalf("expected min bytes %d, got %d", 1, cfg.MinBytes)
	}
	if cfg.MaxBytes != 10 {
		t.Fatalf("expected max bytes %d, got %d", 10, cfg.MaxBytes)
	}
	if cfg.StartOffset != kafka.LastOffset {
		t.Fatalf("expected start offset %d, got %d", kafka.LastOffset, cfg.StartOffset)
	}
}
