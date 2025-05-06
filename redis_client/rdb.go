package redis_client

import (
	"sync"

	"github.com/oj-lab/go-webmods/app"
	"github.com/redis/go-redis/v9"
)

const (
	configKeyRedisUrls     = "redis.urls"
	configKeyRedisPassword = "redis.password"
)

var (
	rdbInitMutex sync.Mutex
)

var redisClient redis.UniversalClient

func GetRDB() redis.UniversalClient {
	if redisClient == nil {
		rdbInitMutex.Lock()
		defer rdbInitMutex.Unlock()
		if redisClient != nil {
			return redisClient
		}

		redisUrls := app.Config().GetStringSlice(configKeyRedisUrls)
		password := app.Config().GetString(configKeyRedisPassword)

		if len(redisUrls) == 0 {
			panic("No redis hosts configured")
		}
		if len(redisUrls) == 1 {
			redisClient = redis.NewClient(&redis.Options{
				Addr:     redisUrls[0],
				Password: password,
			})
		}
		if len(redisUrls) > 1 {
			redisClient = redis.NewClusterClient(&redis.ClusterOptions{
				Addrs:    redisUrls,
				Password: password,
			})
		}
	}
	return redisClient
}
