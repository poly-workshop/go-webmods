// Package redis_client provides Redis client factories with support for both
// single-node and cluster modes, plus an integrated caching layer with local
// memory cache and distributed cache invalidation.
//
// # Factory Pattern (Recommended)
//
// Create Redis clients using the factory pattern to support multiple independent clients:
//
//	import "github.com/poly-workshop/go-webmods/redis_client"
//
//	func main() {
//	    // Create a Redis client
//	    rdb := redis_client.NewRDB(redis_client.Config{
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
//	primaryRDB := redis_client.NewRDB(redis_client.Config{
//	    Urls:     []string{"primary:6379"},
//	    Password: "secret1",
//	})
//
//	cacheRDB := redis_client.NewRDB(redis_client.Config{
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
//	rdb := redis_client.NewRDB(redis_client.Config{
//	    Urls:     []string{"node1:6379", "node2:6379", "node3:6379"},
//	    Password: "password",
//	})
//
// # Singleton Pattern (Deprecated)
//
// The legacy singleton pattern is still supported for backward compatibility:
//
//	redis_client.SetConfig([]string{"localhost:6379"}, "")
//	rdb := redis_client.GetRDB()
//
// However, this pattern is deprecated. Use NewRDB instead for better control
// and to support multiple clients.
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
//	import "github.com/poly-workshop/go-webmods/redis_client"
//
//	func main() {
//	    rdb := redis_client.NewRDB(redis_client.Config{
//	        Urls:     []string{"localhost:6379"},
//	        Password: "",
//	    })
//
//	    cache := redis_client.NewCache(redis_client.CacheConfig{
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
// Singleton pattern cache creation (deprecated):
//
//	redis_client.SetConfig([]string{"localhost:6379"}, "")
//	cache := redis_client.GetCache()
//
// # Multiple Cache Instances
//
// The factory pattern allows creating multiple independent cache instances:
//
//	primaryRDB := redis_client.NewRDB(redis_client.Config{
//	    Urls: []string{"primary:6379"},
//	})
//	sessionRDB := redis_client.NewRDB(redis_client.Config{
//	    Urls: []string{"session:6379"},
//	})
//
//	primaryCache := redis_client.NewCache(redis_client.CacheConfig{
//	    Redis:               primaryRDB,
//	    RefreshEventChannel: "primary:refresh",
//	})
//	sessionCache := redis_client.NewCache(redis_client.CacheConfig{
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
// Customize the cache refresh channel name (must be done before GetCache):
//
//	redis_client.SetCacheRefreshEventChannel("myapp:cache:refresh")
//	cache := redis_client.GetCache()
//
// # Singleton Pattern (Deprecated)
//
// The legacy GetRDB() and GetCache() functions use the singleton pattern:
//   - First call initializes the client/cache
//   - Subsequent calls return the same instance
//   - Thread-safe initialization using sync.Once
//
// This means SetConfig() and SetCacheRefreshEventChannel() must be called
// before the first GetRDB() or GetCache() call.
//
// For new code, prefer using NewRDB and NewCache instead for better control
// and to support multiple independent clients.
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
//	import "github.com/poly-workshop/go-webmods/app"
//
//	app.Init(".")
//	redis_client.SetConfig(
//	    app.Config().GetStringSlice("redis.urls"),
//	    app.Config().GetString("redis.password"),
//	)
//
// # Best Practices
//
//   - Use NewRDB and NewCache factory functions for new code (not the deprecated singleton pattern)
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
// - GetRDB panics if no Redis hosts are configured (deprecated singleton pattern)
// - Cache operations return errors that should be handled:
//
//	err := cache.Get(ctx, key, &value)
//	if err != nil {
//	    // Handle cache miss or error
//	    // Load from database, etc.
//	}
package redis_client
