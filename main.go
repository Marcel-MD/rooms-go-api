package main

import (
	"log"

	"github.com/Marcel-MD/rooms-go-api/handlers"
	"github.com/Marcel-MD/rooms-go-api/models"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	models.InitDB()
	r := handlers.InitRouter()
	r.Run()
}
