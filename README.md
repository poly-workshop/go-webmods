# Golang Web Modules

[![Go Reference](https://pkg.go.dev/badge/github.com/poly-workshop/go-webmods.svg)](https://pkg.go.dev/github.com/poly-workshop/go-webmods)
[![Go Report Card](https://goreportcard.com/badge/github.com/poly-workshop/go-webmods)](https://goreportcard.com/report/github.com/poly-workshop/go-webmods)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

Common solution collection for quickly setting up Golang web applications. This library provides a modular plugin architecture with factory functions and config-driven initialization for common web application components.

## 📚 Documentation

**[Full API Documentation on pkg.go.dev →](https://pkg.go.dev/github.com/poly-workshop/go-webmods)**

## 🚀 Quick Start

### Installation

```bash
go get github.com/poly-workshop/go-webmods
```

### Basic Usage

```go
package main

import (
    "log/slog"
    "github.com/poly-workshop/go-webmods/app"
    "github.com/poly-workshop/go-webmods/gorm_client"
)

func main() {
    // Initialize application with configuration
    app.SetCMDName("myapp")
    app.Init(".")  // Loads config from ./configs/

    // Initialize database
    db := gorm_client.NewDB(gorm_client.Config{
        Driver:   "postgres",
        Host:     "localhost",
        Port:     5432,
        Username: "user",
        Password: "password",
        Name:     "mydb",
        SSLMode:  "disable",
    })

    slog.Info("Application started successfully")
    _ = db
}
```

## 📦 Packages

### [app](https://pkg.go.dev/github.com/poly-workshop/go-webmods/app)
Core application utilities including:
- Layered configuration management (Viper)
- Structured logging with context propagation (slog)
- Application initialization helpers

### [gorm_client](https://pkg.go.dev/github.com/poly-workshop/go-webmods/gorm_client)
Database client factory supporting:
- PostgreSQL
- MySQL
- SQLite
- Connection pooling configuration

### [mongo_client](https://pkg.go.dev/github.com/poly-workshop/go-webmods/mongo_client)
MongoDB client factory using the v2 driver:
- MongoDB Atlas support
- Connection pooling and timeouts
- Ping verification on startup

### [redis_client](https://pkg.go.dev/github.com/poly-workshop/go-webmods/redis_client)
Redis client with:
- Single-node and cluster mode support
- Two-level caching (local + distributed)
- Automatic cache invalidation via pub/sub

### [object_storage](https://pkg.go.dev/github.com/poly-workshop/go-webmods/object_storage)
Unified object storage interface supporting:
- Local filesystem
- MinIO / S3-compatible storage
- Volcengine TOS

### [grpc_utils](https://pkg.go.dev/github.com/poly-workshop/go-webmods/grpc_utils)
gRPC server interceptors for:
- Structured logging
- Request ID generation and propagation
- Context-aware tracing

### [smtp_mailer](https://pkg.go.dev/github.com/poly-workshop/go-webmods/smtp_mailer)
SMTP email client with:
- TLS support
- HTML and plain text emails
- Multiple recipients support

## 🏗️ Architecture Patterns

### Factory Pattern
Every component uses a Config struct and factory function:
```go
component := package.NewComponent(package.Config{...})
```

### Provider Pattern
Multi-backend support with unified interfaces:
```go
storage, err := object_storage.NewObjectStorage(object_storage.Config{
    ProviderType: object_storage.ProviderLocal,
    ProviderConfig: object_storage.ProviderConfig{
        BasePath: "/data",
    },
})
```

### Context-Aware Logging
Structured logging with automatic context propagation:
```go
ctx = app.WithLogAttrs(ctx, slog.String("user_id", "123"))
slog.InfoContext(ctx, "User logged in") // Includes user_id automatically
```

## ⚙️ Configuration

Configuration uses Viper with layered loading:
1. `configs/default.yaml` - Base configuration
2. `configs/{MODE}.yaml` - Environment-specific overrides
3. Environment variables - Final overrides (using `__` separator)

Example `configs/default.yaml`:
```yaml
log:
  level: info
  format: tint

database:
  driver: postgres
  host: localhost
  port: 5432
  username: user
  password: pass
  name: mydb
  sslmode: disable
```

Set environment mode:
```bash
export MODE=production  # Loads configs/production.yaml
```

Override with environment variables:
```bash
export LOG__LEVEL=debug
export DATABASE__HOST=prod-db
```

## 📖 Examples

See the [examples in pkg.go.dev](https://pkg.go.dev/github.com/poly-workshop/go-webmods#pkg-examples) for detailed usage examples of each package.

## 🧪 Development

### Run Tests
```bash
go test ./...
```

### Format Code
```bash
make fmt
```

### Lint Code
```bash
make lint
```

## 📄 License

This repo is granted under the [MIT License](LICENSE). Feel free to use it in your projects.
