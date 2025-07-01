package interceptor

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
	app "github.com/oj-lab/go-webmods/app"
	"google.golang.org/grpc"
)

type requestIdContextKey string

const (
	RequestIdKey requestIdContextKey = "request_id"
)

// RequestIdInterceptor is a unary server interceptor that adds a request ID to the context and logs
func RequestIdInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	requestId := uuid.New().String()
	ctx = context.WithValue(ctx, RequestIdKey, requestId)
	ctx = app.WithLogAttrs(ctx, slog.String("request_id", requestId))
	return handler(ctx, req)
}
