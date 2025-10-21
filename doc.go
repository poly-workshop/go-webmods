// Package go-webmods provides a collection of common web application components
// for Go services. It follows a modular plugin architecture where each package
// provides factory functions and config-driven initialization.
//
// # Overview
//
// This library includes the following packages:
//
//   - app: Core application utilities for configuration, logging, and context management
//   - gorm_client: Database client factory supporting PostgreSQL and SQLite
//   - redis_client: Redis client with caching support and cluster mode
//   - object_storage: Multi-provider object storage interface (local, MinIO, Volcengine TOS)
//   - grpc_utils: gRPC middleware and interceptors for logging and request ID tracking
//   - smtp_mailer: SMTP email sender with TLS support
//
// # Installation
//
// To install the package, use:
//
//	go get github.com/poly-workshop/go-webmods
//
// # Quick Start
//
// The typical usage pattern follows these steps:
//
// 1. Initialize the application with configuration:
//
//	import "github.com/poly-workshop/go-webmods/app"
//
//	func main() {
//	    app.Init("/path/to/workdir") // Loads configs from workdir/configs/
//	    // ... rest of application
//	}
//
// 2. Initialize components using their factory functions:
//
//	import "github.com/poly-workshop/go-webmods/gorm_client"
//
//	db := gorm_client.NewDB(gorm_client.Config{
//	    Driver:   "postgres",
//	    Host:     "localhost",
//	    Port:     5432,
//	    Username: "user",
//	    Password: "pass",
//	    Name:     "dbname",
//	    SSLMode:  "disable",
//	})
//
// # Configuration
//
// The app package uses Viper for layered configuration:
//   - default.yaml: Default configuration
//   - {MODE}.yaml: Environment-specific configuration (e.g., development.yaml, production.yaml)
//   - Environment variables: Override with LOG__LEVEL, DB__HOST, etc. (using __ separator)
//
// Set the MODE environment variable to control which config file is loaded:
//
//	export MODE=production
//
// # Architecture Patterns
//
// Factory Pattern: Every component uses a Config struct and factory function:
//
//	component := package.NewComponent(package.Config{...})
//
// Provider Pattern: Multi-backend support (e.g., object storage):
//
//	storage, err := object_storage.NewObjectStorage(object_storage.Config{
//	    ProviderType: object_storage.ProviderLocal,
//	    ProviderConfig: object_storage.ProviderConfig{
//	        BasePath: "/data",
//	    },
//	})
//
// Context-Aware Logging: Structured logging with context propagation:
//
//	ctx = app.WithLogAttrs(ctx, slog.String("user_id", "123"))
//	slog.InfoContext(ctx, "User logged in") // Automatically includes user_id
//
// # Example Application Structure
//
// A typical application using go-webmods might look like:
//
//	package main
//
//	import (
//	    "log/slog"
//	    "github.com/poly-workshop/go-webmods/app"
//	    "github.com/poly-workshop/go-webmods/gorm_client"
//	    "github.com/poly-workshop/go-webmods/redis_client"
//	)
//
//	func main() {
//	    // Initialize application
//	    app.SetCMDName("myapp")
//	    app.Init(".")
//
//	    // Initialize database
//	    db := gorm_client.NewDB(gorm_client.Config{
//	        Driver: "sqlite",
//	        Name:   "data/app.db",
//	    })
//
//	    // Initialize Redis
//	    redis_client.SetConfig([]string{"localhost:6379"}, "")
//	    rdb := redis_client.GetRDB()
//
//	    slog.Info("Application started successfully")
//	    // ... rest of application logic
//	}
//
// # License
//
// This project is licensed under the MIT License.
package go_webmods
