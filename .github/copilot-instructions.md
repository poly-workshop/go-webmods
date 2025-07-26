# Copilot Instructions for go-webmods

## Overview
This is a shared Go module library providing common web application components for OJ Lab services. It follows a **modular plugin architecture** where each package provides factory functions and config-driven initialization.

## Architecture Principles

### 1. Factory Pattern with Config Structs
Every component uses a `Config` struct and factory function:
- Database: `gorm_client.NewDB(gorm_client.Config{})`
- Redis: `redis_client.NewRDB(redis_client.Config{})`  
- Object Storage: `object_storage.NewObjectStorage(object_storage.Config{})`

### 2. Provider Pattern for Multi-Backend Support
Object storage demonstrates the provider pattern:
```go
// Config specifies both provider type and provider-specific config
type Config struct {
    ProviderType   // "local" or "volcengine" 
    ProviderConfig // unified config struct for all providers
}
```

### 3. Context-Aware Logging
The `app` package provides centralized logging with context propagation:
- Use `app.WithLogAttrs(ctx, slog.String("key", "value"))` to add fields
- Log handler automatically adds `cmd` and `hostname` to all logs
- Supports JSON, plain-text, and tint formats via config

## Key Development Patterns

### Configuration Management
- Uses Viper with layered config: `default.yaml` → `{MODE}.yaml` → environment variables
- Environment variables use `__` separator: `log.level` → `LOG__LEVEL`
- Initialize with `app.Init(workdir)` where `workdir/configs/` contains config files

### Database Pattern
```go
// Always check driver type in switch statement
switch cfg.Driver {
case "postgres":
    return openPostgres(cfg)
case "sqlite": 
    return openSqlite(cfg)
default:
    panic(fmt.Sprintf("unsupported database driver: %s", cfg.Driver))
}
```

### Interface Design
Follow the object storage interface pattern:
- Define interface with core operations (`Save`, `List`, `Open`, `Stat`, `Delete`)
- Create wrapper types that implement `io.ReadSeekCloser` 
- Use `os.FileInfo` for metadata consistency

### gRPC Integration
Use the provided interceptors for consistent request handling:
```go
grpcutils.BuildLogInterceptor(logger)      // Structured logging
grpcutils.BuildRequestIDInterceptor()      // Request ID propagation
```

## Development Workflow

### Testing
- Run tests: `go test ./...`
- Each package has comprehensive tests demonstrating usage patterns
- Use temporary directories for file-based tests: `os.MkdirTemp("", "test_prefix")`

### Code Quality
- Format: `make fmt` (uses golines + gofumpt)  
- Lint: `make lint` (golangci-lint with --fix)
- CI runs `golangci-lint` via GitHub Actions

### Adding New Providers
1. Add provider type constant to main config
2. Implement the interface with provider-specific struct
3. Add case to factory function switch statement
4. Create comprehensive tests covering all interface methods
5. Update main `NewXXX` factory function

## File Organization
- `app/`: Core application utilities (config, logging, context)
- `*_client/`: Database and cache clients with simple factory patterns
- `object_storage/`: Multi-provider storage with interface-based design
- `grpc_utils/`: gRPC middleware and interceptors

## Critical Dependencies
- **Logging**: `github.com/lmittmann/tint` for pretty console output
- **Config**: `github.com/spf13/viper` for layered configuration
- **Database**: GORM with postgres/sqlite drivers  
- **Redis**: `github.com/redis/go-redis/v9` with cluster support
- **gRPC**: Standard gRPC middleware for logging and request ID tracking
