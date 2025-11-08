package app

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"os"
	"testing"
	"time"
)

func TestWithLogFields(t *testing.T) {
	// Test cases
	tests := []struct {
		name     string
		initial  []slog.Attr
		add      []slog.Attr
		expected []slog.Attr
	}{
		{
			name:    "add_first_fields",
			initial: nil,
			add: []slog.Attr{
				slog.String("key1", "value1"),
			},
			expected: []slog.Attr{
				slog.String("key1", "value1"),
			},
		},
		{
			name: "add_multiple_fields",
			initial: []slog.Attr{
				slog.String("key1", "value1"),
			},
			add: []slog.Attr{
				slog.String("key2", "value2"),
				slog.Int("key3", 3),
			},
			expected: []slog.Attr{
				slog.String("key1", "value1"),
				slog.String("key2", "value2"),
				slog.Int("key3", 3),
			},
		},
		{
			name: "override_existing_field",
			initial: []slog.Attr{
				slog.String("key1", "old_value"),
			},
			add: []slog.Attr{
				slog.String("key1", "new_value"),
			},
			expected: []slog.Attr{
				slog.String("key1", "old_value"),
				slog.String("key1", "new_value"),
			},
		},
		{
			name: "mixed_type_fields",
			initial: []slog.Attr{
				slog.Int("count", 1),
				slog.String("status", "pending"),
			},
			add: []slog.Attr{
				slog.Float64("price", 99.99),
				slog.Bool("active", true),
			},
			expected: []slog.Attr{
				slog.Int("count", 1),
				slog.String("status", "pending"),
				slog.Float64("price", 99.99),
				slog.Bool("active", true),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			// Set initial fields if any
			if tt.initial != nil {
				ctx = WithLogAttrs(ctx, tt.initial...)
			}

			// Add new fields
			ctx = WithLogAttrs(ctx, tt.add...)

			// Get fields from context
			attrs, ok := ctx.Value(logAttrsKey).([]slog.Attr)
			if !ok {
				t.Fatal("Failed to get log fields from context")
			}

			// Check number of fields
			if len(attrs) != len(tt.expected) {
				t.Errorf("Expected %d fields, got %d", len(tt.expected), len(attrs))
			}

			// Verify each field's key and value
			for i, want := range tt.expected {
				if i >= len(attrs) {
					t.Errorf("Missing expected field at index %d: %v", i, want)
					continue
				}
				got := attrs[i]
				if got.Key != want.Key {
					t.Errorf("Field %d: expected key %q, got %q", i, want.Key, got.Key)
				}
				if !compareAttrValues(got.Value, want.Value) {
					t.Errorf(
						"Field %d (%s): expected value %v, got %v",
						i,
						got.Key,
						want.Value,
						got.Value,
					)
				}
			}
		})
	}
}

// Helper function to compare slog.Value
func compareAttrValues(got, want slog.Value) bool {
	if got.Kind() != want.Kind() {
		return false
	}
	switch got.Kind() {
	case slog.KindString:
		return got.String() == want.String()
	case slog.KindInt64:
		return got.Int64() == want.Int64()
	case slog.KindFloat64:
		return got.Float64() == want.Float64()
	case slog.KindBool:
		return got.Bool() == want.Bool()
	default:
		return got.String() == want.String()
	}
}

func TestLogHandler(t *testing.T) {
	// Setup test environment
	cmdName = "test_cmd"
	hostname, _ = os.Hostname()

	// Create a buffer to capture output
	var buf bytes.Buffer

	// Create a test handler
	handler := &logHandler{
		Handler: slog.NewJSONHandler(&buf, &slog.HandlerOptions{Level: slog.LevelDebug}),
	}

	// Test context with fields
	ctx := WithLogAttrs(context.Background(),
		slog.String("test_key", "test_value"),
		slog.Int("test_number", 42),
	)

	// Create a test record
	record := slog.Record{
		Time:    time.Now(),
		Message: "test message",
		Level:   slog.LevelInfo,
	}

	// Test handling
	err := handler.Handle(ctx, record)
	if err != nil {
		t.Errorf("Handler.Handle() error = %v", err)
	}

	// Parse and verify the JSON output
	var logEntry map[string]any
	if err := json.Unmarshal(buf.Bytes(), &logEntry); err != nil {
		t.Fatalf("Failed to parse JSON output: %v", err)
	}

	// Verify required fields
	tests := []struct {
		field    string
		expected any
	}{
		{"msg", "test message"},
		{"level", "INFO"},
		{"test_key", "test_value"},
		{"test_number", float64(42)}, // JSON numbers are float64
		{"cmd", "test_cmd"},          // Verify command name
		{"hostname", hostname},       // Verify hostname
	}

	for _, tt := range tests {
		got, exists := logEntry[tt.field]
		if !exists {
			t.Errorf("Field %q not found in log output", tt.field)
			continue
		}
		if got != tt.expected {
			t.Errorf("Field %q = %v, want %v", tt.field, got, tt.expected)
		}
	}

	// Verify hostname and cmd if they are set
	if hostname != "" {
		if got, ok := logEntry["hostname"]; !ok || got == "" {
			t.Error("Expected hostname in log output")
		} else if got != hostname {
			t.Errorf("Expected hostname to be %s, got %s", hostname, got)
		}
	}
	if cmdName != "" {
		if got, ok := logEntry["cmd"]; !ok || got == "" {
			t.Error("Expected cmd in log output")
		} else if got != "test_cmd" {
			t.Errorf("Expected cmd to be test_cmd, got %s", got)
		}
	}
}
