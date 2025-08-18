package redis_client

import (
	"context"
	"log/slog"
	"sync"
	"time"

	"github.com/go-redis/cache/v9"
)

var (
	cacheRefreshEventChannel = "cacheRefreshEventChannel"
	cacheClient              *cache.Cache
	cacheInstance            *Cache
	cacheInstanceOnce        sync.Once
)

// Customize name of the cache refresh event channel
func SetCacheRefreshEventChannel(channel string) {
	cacheRefreshEventChannel = channel
}

type Cache struct {
	*cache.Cache
}

func GetCache() *Cache {
	cacheInstanceOnce.Do(func() {
		cacheClient = cache.New(&cache.Options{
			Redis:      GetRDB(),
			LocalCache: cache.NewTinyLFU(1000, time.Minute),
		})
		cacheInstance = &Cache{Cache: cacheClient}

		// Subscribe cache refresh event
		ctx := context.Background()
		_, err := GetRDB().Set(ctx, cacheRefreshEventChannel, cacheRefreshEventChannel, 0).Result()
		if err != nil {
			panic(err)
		}
		go func() {
			pubsub := GetRDB().Subscribe(ctx, cacheRefreshEventChannel)
			defer func() {
				err := pubsub.Close()
				if err != nil {
					slog.Error("Error closing pubsub", "error", err)
				}
			}()
			slog.Info(
				"Subscribed to cache refresh event channel", "channel", cacheRefreshEventChannel)
			ch := pubsub.Channel()
			for {
				select {
				case msg := <-ch:
					if msg.Payload == cacheRefreshEventChannel {
						slog.Info("Cache refresh event received", "key", msg.Payload)
						cacheInstance.DeleteFromLocalCache(msg.Payload)
					}
				case <-ctx.Done():
					return
				}
			}
		}()
	})
	return cacheInstance
}

func publishCacheRefreshEvent(ctx context.Context, key string) error {
	return GetRDB().Publish(ctx, cacheRefreshEventChannel, key).Err()
}

func (c *Cache) Get(ctx context.Context, key string, value any) error {
	return c.Cache.Get(ctx, key, value)
}

func (c *Cache) Set(ctx context.Context, key string, value any, expiration time.Duration) error {
	if err := c.Cache.Set(&cache.Item{
		Ctx:   ctx,
		Key:   key,
		Value: value,
		TTL:   expiration,
	}); err != nil {
		return err
	}
	return publishCacheRefreshEvent(ctx, key)
}

func (c *Cache) Delete(ctx context.Context, key string) error {
	if err := c.Cache.Delete(ctx, key); err != nil {
		return err
	}
	return publishCacheRefreshEvent(ctx, key)
}
