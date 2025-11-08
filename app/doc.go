// Package app provides core application utilities including configuration management,
// structured logging with context propagation, and initialization helpers.
//
// # Configuration Management
//
// The package uses Viper to provide layered configuration from multiple sources.
// Configuration is loaded in the following priority order (highest priority last):
//
// 1. Global base: default.yaml
// 2. Global environment: {MODE}.yaml (development.yaml, production.yaml, etc.)
// 3. Command base: {cmdName}/default.yaml (replaces all previous config if exists)
// 4. Command environment: {cmdName}/{MODE}.yaml (merges with command base)
// 5. Environment variables: Final overrides using __ as separator (e.g., LOG__LEVEL)
//
// IMPORTANT: Command-specific default.yaml completely replaces global configuration
// rather than merging with it. Use command-specific configs only when you need
// completely different settings for specific commands.
//
// Configuration files should be placed in a configs/ directory. The Init() function
// automatically looks for configs in the working directory, or you can specify a
// custom path with InitWithConfigPath().
//
// Example config structure:
//
//	configs/
//	├── default.yaml              # Global defaults
//	├── development.yaml          # Global development overrides
//	├── production.yaml           # Global production overrides
//	├── myapp/
//	│   ├── default.yaml          # App-specific config (replaces global!)
//	│   ├── development.yaml      # App-specific development overrides
//	│   └── production.yaml       # App-specific production overrides
//	└── worker/
//	    ├── default.yaml          # Worker-specific config (replaces global!)
//	    └── production.yaml       # Worker-specific production overrides
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
// Initialize the application at startup with one of two methods:
//
//	func main() {
//	    // Method 1: Auto-detect config path from working directory
//	    app.Init("myapp")                    // Loads config from ./configs/
//
//	    // Method 2: Specify custom config path
//	    app.InitWithConfigPath("myapp", "/etc/myapp/configs")
//	    // ... rest of application
//	}
//
// The command name passed to Init() functions serves multiple purposes:
//   - Appears in all log messages as the "cmd" field
//   - Used to find command-specific configuration files
//   - Helps identify different services/commands in logs
//
// The MODE environment variable controls which config file is loaded:
//
//	export MODE=production  # Loads production.yaml + myapp/production.yaml
//	export MODE=development # Loads development.yaml + myapp/development.yaml (default)
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
//   - cmd: Command name (passed to Init functions)
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
//   - Call Init() or InitWithConfigPath() once at application startup
//   - Use descriptive command names to identify different services/commands in logs
//   - Use global configs for shared settings, command-specific configs sparingly
//   - Remember that command-specific default.yaml replaces (not merges) global config
//   - Add request/user context via WithLogAttrs for better traceability
//   - Use structured logging (key-value pairs) instead of string formatting
//   - Configure log format and level via config files, not hardcoded
package app
