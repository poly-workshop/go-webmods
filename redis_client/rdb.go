package redis_client

import (
	"sync"

	"github.com/redis/go-redis/v9"
)

type Config struct {
	Urls     []string
	Password string
}

var (
	rdbInitOnce sync.Once
	cfg         Config
	redisClient redis.UniversalClient
)

// NewRDB creates a new Redis client with the provided configuration.
// This function supports creating multiple Redis clients with different configurations.
//
// The client automatically detects the mode based on the number of URLs:
//   - Single URL: Creates a standard Redis client
//   - Multiple URLs: Creates a Redis cluster client
//
// Example:
//
//	rdb := redis_client.NewRDB(redis_client.Config{
//	    Urls:     []string{"localhost:6379"},
//	    Password: "",
//	})
func NewRDB(cfg Config) redis.UniversalClient {
	if len(cfg.Urls) == 0 {
		panic("redis_client: no redis hosts configured")
	}
	if len(cfg.Urls) == 1 {
		return redis.NewClient(&redis.Options{
			Addr:     cfg.Urls[0],
			Password: cfg.Password,
		})
	}
	return redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:    cfg.Urls,
		Password: cfg.Password,
	})
}

// SetConfig configures the singleton Redis client.
// This must be called before GetRDB() if using the singleton pattern.
// Deprecated: Use NewRDB instead for better control and to support multiple clients.
func SetConfig(urls []string, password string) {
	cfg = Config{
		Urls:     urls,
		Password: password,
	}
}

// GetRDB returns the singleton Redis client instance.
// The client is initialized on first call using the configuration set by SetConfig.
// Subsequent calls return the same instance (thread-safe).
// Deprecated: Use NewRDB instead for better control and to support multiple clients.
func GetRDB() redis.UniversalClient {
	rdbInitOnce.Do(func() {
		redisClient = NewRDB(cfg)
	})
	return redisClient
}
