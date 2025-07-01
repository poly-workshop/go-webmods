package interceptor

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
	app "github.com/oj-lab/go-webmods/app"
	"google.golang.org/grpc"
)

type requestIDContextKey string

const (
	requestIDKey requestIDContextKey = "request_id"
)

// RequestIDInterceptor is a unary server interceptor that adds a request ID to the context and logs
func RequestIDInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	requestID := uuid.New().String()
	ctx = context.WithValue(ctx, requestIDKey, requestID)
	ctx = app.WithLogAttrs(ctx, slog.String("request_id", requestID))
	return handler(ctx, req)
}

// GetRequestID extracts the request ID from the context
func GetRequestID(ctx context.Context) string {
	if requestID, ok := ctx.Value(requestIDKey).(string); ok {
		return requestID
	}
	return ""
}
