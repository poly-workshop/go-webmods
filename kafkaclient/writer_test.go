package kafkaclient

import (
	"testing"
	"time"

	"github.com/segmentio/kafka-go"
)

func TestNewWriterRequiresBrokers(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic when brokers are missing")
		}
	}()

	NewWriter(WriterConfig{
		Topic: "test-topic",
	})
}

func TestNewWriterRequiresTopic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic when topic is missing")
		}
	}()

	NewWriter(WriterConfig{
		Brokers: []string{"localhost:9092"},
	})
}

func TestNewWriterAppliesConfig(t *testing.T) {
	writer := NewWriter(WriterConfig{
		Brokers:      []string{"localhost:9092"},
		Topic:        "test-topic",
		Async:        true,
		RequiredAcks: kafka.RequireAll,
		BatchSize:    2,
		BatchTimeout: 2 * time.Second,
	})
	t.Cleanup(func() {
		if err := writer.Close(); err != nil {
			t.Logf("failed to close writer: %v", err)
		}
	})

	stats := writer.Stats()
	if stats.Topic != "test-topic" {
		t.Fatalf("expected topic %q, got %q", "test-topic", stats.Topic)
	}
	if !stats.Async {
		t.Fatal("expected async writer to be true")
	}
	if stats.RequiredAcks != int64(kafka.RequireAll) {
		t.Fatalf("expected required acks %d, got %d", kafka.RequireAll, stats.RequiredAcks)
	}
	if stats.MaxBatchSize != 2 {
		t.Fatalf("expected max batch size %d, got %d", 2, stats.MaxBatchSize)
	}
	if stats.BatchTimeout != 2*time.Second {
		t.Fatalf("expected batch timeout %v, got %v", 2*time.Second, stats.BatchTimeout)
	}
}
