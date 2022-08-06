package handlers

import (
	"net/http"

	"github.com/Marcel-MD/rooms-go-api/middleware"
	"github.com/Marcel-MD/rooms-go-api/services"
	"github.com/Marcel-MD/rooms-go-api/websockets"
	"github.com/gin-gonic/gin"
)

type webSocketHandler struct {
	service services.IMessageService
	server  websockets.IServer
}

func routeWebSocketHandler(router *gin.RouterGroup) {
	h := &webSocketHandler{
		service: services.GetMessageService(),
		server:  websockets.GetServer(),
	}

	r := router.Group("/ws").Use(middleware.JwtAuth())
	r.GET("/:room_id", h.connect)
}

func (h *webSocketHandler) connect(c *gin.Context) {
	roomID := c.Param("room_id")
	userID := c.GetString("user_id")

	err := h.service.VerifyUserInRoom(roomID, userID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	h.server.ServeWS(c.Writer, c.Request, roomID, userID)
}
