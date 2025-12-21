package mongo_client

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type Config struct {
	URI            string
	Database       string
	ConnectTimeout time.Duration
	PingTimeout    time.Duration
}

func NewClient(cfg Config) *mongo.Client {
	// Set default timeouts if not provided
	if cfg.ConnectTimeout == 0 {
		cfg.ConnectTimeout = 10 * time.Second
	}
	if cfg.PingTimeout == 0 {
		cfg.PingTimeout = 5 * time.Second
	}

	// Create MongoDB client options
	clientOptions := options.Client().ApplyURI(cfg.URI)

	// Connect to MongoDB
	client, err := mongo.Connect(clientOptions)
	if err != nil {
		panic(fmt.Sprintf("failed to connect to MongoDB: %v", err))
	}

	// Ping to verify connection
	pingCtx, pingCancel := context.WithTimeout(context.Background(), cfg.PingTimeout)
	defer pingCancel()

	if err := client.Ping(pingCtx, nil); err != nil {
		panic(fmt.Sprintf("failed to ping MongoDB: %v", err))
	}

	return client
}

// NewDatabase creates a new MongoDB database connection.
// Note: This function returns only the database reference. The underlying client
// connection will remain open for the lifetime of the application. This is the
// intended usage for long-running applications (servers, daemons).
//
// If you need to explicitly manage the client lifecycle (e.g., for short-lived
// scripts or need to call Disconnect), use NewClient instead and access the
// database via client.Database(name).
func NewDatabase(cfg Config) *mongo.Database {
	client := NewClient(cfg)
	return client.Database(cfg.Database)
}
