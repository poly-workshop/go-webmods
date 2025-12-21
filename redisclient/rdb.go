package redisclient

import (
	"github.com/redis/go-redis/v9"
)

type Config struct {
	Urls     []string
	Password string
}

// NewRDB creates a new Redis client with the provided configuration.
// This function supports creating multiple Redis clients with different configurations.
//
// The client automatically detects the mode based on the number of URLs:
//   - Single URL: Creates a standard Redis client
//   - Multiple URLs: Creates a Redis cluster client
//
// Example:
//
//	rdb := redisclient.NewRDB(redisclient.Config{
//	    Urls:     []string{"localhost:6379"},
//	    Password: "",
//	})
func NewRDB(cfg Config) redis.UniversalClient {
	if len(cfg.Urls) == 0 {
		panic("redisclient: no redis hosts configured")
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


