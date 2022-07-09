package handlers

import "github.com/gin-gonic/gin"

func InitRouter() {
	e := gin.Default()
	r := e.Group("/api")

	NewUserHandler().Route(r)
	NewRoomHandler().Route(r)
	NewMessageHandler().Route(r)
	NewWebSocketHandler().Route(r)

	e.Run()
}
