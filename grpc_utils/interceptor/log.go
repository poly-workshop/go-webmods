package interceptor

import (
	"context"
	"log/slog"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"google.golang.org/grpc"
)

// Creates a gRPC interceptor that logs messages using the provided slog.Logger.
func BuildLogInterceptor(l *slog.Logger) grpc.UnaryServerInterceptor {
	return logging.UnaryServerInterceptor(logging.LoggerFunc(
		func(ctx context.Context, lvl logging.Level, msg string, fields ...any) {
			l.Log(ctx, slog.Level(lvl), msg, fields...)
		},
	))
}
