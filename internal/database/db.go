package database

import (
	"log"
	"os"
	"sync"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB
var once sync.Once

func Get() *gorm.DB {
	once.Do(func() {
		_ = godotenv.Load()

		dsn := os.Getenv("DATABASE_URL")
		var err error
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			log.Fatal("Failed to connect to database:", err)
		}
	})

	return db
}
