package handlers

import (
	"os"
	"sync"

	"github.com/Marcel-MD/rooms-go-api/middleware"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

var once sync.Once

func InitRouter() {
	once.Do(func() {
		log.Info().Msg("Initializing router")
		e := gin.Default()
		e.Use(middleware.CORS())

		env := os.Getenv("ENVIRONMENT")
		if env == "dev" {
			pprof.Register(e)
		}

		r := e.Group("/api")

		routeUserHandler(r)
		routeRoomHandler(r)
		routeMessageHandler(r)
		routeWebSocketHandler(r)

		port := os.Getenv("PORT")
		if port == "" {
			port = "8080"
		}

		e.Run(":" + port)
	})
}
