package main

import (
	"github.com/Marcel-MD/rooms-go-api/handlers"
	"github.com/Marcel-MD/rooms-go-api/logger"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
)

func main() {

	err := godotenv.Load(".env")
	if err != nil {
		log.Err(err).Msg("Error loading .env file")
	}

	logger.Config()
	handlers.InitRouter()
}
