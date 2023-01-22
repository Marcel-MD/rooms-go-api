package models

import (
	"os"
	"sync"

	"github.com/rs/zerolog/log"
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
	log.Info().Msg("Initializing database")

	dsn := os.Getenv("DATABASE_URL")

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to database")
	}

	db.AutoMigrate(&User{})
	db.AutoMigrate(&Room{})
	db.AutoMigrate(&Message{})
	db.AutoMigrate(&AllowedUser{})

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
