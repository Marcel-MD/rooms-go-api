package handlers

import (
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type handler interface {
	route(r *gin.RouterGroup)
}

var once sync.Once

func InitRouter() {
	once.Do(func() {
		log.Info().Msg("Initializing router")
		e := gin.Default()
		r := e.Group("/api")

		newUserHandler().route(r)
		newRoomHandler().route(r)
		newMessageHandler().route(r)
		newWebSocketHandler().route(r)

		e.Run()
	})
}
