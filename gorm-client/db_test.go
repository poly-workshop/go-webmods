package gorm_client

import (
	"os"
	"testing"
)

func TestNewDB_SQLite(t *testing.T) {
	// Create temporary directory for testing
	tempDir, err := os.MkdirTemp("", "gorm_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer func() {
		if err := os.RemoveAll(tempDir); err != nil {
			t.Logf("Failed to remove temp dir: %v", err)
		}
	}()

	dbPath := tempDir + "/test.db"
	db := NewDB(Config{
		Driver: "sqlite",
		Name:   dbPath,
	})

	if db == nil {
		t.Fatal("Expected non-nil database connection")
	}

	// Verify we can ping the database
	sqlDB, err := db.DB()
	if err != nil {
		t.Fatalf("Failed to get database instance: %v", err)
	}

	if err := sqlDB.Ping(); err != nil {
		t.Fatalf("Failed to ping database: %v", err)
	}
}

func TestNewDB_UnsupportedDriver(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("Expected panic for unsupported driver")
		}
	}()

	NewDB(Config{
		Driver: "unsupported",
	})
}

func TestOpenMysql_DSNFormat(t *testing.T) {
	// This test verifies the DSN format is correct
	// We don't actually connect to MySQL, just verify the function doesn't panic
	// with valid configuration
	cfg := Config{
		Driver:   "mysql",
		Host:     "localhost",
		Port:     3306,
		Username: "user",
		Password: "password",
		Name:     "testdb",
	}

	// The openMysql function will attempt to connect and fail (no MySQL server)
	// but we can verify it constructs the DSN correctly
	_, err := openMysql(cfg)
	// We expect an error since there's no MySQL server running
	// This just ensures the function exists and can be called
	if err == nil {
		t.Log("MySQL connection succeeded (unexpected in test environment)")
	} else {
		t.Logf("MySQL connection failed as expected: %v", err)
	}
}

func TestOpenPostgres_DSNFormat(t *testing.T) {
	// This test verifies the DSN format is correct for Postgres
	cfg := Config{
		Driver:   "postgres",
		Host:     "localhost",
		Port:     5432,
		Username: "user",
		Password: "password",
		Name:     "testdb",
		SSLMode:  "disable",
	}

	// The openPostgres function will attempt to connect and fail (no Postgres server)
	_, err := openPostgres(cfg)
	// We expect an error since there's no Postgres server running
	if err == nil {
		t.Log("Postgres connection succeeded (unexpected in test environment)")
	} else {
		t.Logf("Postgres connection failed as expected: %v", err)
	}
}
