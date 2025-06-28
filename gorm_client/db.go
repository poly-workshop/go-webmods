package gorm_client

import (
	"fmt"

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
	db, err = gorm.Open(sqlite.Open(cfg.Name))
	if err != nil {
		return nil, err
	}
	return db, nil
}
