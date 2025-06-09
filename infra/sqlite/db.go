package sqlite

import (
	"log"
	"sync"

	"github.com/ToshihiroOgino/elib/env"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	db     *gorm.DB
	dbOnce sync.Once
)

type Config struct {
	DBPath string
}

func DefaultConfig() Config {
	return Config{
		DBPath: env.Get().DBFile,
	}
}

func GetDB() *gorm.DB {
	dbOnce.Do(func() {
		var err error
		config := DefaultConfig()

		db, err = gorm.Open(sqlite.Open(config.DBPath), &gorm.Config{})
		if err != nil {
			log.Fatalf("Failed to connect to database: %v", err)
		}

		log.Println("Database connection established")
	})

	return db
}

func CloseDB() error {
	if db != nil {
		sqlDB, err := db.DB()
		if err != nil {
			return err
		}
		return sqlDB.Close()
	}
	return nil
}
