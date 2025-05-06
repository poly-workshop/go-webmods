package gorm

import (
	"fmt"
	"sync"

	"github.com/oj-lab/go-webmods/app"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

const (
	configKeyDatabaseDriver = "gorm.database.driver"
	configKeyDatabaseHost   = "gorm.database.host"
	configKeyDatabasePort   = "gorm.database.port"
	configKeyDatabaseUser   = "gorm.database.username"
	configKeyDatabaseName   = "gorm.database.name"
	configKeyDatabasePass   = "gorm.database.password"
	configKeyDatabaseSSL    = "gorm.database.sslmode"
)

var (
	initMutx sync.Mutex
	db       *gorm.DB
)

func GetDB() *gorm.DB {
	if db == nil {
		initMutx.Lock()
		defer initMutx.Unlock()
		if db != nil {
			return db
		}
		var err error
		driver := app.Config().GetString(configKeyDatabaseDriver)
		switch driver {
		case "postgres":
			db, err = openPostgres()
			if err != nil {
				panic(err)
			}
		case "sqlite":
			db, err = openSqlite()
			if err != nil {
				panic(err)
			}
		default:
			panic(fmt.Sprintf("unsupported database driver: %s", driver))
		}
	}
	return db
}

func openPostgres() (db *gorm.DB, err error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		app.Config().GetString(configKeyDatabaseHost),
		app.Config().GetString(configKeyDatabasePort),
		app.Config().GetString(configKeyDatabaseUser),
		app.Config().GetString(configKeyDatabaseName),
		app.Config().GetString(configKeyDatabasePass),
		app.Config().GetString(configKeyDatabaseSSL),
	)
	db, err = gorm.Open(postgres.Open(dsn))
	if err != nil {
		return nil, err
	}
	return db, nil
}

func openSqlite() (db *gorm.DB, err error) {
	db, err = gorm.Open(sqlite.Open(app.Config().GetString(configKeyDatabaseName)))
	if err != nil {
		return nil, err
	}
	return db, nil
}
