package models

import (
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var database *gorm.DB

func GetDB() *gorm.DB {
	return database
}

func InitDB() {

	dsn := "postgres://postgres:password@localhost:5432/rooms"

	// Get .env connection url if exists
	err := godotenv.Load(".env")
	if err == nil {
		dsn = os.Getenv("DATABASE_URL")
	}

	// Open connection
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		panic("failed to connect database")
	}

	// Migrate models
	db.AutoMigrate(&User{})

	database = db
}
