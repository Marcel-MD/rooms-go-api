package main

import (
	"log"

	"github.com/Marcel-MD/rooms-go-api/handlers"
	"github.com/Marcel-MD/rooms-go-api/websockets"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	websockets.InitHub()
	handlers.InitRouter()
}
