package redis_client_test

import (
	"context"
	"testing"
	"time"

	"github.com/poly-workshop/go-webmods/redis_client"
)

func TestNewCache(t *testing.T) {
	addr, cleanup := startRedisContainer(t)
	defer cleanup()

	t.Run("BasicCacheOperations", func(t *testing.T) {
		rdb := redis_client.NewRDB(redis_client.Config{
			Urls:     []string{addr},
			Password: "",
		})

		cache := redis_client.NewCache(redis_client.CacheConfig{
			Redis: rdb,
		})

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// Test Set and Get
		type TestStruct struct {
			Name string
			Age  int
		}

		testData := TestStruct{Name: "Alice", Age: 30}
		err := cache.Set(ctx, "user:123", testData, 5*time.Minute)
		if err != nil {
			t.Fatalf("Failed to set cache: %v", err)
		}

		var result TestStruct
		err = cache.Get(ctx, "user:123", &result)
		if err != nil {
			t.Fatalf("Failed to get cache: %v", err)
		}

		if result.Name != testData.Name || result.Age != testData.Age {
			t.Errorf("Expected %+v, got %+v", testData, result)
		}

		// Test Delete
		err = cache.Delete(ctx, "user:123")
		if err != nil {
			t.Fatalf("Failed to delete cache: %v", err)
		}

		// Verify deletion
		err = cache.Get(ctx, "user:123", &result)
		if err == nil {
			t.Error("Expected error after deletion, got nil")
		}
	})

	t.Run("MultipleCacheInstances", func(t *testing.T) {
		rdb1 := redis_client.NewRDB(redis_client.Config{
			Urls:     []string{addr},
			Password: "",
		})

		rdb2 := redis_client.NewRDB(redis_client.Config{
			Urls:     []string{addr},
			Password: "",
		})

		cache1 := redis_client.NewCache(redis_client.CacheConfig{
			Redis:               rdb1,
			RefreshEventChannel: "cache1:refresh",
		})

		cache2 := redis_client.NewCache(redis_client.CacheConfig{
			Redis:               rdb2,
			RefreshEventChannel: "cache2:refresh",
		})

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// Set value in cache1
		err := cache1.Set(ctx, "shared:key", "value1", 5*time.Minute)
		if err != nil {
			t.Fatalf("Failed to set in cache1: %v", err)
		}

		// Read from cache2 (should be able to access the same Redis backend)
		var result string
		err = cache2.Get(ctx, "shared:key", &result)
		if err != nil {
			t.Fatalf("Failed to get from cache2: %v", err)
		}

		if result != "value1" {
			t.Errorf("Expected 'value1', got '%s'", result)
		}

		// Set value in cache2 using a different key to avoid local cache conflicts
		err = cache2.Set(ctx, "cache2:key", "value2", 5*time.Minute)
		if err != nil {
			t.Fatalf("Failed to set in cache2: %v", err)
		}

		// Read from cache1 using the key set by cache2
		err = cache1.Get(ctx, "cache2:key", &result)
		if err != nil {
			t.Fatalf("Failed to get from cache1: %v", err)
		}

		if result != "value2" {
			t.Errorf("Expected 'value2', got '%s'", result)
		}
	})

	t.Run("CustomConfiguration", func(t *testing.T) {
		rdb := redis_client.NewRDB(redis_client.Config{
			Urls:     []string{addr},
			Password: "",
		})

		cache := redis_client.NewCache(redis_client.CacheConfig{
			Redis:               rdb,
			RefreshEventChannel: "custom:channel",
			LocalCacheSize:      500,
			LocalCacheTTL:       30 * time.Second,
		})

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// Test basic operation to ensure custom config works
		err := cache.Set(ctx, "test:key", "test:value", time.Minute)
		if err != nil {
			t.Fatalf("Failed to set cache with custom config: %v", err)
		}

		var result string
		err = cache.Get(ctx, "test:key", &result)
		if err != nil {
			t.Fatalf("Failed to get cache with custom config: %v", err)
		}

		if result != "test:value" {
			t.Errorf("Expected 'test:value', got '%s'", result)
		}
	})

	t.Run("PanicOnNilRedisClient", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("Expected panic for nil Redis client, but didn't panic")
			}
		}()

		redis_client.NewCache(redis_client.CacheConfig{
			Redis: nil,
		})
	})
}
