// Package mongo_client provides a factory function for creating MongoDB
// connections using the MongoDB Go Driver v2.
//
// # Basic Usage
//
// Create a MongoDB client:
//
//	import "github.com/poly-workshop/go-webmods/mongo_client"
//
//	client := mongo_client.NewClient(mongo_client.Config{
//	    URI:      "mongodb://localhost:27017",
//	    Database: "mydb",
//	})
//	defer client.Disconnect(context.Background())
//
// Create a MongoDB database directly:
//
//	db := mongo_client.NewDatabase(mongo_client.Config{
//	    URI:      "mongodb://localhost:27017",
//	    Database: "mydb",
//	})
//
// # Using with Viper Configuration
//
// The Config struct is designed to work seamlessly with Viper configuration:
//
//	import (
//	    "github.com/poly-workshop/go-webmods/app"
//	    "github.com/poly-workshop/go-webmods/mongo_client"
//	)
//
//	func main() {
//	    app.Init(".")
//
//	    cfg := mongo_client.Config{
//	        URI:            app.Config().GetString("mongodb.uri"),
//	        Database:       app.Config().GetString("mongodb.database"),
//	        ConnectTimeout: app.Config().GetDuration("mongodb.connect_timeout"),
//	        PingTimeout:    app.Config().GetDuration("mongodb.ping_timeout"),
//	    }
//
//	    db := mongo_client.NewDatabase(cfg)
//	    // Use db for MongoDB operations
//	}
//
// Example config file (configs/default.yaml):
//
//	mongodb:
//	  uri: mongodb://localhost:27017
//	  database: mydb
//	  connect_timeout: 10s
//	  ping_timeout: 5s
//
// # Working with MongoDB
//
// The returned *mongo.Client or *mongo.Database can be used with all standard MongoDB operations:
//
//	// Get a collection
//	collection := db.Collection("users")
//
//	// Insert a document
//	result, err := collection.InsertOne(context.TODO(), bson.D{
//	    {Key: "name", Value: "Alice"},
//	    {Key: "email", Value: "alice@example.com"},
//	})
//
//	// Find documents
//	cursor, err := collection.Find(context.TODO(), bson.D{{Key: "age", Value: bson.D{{Key: "$gt", Value: 18}}}})
//	defer cursor.Close(context.TODO())
//
//	for cursor.Next(context.TODO()) {
//	    var user User
//	    if err := cursor.Decode(&user); err != nil {
//	        log.Fatal(err)
//	    }
//	    fmt.Printf("User: %+v\n", user)
//	}
//
//	// Update a document
//	filter := bson.D{{Key: "name", Value: "Alice"}}
//	update := bson.D{{Key: "$set", Value: bson.D{{Key: "age", Value: 25}}}}
//	result, err := collection.UpdateOne(context.TODO(), filter, update)
//
//	// Delete a document
//	filter := bson.D{{Key: "name", Value: "Alice"}}
//	result, err := collection.DeleteOne(context.TODO(), filter)
//
// For more MongoDB features, see https://www.mongodb.com/docs/drivers/go/current/
//
// # Connection Strings
//
// MongoDB connection strings support various options:
//
//	// Simple connection
//	mongodb://localhost:27017
//
//	// With authentication
//	mongodb://username:password@localhost:27017
//
//	// MongoDB Atlas
//	mongodb+srv://username:password@cluster.mongodb.net/
//
//	// With options
//	mongodb://localhost:27017/?maxPoolSize=20&w=majority
//
//	// Replica set
//	mongodb://node1:27017,node2:27017,node3:27017/?replicaSet=rs0
//
// # Error Handling
//
// NewClient and NewDatabase panic if:
//   - Connection to MongoDB fails
//   - Ping operation times out or fails
//
// In production, consider recovering from panics or validating configuration
// before calling these functions.
//
// # Timeouts and Connection Pooling
//
// Default timeouts:
//   - ConnectTimeout: 10 seconds (if not specified)
//   - PingTimeout: 5 seconds (if not specified)
//
// For production use, configure connection pooling and timeouts:
//
//	clientOptions := options.Client().
//	    ApplyURI(cfg.URI).
//	    SetMaxPoolSize(100).
//	    SetMinPoolSize(10).
//	    SetMaxConnIdleTime(30 * time.Second)
//
// # Best Practices
//
//   - Use connection strings from configuration or environment variables
//   - Always close the client when done: defer client.Disconnect(context.Background())
//   - Use context with timeout for all database operations
//   - Configure appropriate connection pool sizes for your workload
//   - Use MongoDB Atlas for production deployments
//   - Enable authentication and TLS for production connections
//   - Use indexes to optimize query performance
//   - Monitor connection pool metrics in production
package mongo_client
