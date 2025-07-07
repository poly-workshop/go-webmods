package app

import (
	"context"
	"log/slog"
	"os"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/lmittmann/tint"
)

type contextKey string

const (
	configKeyLogLevel  = "log.level"
	configKeyLogFormat = "log.format"

	logFormatJSON      = "json"
	logFormatPlainText = "plain-text"
	logFormatTint      = "tint"

	logAttrsKey contextKey = "log_attrs"
)

func stringToSlogLevel(
	level string,
) slog.Level {
	switch level {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

type logHandler struct {
	slog.Handler
}

func WithLogAttrs(ctx context.Context, attrs ...slog.Attr) context.Context {
	existingAttrs, ok := ctx.Value(logAttrsKey).([]slog.Attr)
	if !ok {
		existingAttrs = []slog.Attr{}
	}
	return context.WithValue(ctx, logAttrsKey, append(existingAttrs, attrs...))
}

func (h *logHandler) Handle(ctx context.Context, r slog.Record) error {
	if cmdName != "" {
		r.AddAttrs(slog.String("cmd", cmdName))
	}
	if hostname != "" {
		r.AddAttrs(slog.String("hostname", hostname))
	}

	// Add log fields from context (app-specific)
	if attrs, ok := ctx.Value(logAttrsKey).([]slog.Attr); ok {
		for _, v := range attrs {
			r.AddAttrs(v)
		}
	}

	// Add log fields from go-grpc-middleware context
	if grpcFields := logging.ExtractFields(ctx); len(grpcFields) > 0 {
		for i := 0; i < len(grpcFields); i += 2 {
			if i+1 < len(grpcFields) {
				key, ok := grpcFields[i].(string)
				if ok {
					value := grpcFields[i+1]
					r.AddAttrs(slog.Any(key, value))
				}
			}
		}
	}

	return h.Handler.Handle(ctx, r)
}

func initLog() {
	logLevel := stringToSlogLevel(config.GetString(configKeyLogLevel))

	var handler slog.Handler
	switch config.GetString(configKeyLogFormat) {
	case logFormatJSON:
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: logLevel})
	case logFormatPlainText:
		handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: logLevel})
	default:
		handler = tint.NewHandler(os.Stdout, &tint.Options{Level: logLevel})
	}
	handler = &logHandler{handler}
	slog.SetDefault(slog.New(handler))
}
