package handlers

import (
	"net/http"

	"github.com/Marcel-MD/rooms-go-api/middleware"
	"github.com/Marcel-MD/rooms-go-api/services"
	"github.com/Marcel-MD/rooms-go-api/token"
	"github.com/Marcel-MD/rooms-go-api/websockets"
	"github.com/gin-gonic/gin"
)

type IWebSocketHandler interface {
	Route(r *gin.RouterGroup)
}

type WebSocketHandler struct {
	MessageService services.IMessageService
}

func NewWebSocketHandler() IWebSocketHandler {
	return &WebSocketHandler{
		MessageService: services.GetMessageService(),
	}
}

func (h *WebSocketHandler) Route(router *gin.RouterGroup) {
	r := router.Group("/ws").Use(middleware.JwtAuth())
	r.GET("/:room_id", h.Connect)
}

func (h *WebSocketHandler) Connect(c *gin.Context) {
	roomID := c.Param("room_id")

	userID, err := token.ExtractID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	err = h.MessageService.VerifyUserInRoom(roomID, userID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	websockets.ServeWs(c.Writer, c.Request, roomID, userID)
}
