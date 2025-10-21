package app_test

import (
	"context"
	"log/slog"

	"github.com/poly-workshop/go-webmods/app"
)

// Example demonstrates basic application initialization with configuration
// and logging setup.
func Example() {
	// Set command name (appears in all logs)
	app.SetCMDName("myapp")

	// Initialize application (loads config from ./configs/)
	// Note: This would normally be called with a valid config directory
	// app.Init(".")

	// Access configuration
	// cfg := app.Config()
	// dbHost := cfg.GetString("database.host")
	// logLevel := cfg.GetString("log.level")

	// Use structured logging
	slog.Info("Application started", "version", "1.0.0")
}

// Example_withLogAttrs demonstrates adding contextual information to logs
// that automatically propagates through the request lifecycle.
func Example_withLogAttrs() {
	ctx := context.Background()

	// Add user context
	ctx = app.WithLogAttrs(ctx, slog.String("user_id", "123"))
	ctx = app.WithLogAttrs(ctx, slog.String("action", "login"))

	// All logs with this context will include user_id and action
	slog.InfoContext(ctx, "User logged in")

	// Add more context as needed
	ctx = app.WithLogAttrs(ctx, slog.String("session_id", "abc-def"))
	slog.InfoContext(ctx, "Session created")
}

// Example_configuration demonstrates accessing configuration values from
// different sources (config files and environment variables).
func Example_configuration() {
	// Initialize application
	// app.Init(".")

	// Get configuration
	// cfg := app.Config()

	// Access string values
	// dbHost := cfg.GetString("database.host")
	// fmt.Println("Database host:", dbHost)

	// Access integer values
	// dbPort := cfg.GetInt("database.port")
	// fmt.Println("Database port:", dbPort)

	// Access boolean values
	// enableDebug := cfg.GetBool("debug")
	// fmt.Println("Debug enabled:", enableDebug)

	// Access with defaults
	// timeout := cfg.GetInt("timeout")
	// if timeout == 0 {
	//     timeout = 30
	// }

	// Environment variables override config files
	// export DATABASE__HOST=prod-db  # Overrides database.host
	// export LOG__LEVEL=debug        # Overrides log.level
}
