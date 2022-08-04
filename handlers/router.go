package handlers

import (
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

var once sync.Once

func InitRouter() {
	once.Do(func() {
		log.Info().Msg("Initializing router")
		e := gin.Default()
		r := e.Group("/api")

		routeUserHandler(r)
		routeRoomHandler(r)
		routeMessageHandler(r)
		routeWebSocketHandler(r)

		e.Run()
	})
}
