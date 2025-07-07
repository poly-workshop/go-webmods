package interceptor

import (
	"context"
	"testing"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"google.golang.org/grpc"
)

func TestRequestIDInterceptor(t *testing.T) {
	// Create a mock handler
	mockHandler := func(ctx context.Context, req any) (any, error) {
		// Check if the request ID is in the context using go-grpc-middleware's ExtractFields
		fields := logging.ExtractFields(ctx)
		if len(fields) == 0 {
			t.Error("Expected fields to be injected into context, but none found")
			return nil, nil
		}
		
		// Check if request_id is in the fields
		found := false
		for i := 0; i < len(fields); i += 2 {
			if i+1 < len(fields) && fields[i] == "request_id" {
				found = true
				if requestID, ok := fields[i+1].(string); ok {
					if requestID == "" {
						t.Error("Expected request_id to be non-empty")
					}
				} else {
					t.Error("Expected request_id to be a string")
				}
				break
			}
		}
		
		if !found {
			t.Error("Expected request_id field to be present in context")
		}
		
		return "test response", nil
	}

	// Create a mock grpc.UnaryServerInfo
	info := &grpc.UnaryServerInfo{
		Server:     nil,
		FullMethod: "/test.Service/TestMethod",
	}

	// Call the interceptor
	ctx := context.Background()
	resp, err := RequestIDInterceptor(ctx, "test request", info, mockHandler)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	
	if resp != "test response" {
		t.Errorf("Expected response 'test response', got %v", resp)
	}
}

func TestRequestIDInterceptorIntegration(t *testing.T) {
	// This test verifies that the RequestIDInterceptor works with the app logging system
	// by checking that the fields are properly injected into the context
	var capturedRequestID string
	
	// Create a handler that captures the context and extracts the request ID
	mockHandler := func(ctx context.Context, req any) (any, error) {
		// Extract fields from go-grpc-middleware context
		fields := logging.ExtractFields(ctx)
		for i := 0; i < len(fields); i += 2 {
			if i+1 < len(fields) && fields[i] == "request_id" {
				if requestID, ok := fields[i+1].(string); ok {
					capturedRequestID = requestID
				}
				break
			}
		}
		
		return "test response", nil
	}

	// Create a mock grpc.UnaryServerInfo
	info := &grpc.UnaryServerInfo{
		Server:     nil,
		FullMethod: "/test.Service/TestMethod",
	}

	// Call the interceptor
	ctx := context.Background()
	_, err := RequestIDInterceptor(ctx, "test request", info, mockHandler)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	
	// Verify that the request_id was captured and is a valid UUID
	if capturedRequestID == "" {
		t.Error("Expected request_id to be captured, but it was empty")
	}
	
	// Simple validation that it looks like a UUID (has dashes and is the right length)
	if len(capturedRequestID) != 36 {
		t.Errorf("Expected request_id to be 36 characters (UUID format), got %d", len(capturedRequestID))
	}
}