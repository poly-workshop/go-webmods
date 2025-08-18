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

func SetConfig(urls []string, password string) {
	cfg = Config{
		Urls:     urls,
		Password: password,
	}
}

func GetRDB() redis.UniversalClient {
	rdbInitOnce.Do(func() {
		if len(cfg.Urls) == 0 {
			panic("No redis hosts configured")
		}
		if len(cfg.Urls) == 1 {
			redisClient = redis.NewClient(&redis.Options{
				Addr:     cfg.Urls[0],
				Password: cfg.Password,
			})
		}
		if len(cfg.Urls) > 1 {
			redisClient = redis.NewClusterClient(&redis.ClusterOptions{
				Addrs:    cfg.Urls,
				Password: cfg.Password,
			})
		}
	})
	return redisClient
}
