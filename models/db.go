package models

import (
	"os"
	"sync"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var once sync.Once
var database *gorm.DB

func GetDB() *gorm.DB {
	once.Do(func() {
		database = initDB()
	})
	return database
}

func initDB() *gorm.DB {

	dsn := os.Getenv("DATABASE_URL")

	// Open connection
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Migrate models
	db.AutoMigrate(&User{})
	db.AutoMigrate(&Room{})
	db.AutoMigrate(&Message{})

	return db
}

func Paginate(page int, size int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		switch {
		case size > 100:
			size = 100
		case size <= 0:
			size = 10
		}

		offset := (page - 1) * size
		return db.Offset(offset).Limit(size)
	}
}
