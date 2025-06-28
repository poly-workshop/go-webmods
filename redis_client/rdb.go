package redis_client

import (
	"github.com/redis/go-redis/v9"
)

type Config struct {
	Urls     []string
	Password string
}

func NewRDB(cfg Config) redis.UniversalClient {
	if len(cfg.Urls) == 0 {
		panic("No redis hosts configured")
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
