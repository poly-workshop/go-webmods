package app_test

import (
	"context"
	"log/slog"

	"github.com/poly-workshop/go-webmods/app"
)

// Example demonstrates basic application initialization with configuration
// and logging setup.
func Example() {
	// Initialize application with command name (loads config from ./configs/)
	// Note: This would normally be called with a valid config directory
	// app.Init("myapp")

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
	// Initialize application with command name
	// app.Init("myapp")

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

// Example_customConfigPath demonstrates using a custom configuration directory
// path instead of the default ./configs/ location.
func Example_customConfigPath() {
	// Initialize with custom config path
	// app.InitWithConfigPath("worker", "/etc/myapp/configs")

	// This will load configuration in the following order:
	// 1. /etc/myapp/configs/default.yaml
	// 2. /etc/myapp/configs/{MODE}.yaml (e.g., production.yaml)
	// 3. /etc/myapp/configs/worker/default.yaml
	// 4. /etc/myapp/configs/worker/{MODE}.yaml
	// 5. Environment variables (highest priority)

	// Access the merged configuration
	// cfg := app.Config()
	// workerThreads := cfg.GetInt("worker.threads")
	// slog.Info("Worker initialized", "threads", workerThreads)
}
