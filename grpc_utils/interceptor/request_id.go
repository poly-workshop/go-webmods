package interceptor

import (
	"context"

	"github.com/google/uuid"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"google.golang.org/grpc"
)

// RequestIDInterceptor is a unary server interceptor that adds a request ID to the context and logs
func RequestIDInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	requestID := uuid.New().String()
	// Use go-grpc-middleware's InjectFields for standardized field injection
	ctx = logging.InjectFields(ctx, logging.Fields{"request_id", requestID})
	return handler(ctx, req)
}
