package redis_client_test

import (
	"context"
	"testing"
	"time"

	"github.com/poly-workshop/go-webmods/redis_client"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

// startRedisContainer starts a Redis container for testing.
func startRedisContainer(t *testing.T) (string, func()) {
	t.Helper()

	ctx := context.Background()
	req := testcontainers.ContainerRequest{
		Image:        "redis:7-alpine",
		ExposedPorts: []string{"6379/tcp"},
		WaitingFor:   wait.ForLog("Ready to accept connections"),
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		t.Fatalf("Failed to start Redis container: %v", err)
	}

	host, err := container.Host(ctx)
	if err != nil {
		t.Fatalf("Failed to get container host: %v", err)
	}

	port, err := container.MappedPort(ctx, "6379")
	if err != nil {
		t.Fatalf("Failed to get container port: %v", err)
	}

	addr := host + ":" + port.Port()

	cleanup := func() {
		if err := container.Terminate(ctx); err != nil {
			t.Errorf("Failed to terminate container: %v", err)
		}
	}

	return addr, cleanup
}

func TestNewRDB(t *testing.T) {
	addr, cleanup := startRedisContainer(t)
	defer cleanup()

	// Test creating a single-node client
	t.Run("SingleNode", func(t *testing.T) {
		rdb := redis_client.NewRDB(redis_client.Config{
			Urls:     []string{addr},
			Password: "",
		})

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// Test basic operations
		err := rdb.Set(ctx, "test_key", "test_value", time.Minute).Err()
		if err != nil {
			t.Fatalf("Failed to set key: %v", err)
		}

		val, err := rdb.Get(ctx, "test_key").Result()
		if err != nil {
			t.Fatalf("Failed to get key: %v", err)
		}

		if val != "test_value" {
			t.Errorf("Expected value 'test_value', got '%s'", val)
		}
	})

	// Test creating multiple independent clients
	t.Run("MultipleClients", func(t *testing.T) {
		rdb1 := redis_client.NewRDB(redis_client.Config{
			Urls:     []string{addr},
			Password: "",
		})

		rdb2 := redis_client.NewRDB(redis_client.Config{
			Urls:     []string{addr},
			Password: "",
		})

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// Set value with client 1
		err := rdb1.Set(ctx, "client1_key", "value1", time.Minute).Err()
		if err != nil {
			t.Fatalf("Failed to set key with client1: %v", err)
		}

		// Set value with client 2
		err = rdb2.Set(ctx, "client2_key", "value2", time.Minute).Err()
		if err != nil {
			t.Fatalf("Failed to set key with client2: %v", err)
		}

		// Read value set by client 1 using client 2
		val, err := rdb2.Get(ctx, "client1_key").Result()
		if err != nil {
			t.Fatalf("Failed to get key with client2: %v", err)
		}
		if val != "value1" {
			t.Errorf("Expected value 'value1', got '%s'", val)
		}

		// Read value set by client 2 using client 1
		val, err = rdb1.Get(ctx, "client2_key").Result()
		if err != nil {
			t.Fatalf("Failed to get key with client1: %v", err)
		}
		if val != "value2" {
			t.Errorf("Expected value 'value2', got '%s'", val)
		}
	})

	// Test panic on empty URLs
	t.Run("EmptyURLs", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("Expected panic for empty URLs, but didn't panic")
			}
		}()

		redis_client.NewRDB(redis_client.Config{
			Urls:     []string{},
			Password: "",
		})
	})
}

func TestGetRDB_BackwardCompatibility(t *testing.T) {
	addr, cleanup := startRedisContainer(t)
	defer cleanup()

	// Test backward compatibility with singleton pattern
	redis_client.SetConfig([]string{addr}, "")
	rdb := redis_client.GetRDB()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := rdb.Set(ctx, "singleton_key", "singleton_value", time.Minute).Err()
	if err != nil {
		t.Fatalf("Failed to set key: %v", err)
	}

	val, err := rdb.Get(ctx, "singleton_key").Result()
	if err != nil {
		t.Fatalf("Failed to get key: %v", err)
	}

	if val != "singleton_value" {
		t.Errorf("Expected value 'singleton_value', got '%s'", val)
	}

	// Test that subsequent calls return the same instance
	rdb2 := redis_client.GetRDB()
	if rdb != rdb2 {
		t.Error("GetRDB should return the same singleton instance")
	}
}
