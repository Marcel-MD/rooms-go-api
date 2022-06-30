package handlers

import "github.com/gin-gonic/gin"

func InitRouter() *gin.Engine {
	e := gin.Default()
	r := e.Group("/api")

	NewUserHandler().Route(r)
	NewRoomHandler().Route(r)

	return e
}
