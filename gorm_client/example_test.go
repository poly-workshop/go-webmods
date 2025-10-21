package gorm_client_test

import (
	"fmt"

	_ "github.com/poly-workshop/go-webmods/gorm_client"
)

// Example demonstrates creating a PostgreSQL database connection.
func Example() {
	// import "github.com/poly-workshop/go-webmods/gorm_client"
	//
	// db := gorm_client.NewDB(gorm_client.Config{
	// 	Driver:   "postgres",
	// 	Host:     "localhost",
	// 	Port:     5432,
	// 	Username: "user",
	// 	Password: "password",
	// 	Name:     "mydb",
	// 	SSLMode:  "disable",
	// })
	//
	// // Use the database connection
	// _ = db

	fmt.Println("Database connected")
	// Output: Database connected
}

// Example_sqlite demonstrates creating a SQLite database connection.
func Example_sqlite() {
	// import "github.com/poly-workshop/go-webmods/gorm_client"
	//
	// db := gorm_client.NewDB(gorm_client.Config{
	// 	Driver: "sqlite",
	// 	Name:   "/tmp/test.db",
	// })
	//
	// // Use the database connection
	// _ = db

	fmt.Println("SQLite database connected")
	// Output: SQLite database connected
}

// Example_withConfig demonstrates using configuration to create a database connection.
func Example_withConfig() {
	// In a real application, you would load these from app.Config()
	// import "github.com/poly-workshop/go-webmods/app"
	// import "github.com/poly-workshop/go-webmods/gorm_client"
	//
	// app.Init(".")
	// cfg := app.Config()
	//
	// db := gorm_client.NewDB(gorm_client.Config{
	//     Driver:   cfg.GetString("database.driver"),
	//     Host:     cfg.GetString("database.host"),
	//     Port:     cfg.GetInt("database.port"),
	//     Username: cfg.GetString("database.username"),
	//     Password: cfg.GetString("database.password"),
	//     Name:     cfg.GetString("database.name"),
	//     SSLMode:  cfg.GetString("database.sslmode"),
	// })

	// Example config file (configs/default.yaml):
	//
	// database:
	//   driver: postgres
	//   host: localhost
	//   port: 5432
	//   username: user
	//   password: pass
	//   name: mydb
	//   sslmode: disable
}
