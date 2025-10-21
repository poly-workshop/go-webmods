// Package app provides core application utilities including configuration management,
// structured logging with context propagation, and initialization helpers.
//
// # Configuration Management
//
// The package uses Viper to provide layered configuration from multiple sources:
//   - default.yaml: Base configuration
//   - {MODE}.yaml: Environment-specific overrides (development.yaml, production.yaml, etc.)
//   - Environment variables: Final overrides using __ as separator (e.g., LOG__LEVEL)
//
// Configuration files should be placed in a configs/ directory relative to the
// working directory passed to Init().
//
// Example config structure:
//
//	workdir/
//	├── configs/
//	│   ├── default.yaml
//	│   ├── development.yaml
//	│   └── production.yaml
//	└── main.go
//
// Example default.yaml:
//
//	log:
//	  level: info
//	  format: tint
//	database:
//	  driver: postgres
//	  host: localhost
//	  port: 5432
//
// # Initialization
//
// Initialize the application at startup:
//
//	func main() {
//	    app.SetCMDName("myapp")  // Optional: sets command name in logs
//	    app.Init(".")            // Loads config from ./configs/
//	    // ... rest of application
//	}
//
// The MODE environment variable controls which config file is loaded:
//
//	export MODE=production  # Loads production.yaml
//	export MODE=development # Loads development.yaml (default)
//
// # Logging
//
// The package provides structured logging using Go's standard log/slog with
// automatic context propagation and enhanced features:
//
// Supported log formats (configured via log.format in config):
//   - tint: Pretty colored console output (default)
//   - json: JSON structured logs
//   - plain-text: Plain text logs
//
// Log levels: debug, info, warn, error
//
// Basic logging:
//
//	import "log/slog"
//
//	slog.Info("Server started", "port", 8080)
//	slog.Error("Database connection failed", "error", err)
//
// Context-aware logging:
//
//	ctx = app.WithLogAttrs(ctx, slog.String("user_id", "123"))
//	ctx = app.WithLogAttrs(ctx, slog.String("request_id", "abc"))
//	slog.InfoContext(ctx, "Processing request") // Includes user_id and request_id
//
// All log messages automatically include:
//   - cmd: Command name (if set via SetCMDName)
//   - hostname: Current hostname
//   - Any attributes added to the context via WithLogAttrs
//
// # Configuration Access
//
// Access configuration values using the Config() function:
//
//	cfg := app.Config()
//	dbHost := cfg.GetString("database.host")
//	dbPort := cfg.GetInt("database.port")
//	enableFeature := cfg.GetBool("features.new_feature")
//
// Environment variables override config file values:
//
//	export LOG__LEVEL=debug          # Overrides log.level
//	export DATABASE__HOST=prod-db    # Overrides database.host
//
// # Best Practices
//
//   - Call Init() once at application startup
//   - Use SetCMDName() to identify different services/commands in logs
//   - Add request/user context via WithLogAttrs for better traceability
//   - Use structured logging (key-value pairs) instead of string formatting
//   - Configure log format and level via config files, not hardcoded
package app
