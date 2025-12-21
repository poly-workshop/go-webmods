// Package gormclient provides a factory function for creating GORM database
// connections with support for multiple database drivers.
//
// # Supported Databases
//
//   - PostgreSQL: Production-grade relational database
//   - SQLite: Lightweight file-based database for development and testing
//
// # Basic Usage
//
// Create a PostgreSQL connection:
//
//	import gormclient "github.com/poly-workshop/go-webmods/gormclient"
//
//	db := gormclient.NewDB(gormclient.Config{
//	    Driver:   "postgres",
//	    Host:     "localhost",
//	    Port:     5432,
//	    Username: "user",
//	    Password: "password",
//	    DbName:   "mydb",
//	    SSLMode:  "disable",
//	})
//
// Create a SQLite connection:
//
//	db := gormclient.NewDB(gormclient.Config{
//	    Driver: "sqlite",
//	    DbName: "data/app.db",  // File path for SQLite database
//	})
//
// The function automatically creates parent directories for SQLite databases.
//
// # Using with Viper Configuration
//
// The Config struct is designed to work seamlessly with Viper configuration:
//
//	import (
//	    "github.com/poly-workshop/go-webmods/app"
//	    gormclient "github.com/poly-workshop/go-webmods/gormclient"
//	)
//
//	func main() {
//	    app.Init(".")
//
//	    cfg := gormclient.Config{
//	        Driver:   app.Config().GetString("database.driver"),
//	        Host:     app.Config().GetString("database.host"),
//	        Port:     app.Config().GetInt("database.port"),
//	        Username: app.Config().GetString("database.username"),
//	        Password: app.Config().GetString("database.password"),
//	        DbName:   app.Config().GetString("database.name"),
//	        SSLMode:  app.Config().GetString("database.sslmode"),
//	    }
//
//	    db := gormclient.NewDB(cfg)
//	    // Use db for GORM operations
//	}
//
// Example config file (configs/default.yaml):
//
//	database:
//	  driver: postgres
//	  host: localhost
//	  port: 5432
//	  username: myuser
//	  password: mypass
//	  name: mydb
//	  sslmode: disable
//
// # Working with GORM
//
// The returned *gorm.DB can be used with all standard GORM operations:
//
//	// Auto-migrate tables
//	db.AutoMigrate(&User{}, &Post{})
//
//	// Create records
//	db.Create(&User{Name: "Alice", Email: "alice@example.com"})
//
//	// Query records
//	var users []User
//	db.Where("age > ?", 18).Find(&users)
//
//	// Update records
//	db.Model(&user).Update("Age", 25)
//
//	// Delete records
//	db.Delete(&user)
//
// For more GORM features, see https://gorm.io/docs/
//
// # Error Handling
//
// NewDB panics if:
//   - An unsupported driver is specified
//   - Database connection fails
//   - Directory creation fails for SQLite
//
// In production, consider recovering from panics or validating configuration
// before calling NewDB.
//
// # Connection Pooling
//
// For production use, configure connection pooling:
//
//	sqlDB, err := db.DB()
//	if err != nil {
//	    panic(err)
//	}
//	sqlDB.SetMaxOpenConns(25)
//	sqlDB.SetMaxIdleConns(5)
//	sqlDB.SetConnMaxLifetime(5 * time.Minute)
//
// # Best Practices
//
//   - Use PostgreSQL for production environments
//   - Use SQLite for development, testing, and small deployments
//   - Store database credentials in environment variables or secure vaults
//   - Enable SSL mode for PostgreSQL in production (sslmode: require)
//   - Configure connection pooling for high-traffic applications
//   - Use migrations for schema management (e.g., golang-migrate or GORM AutoMigrate)
package gormclient
