package app

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"testing"
	"time"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
)

func TestLogHandlerWithGrpcMiddlewareFields(t *testing.T) {
	var buf bytes.Buffer
	
	// Create a JSON handler to capture structured output
	jsonHandler := slog.NewJSONHandler(&buf, &slog.HandlerOptions{Level: slog.LevelDebug})
	handler := &logHandler{jsonHandler}
	
	// Create a context with go-grpc-middleware fields
	ctx := logging.InjectFields(context.Background(), logging.Fields{"request_id", "test-123", "user_id", "user-456"})
	
	// Create a log record
	record := slog.NewRecord(time.Now(), slog.LevelInfo, "test message", 0)
	
	// Handle the log record
	err := handler.Handle(ctx, record)
	if err != nil {
		t.Fatalf("Failed to handle log record: %v", err)
	}
	
	// Parse the JSON output
	var logEntry map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &logEntry); err != nil {
		t.Fatalf("Failed to parse log JSON: %v", err)
	}
	
	// Verify that the grpc-middleware fields are present
	if requestID, ok := logEntry["request_id"]; !ok || requestID != "test-123" {
		t.Errorf("Expected request_id to be 'test-123', got %v", requestID)
	}
	
	if userID, ok := logEntry["user_id"]; !ok || userID != "user-456" {
		t.Errorf("Expected user_id to be 'user-456', got %v", userID)
	}
	
	// Verify the message is still there
	if msg, ok := logEntry["msg"]; !ok || msg != "test message" {
		t.Errorf("Expected msg to be 'test message', got %v", msg)
	}
}