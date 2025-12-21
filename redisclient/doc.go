// Package redisclient provides Redis client factories with support for both
// single-node and cluster modes, plus an integrated caching layer with local
// memory cache and distributed cache invalidation.
//
// # Factory Pattern (Recommended)
//
// Create Redis clients using the factory pattern to support multiple independent clients:
//
//	import redisclient "github.com/poly-workshop/go-webmods/redisclient"
//
//	func main() {
//	    // Create a Redis client
//	    rdb := redisclient.NewRDB(redisclient.Config{
//	        Urls:     []string{"localhost:6379"},
//	        Password: "",
//	    })
//
//	    // Use Redis client
//	    ctx := context.Background()
//	    rdb.Set(ctx, "key", "value", 0)
//	    val, err := rdb.Get(ctx, "key").Result()
//	}
//
// The factory pattern allows creating multiple independent Redis clients:
//
//	primaryRDB := redisclient.NewRDB(redisclient.Config{
//	    Urls:     []string{"primary:6379"},
//	    Password: "secret1",
//	})
//
//	cacheRDB := redisclient.NewRDB(redisclient.Config{
//	    Urls:     []string{"cache:6379"},
//	    Password: "secret2",
//	})
//
// The client automatically detects cluster mode:
//   - Single URL: Creates a standard Redis client
//   - Multiple URLs: Creates a Redis cluster client
//
// Redis Cluster example:
//
//	rdb := redisclient.NewRDB(redisclient.Config{
//	    Urls:     []string{"node1:6379", "node2:6379", "node3:6379"},
//	    Password: "password",
//	})
//
// # Caching Layer
//
// The cache provides a two-level caching system:
//   - Local in-memory cache (TinyLFU, configurable size and TTL)
//   - Distributed Redis cache
//   - Automatic cache invalidation via pub/sub
//
// Factory pattern cache creation (recommended):
//
//	import redisclient "github.com/poly-workshop/go-webmods/redisclient"
//
//	func main() {
//	    rdb := redisclient.NewRDB(redisclient.Config{
//	        Urls:     []string{"localhost:6379"},
//	        Password: "",
//	    })
//
//	    cache := redisclient.NewCache(redisclient.CacheConfig{
//	        Redis:               rdb,
//	        RefreshEventChannel: "myapp:cache:refresh",
//	        LocalCacheSize:      2000,
//	        LocalCacheTTL:       2 * time.Minute,
//	    })
//
//	    ctx := context.Background()
//
//	    // Set a value with expiration
//	    err := cache.Set(ctx, "user:123", userData, 5*time.Minute)
//
//	    // Get a value
//	    var user User
//	    err := cache.Get(ctx, "user:123", &user)
//	    if err != nil {
//	        // Cache miss or error
//	    }
//
//	    // Delete from cache
//	    err := cache.Delete(ctx, "user:123")
//	}
//
// # Multiple Cache Instances
//
// The factory pattern allows creating multiple independent cache instances:
//
//	primaryRDB := redisclient.NewRDB(redisclient.Config{
//	    Urls: []string{"primary:6379"},
//	})
//	sessionRDB := redisclient.NewRDB(redisclient.Config{
//	    Urls: []string{"session:6379"},
//	})
//
//	primaryCache := redisclient.NewCache(redisclient.CacheConfig{
//	    Redis:               primaryRDB,
//	    RefreshEventChannel: "primary:refresh",
//	})
//	sessionCache := redisclient.NewCache(redisclient.CacheConfig{
//	    Redis:               sessionRDB,
//	    RefreshEventChannel: "session:refresh",
//	})
//
// # Distributed Cache Invalidation
//
// The cache automatically synchronizes invalidations across multiple instances:
//   - When Set() or Delete() is called, a cache refresh event is published
//   - All instances subscribed to the channel receive the event
//   - Local caches are automatically invalidated
//
// This ensures cache consistency in distributed deployments.
//
// # Working with go-redis
//
// The returned redis.UniversalClient supports all standard go-redis operations:
//
//	// Strings
//	rdb.Set(ctx, "key", "value", time.Hour)
//	rdb.Get(ctx, "key")
//
//	// Lists
//	rdb.LPush(ctx, "queue", "item1", "item2")
//	rdb.RPop(ctx, "queue")
//
//	// Sets
//	rdb.SAdd(ctx, "tags", "go", "redis")
//	rdb.SMembers(ctx, "tags")
//
//	// Hashes
//	rdb.HSet(ctx, "user:123", "name", "Alice")
//	rdb.HGetAll(ctx, "user:123")
//
//	// Sorted Sets
//	rdb.ZAdd(ctx, "leaderboard", redis.Z{Score: 100, Member: "player1"})
//	rdb.ZRange(ctx, "leaderboard", 0, 9)
//
//	// Pub/Sub
//	pubsub := rdb.Subscribe(ctx, "notifications")
//	ch := pubsub.Channel()
//
// For more go-redis features, see https://redis.uptrace.dev/
//
// # Configuration with Viper
//
// Example configuration file (configs/default.yaml):
//
//	redis:
//	  urls:
//	    - localhost:6379
//	  password: ""
//
// Loading configuration:
//
//	import (
//	    "github.com/poly-workshop/go-webmods/app"
//	    redisclient "github.com/poly-workshop/go-webmods/redisclient"
//	)
//
//	app.Init(".")
//	rdb := redisclient.NewRDB(redisclient.Config{
//	    Urls:     app.Config().GetStringSlice("redis.urls"),
//	    Password: app.Config().GetString("redis.password"),
//	})
//
// # Best Practices
//
//   - Use NewRDB and NewCache factory functions for creating clients
//   - Create separate Redis clients for different purposes (e.g., cache, sessions, queues)
//   - Use the cache for frequently accessed, slowly changing data
//   - Set appropriate expiration times to balance freshness and performance
//   - Use Redis Cluster for high availability and horizontal scaling
//   - Monitor cache hit rates and adjust local cache size if needed
//   - Handle cache misses gracefully (load from primary data source)
//   - Use different refresh event channels for independent cache instances
//
// # Error Handling
//
// - NewRDB panics if no Redis hosts are configured
// - NewCache panics if Redis client is nil
// - Cache operations return errors that should be handled:
//
//	err := cache.Get(ctx, key, &value)
//	if err != nil {
//	    // Handle cache miss or error
//	    // Load from database, etc.
//	}
package redisclient
