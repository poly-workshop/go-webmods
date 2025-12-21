// Package go-webmods provides a collection of common web application components
// for Go services. It follows a modular plugin architecture where each package
// provides factory functions and config-driven initialization.
//
// # Overview
//
// This library includes the following packages:
//
//   - app: Core application utilities for configuration, logging, and context management
//   - gormclient: Database client factory supporting PostgreSQL and SQLite
//   - redisclient: Redis client with caching support and cluster mode
//   - kafkaclient: Kafka client factories for readers and writers
//   - objectstorage: Multi-provider object storage interface (local, MinIO, Volcengine TOS)
//   - grpcutils: gRPC middleware and interceptors for logging and request ID tracking
//   - smtpmailer: SMTP email sender with TLS support
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
//	import gormclient "github.com/poly-workshop/go-webmods/gormclient"
//
//	db := gormclient.NewDB(gormclient.Config{
//	    Driver:   "postgres",
//	    Host:     "localhost",
//	    Port:     5432,
//	    Username: "user",
//	    Password: "pass",
//	    DbName:   "dbname",
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
//	storage, err := objectstorage.NewObjectStorage(objectstorage.Config{
//	    ProviderType: objectstorage.ProviderLocal,
//	    ProviderConfig: objectstorage.ProviderConfig{
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
//	    gormclient "github.com/poly-workshop/go-webmods/gormclient"
//	    redisclient "github.com/poly-workshop/go-webmods/redisclient"
//	)
//
//	func main() {
//	    // Initialize application
//	    app.SetCMDName("myapp")
//	    app.Init(".")
//
//	    // Initialize database
//	    db := gormclient.NewDB(gormclient.Config{
//	        Driver: "sqlite",
//	        DbName: "data/app.db",
//	    })
//
//	    // Initialize Redis
//	    redisclient.SetConfig([]string{"localhost:6379"}, "")
//	    rdb := redisclient.GetRDB()
//
//	    slog.Info("Application started successfully")
//	    // ... rest of application logic
//	}
//
// # License
//
// This project is licensed under the MIT License.
package go_webmods
