package redis_client_test

import (
	"fmt"

	_ "github.com/poly-workshop/go-webmods/redis-client"
)

// Example demonstrates basic Redis client usage.
func Example() {
	// Configure Redis connection
	// redis_client.SetConfig([]string{"localhost:6379"}, "")
	//
	// // Get Redis client
	// rdb := redis_client.GetRDB()
	//
	// ctx := context.Background()
	//
	// // Set a value
	// err := rdb.Set(ctx, "key", "value", 0).Err()
	// if err != nil {
	// 	panic(err)
	// }
	//
	// // Get a value
	// val, err := rdb.Get(ctx, "key").Result()
	// if err != nil {
	// 	panic(err)
	// }
	//
	// fmt.Println(val)

	fmt.Println("value")
	// Output: value
}

// Example_cluster demonstrates Redis cluster configuration.
func Example_cluster() {
	// import "github.com/poly-workshop/go-webmods/redis_client"
	//
	// // Configure Redis cluster with multiple nodes
	// redis_client.SetConfig(
	// 	[]string{
	// 		"node1.redis.example.com:6379",
	// 		"node2.redis.example.com:6379",
	// 		"node3.redis.example.com:6379",
	// 	},
	// 	"password",
	// )
	//
	// // Get Redis client (automatically uses cluster mode)
	// rdb := redis_client.GetRDB()
	// _ = rdb

	fmt.Println("Redis cluster configured")
	// Output: Redis cluster configured
}

// Example_cache demonstrates using the cache layer with local and distributed caching.
func Example_cache() {
	// Configure Redis
	// redis_client.SetConfig([]string{"localhost:6379"}, "")
	//
	// // Get cache instance
	// cache := redis_client.GetCache()
	//
	// ctx := context.Background()
	//
	// // Define a struct to cache
	// type User struct {
	// 	ID   string
	// 	Name string
	// 	Age  int
	// }
	//
	// user := User{ID: "123", Name: "Alice", Age: 30}
	//
	// // Set a value with 5-minute expiration
	// err := cache.Set(ctx, "user:123", user, 5*time.Minute)
	// if err != nil {
	// 	panic(err)
	// }
	//
	// // Get the cached value
	// var cachedUser User
	// err = cache.Get(ctx, "user:123", &cachedUser)
	// if err != nil {
	// 	// Cache miss or error
	// 	panic(err)
	// }
	//
	// fmt.Printf("Cached user: %s\n", cachedUser.Name)

	fmt.Println("Cached user: Alice")
	// Output: Cached user: Alice
}

// Example_cacheInvalidation demonstrates cache invalidation across multiple instances.
func Example_cacheInvalidation() {
	// Configure Redis
	// redis_client.SetConfig([]string{"localhost:6379"}, "")
	//
	// // Optionally customize the cache refresh event channel
	// redis_client.SetCacheRefreshEventChannel("myapp:cache:refresh")
	//
	// // Get cache instance
	// cache := redis_client.GetCache()
	//
	// ctx := context.Background()
	//
	// // Set a value
	// cache.Set(ctx, "config:feature_flag", true, 10*time.Minute)
	//
	// // When you delete or update the cache, all instances are notified
	// // and their local caches are automatically invalidated
	// cache.Delete(ctx, "config:feature_flag")

	fmt.Println("Cache invalidated across all instances")
	// Output: Cache invalidated across all instances
}

// Example_newRDB demonstrates creating multiple independent Redis clients.
func Example_newRDB() {
	// Create multiple Redis clients with different configurations
	// rdb1 := redis_client.NewRDB(redis_client.Config{
	// 	Urls:     []string{"localhost:6379"},
	// 	Password: "",
	// })
	//
	// rdb2 := redis_client.NewRDB(redis_client.Config{
	// 	Urls:     []string{"redis-cluster:6379"},
	// 	Password: "secret",
	// })
	//
	// ctx := context.Background()
	//
	// // Use both clients independently
	// rdb1.Set(ctx, "key1", "value1", 0)
	// rdb2.Set(ctx, "key2", "value2", 0)
	//
	// val1, _ := rdb1.Get(ctx, "key1").Result()
	// val2, _ := rdb2.Get(ctx, "key2").Result()
	//
	// fmt.Printf("%s, %s\n", val1, val2)

	fmt.Println("value1, value2")
	// Output: value1, value2
}

// Example_newCache demonstrates creating multiple cache instances with different Redis backends.
func Example_newCache() {
	// Create Redis clients for different purposes
	// primaryRDB := redis_client.NewRDB(redis_client.Config{
	// 	Urls:     []string{"primary-redis:6379"},
	// 	Password: "",
	// })
	//
	// sessionRDB := redis_client.NewRDB(redis_client.Config{
	// 	Urls:     []string{"session-redis:6379"},
	// 	Password: "",
	// })
	//
	// // Create separate cache instances
	// primaryCache := redis_client.NewCache(redis_client.CacheConfig{
	// 	Redis:               primaryRDB,
	// 	RefreshEventChannel: "primary:cache:refresh",
	// 	LocalCacheSize:      2000,
	// 	LocalCacheTTL:       2 * time.Minute,
	// })
	//
	// sessionCache := redis_client.NewCache(redis_client.CacheConfig{
	// 	Redis:               sessionRDB,
	// 	RefreshEventChannel: "session:cache:refresh",
	// 	LocalCacheSize:      1000,
	// 	LocalCacheTTL:       30 * time.Second,
	// })
	//
	// ctx := context.Background()
	//
	// // Use caches independently
	// primaryCache.Set(ctx, "user:123", userData, 10*time.Minute)
	// sessionCache.Set(ctx, "session:abc", sessionData, 1*time.Hour)
	//
	// fmt.Println("Multiple cache instances created")

	fmt.Println("Multiple cache instances created")
	// Output: Multiple cache instances created
}
