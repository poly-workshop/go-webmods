package gorm_client

import (
	"fmt"
	"os"
	"path/filepath"

	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Config struct {
	Driver   string
	Host     string
	Port     int
	Username string
	Password string
	Name     string
	SSLMode  string
}

func NewDB(cfg Config) *gorm.DB {
	driver := cfg.Driver
	switch driver {
	case "postgres":
		db, err := openPostgres(cfg)
		if err != nil {
			panic(err)
		}
		return db
	case "sqlite":
		db, err := openSqlite(cfg)
		if err != nil {
			panic(err)
		}
		return db
	default:
		panic(fmt.Sprintf("unsupported database driver: %s", driver))
	}
}

func openPostgres(cfg Config) (db *gorm.DB, err error) {
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s dbname=%s password=%s sslmode=%s",
		cfg.Host,
		cfg.Port,
		cfg.Username,
		cfg.Name,
		cfg.Password,
		cfg.SSLMode,
	)
	db, err = gorm.Open(postgres.Open(dsn))
	if err != nil {
		return nil, err
	}
	return db, nil
}

func openSqlite(cfg Config) (db *gorm.DB, err error) {
	// Ensure directory exists for SQLite database file
	dbPath := cfg.Name
	if dir := filepath.Dir(dbPath); dir != "." {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return nil, fmt.Errorf("failed to create directory for SQLite database: %w", err)
		}
	}
	
	db, err = gorm.Open(sqlite.Open(dbPath))
	if err != nil {
		return nil, err
	}
	return db, nil
}
