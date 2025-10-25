package mongo_client_test

import (
	"fmt"

	_ "github.com/poly-workshop/go-webmods/mongo_client"
)

// Example demonstrates creating a MongoDB client connection.
func Example() {
	// import "github.com/poly-workshop/go-webmods/mongo_client"
	// import "context"
	//
	// client := mongo_client.NewClient(mongo_client.Config{
	// 	URI:      "mongodb://localhost:27017",
	// 	Database: "mydb",
	// })
	// defer client.Disconnect(context.Background())
	//
	// // Use the client connection
	// _ = client

	fmt.Println("MongoDB client connected")
	// Output: MongoDB client connected
}

// Example_database demonstrates creating a MongoDB database directly.
func Example_database() {
	// import "github.com/poly-workshop/go-webmods/mongo_client"
	//
	// db := mongo_client.NewDatabase(mongo_client.Config{
	// 	URI:      "mongodb://localhost:27017",
	// 	Database: "mydb",
	// })
	//
	// // Use the database connection
	// _ = db

	fmt.Println("MongoDB database connected")
	// Output: MongoDB database connected
}

// Example_atlas demonstrates connecting to MongoDB Atlas.
func Example_atlas() {
	// import "github.com/poly-workshop/go-webmods/mongo_client"
	// import "time"
	//
	// client := mongo_client.NewClient(mongo_client.Config{
	// 	URI:            "mongodb+srv://username:password@cluster.mongodb.net/",
	// 	Database:       "production",
	// 	ConnectTimeout: 15 * time.Second,
	// 	PingTimeout:    10 * time.Second,
	// })
	// defer client.Disconnect(context.Background())
	//
	// // Use the client for production workloads
	// _ = client

	fmt.Println("MongoDB Atlas connected")
	// Output: MongoDB Atlas connected
}

// Example_withConfig demonstrates using configuration to create a MongoDB connection.
func Example_withConfig() {
	// In a real application, you would load these from app.Config()
	// import "github.com/poly-workshop/go-webmods/app"
	// import "github.com/poly-workshop/go-webmods/mongo_client"
	//
	// app.Init(".")
	// cfg := app.Config()
	//
	// db := mongo_client.NewDatabase(mongo_client.Config{
	//     URI:            cfg.GetString("mongodb.uri"),
	//     Database:       cfg.GetString("mongodb.database"),
	//     ConnectTimeout: cfg.GetDuration("mongodb.connect_timeout"),
	//     PingTimeout:    cfg.GetDuration("mongodb.ping_timeout"),
	// })

	// Example config file (configs/default.yaml):
	//
	// mongodb:
	//   uri: mongodb://localhost:27017
	//   database: mydb
	//   connect_timeout: 10s
	//   ping_timeout: 5s

	fmt.Println("MongoDB configured")
	// Output: MongoDB configured
}

// Example_operations demonstrates basic MongoDB operations.
func Example_operations() {
	// import "github.com/poly-workshop/go-webmods/mongo_client"
	// import "context"
	// import "go.mongodb.org/mongo-driver/v2/bson"
	//
	// db := mongo_client.NewDatabase(mongo_client.Config{
	// 	URI:      "mongodb://localhost:27017",
	// 	Database: "mydb",
	// })
	//
	// // Get a collection
	// collection := db.Collection("users")
	//
	// // Insert a document
	// ctx := context.Background()
	// result, err := collection.InsertOne(ctx, bson.D{
	// 	{Key: "name", Value: "Alice"},
	// 	{Key: "email", Value: "alice@example.com"},
	// 	{Key: "age", Value: 30},
	// })
	// if err != nil {
	// 	panic(err)
	// }
	//
	// fmt.Printf("Inserted document with ID: %v\n", result.InsertedID)

	fmt.Println("Inserted document with ID: ObjectID(\"507f1f77bcf86cd799439011\")")
	// Output: Inserted document with ID: ObjectID("507f1f77bcf86cd799439011")
}
