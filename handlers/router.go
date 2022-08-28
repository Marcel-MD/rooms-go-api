package handlers

import (
	"os"
	"sync"

	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

var once sync.Once

func InitRouter() {
	once.Do(func() {
		log.Info().Msg("Initializing router")
		e := gin.Default()

		env := os.Getenv("ENVIRONMENT")
		if env == "dev" {
			pprof.Register(e)
		}

		r := e.Group("/api")

		routeUserHandler(r)
		routeRoomHandler(r)
		routeMessageHandler(r)
		routeWebSocketHandler(r)

		e.Run()
	})
}
