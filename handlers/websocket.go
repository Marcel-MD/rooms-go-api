package handlers

import (
	"net/http"

	"github.com/Marcel-MD/rooms-go-api/middleware"
	"github.com/Marcel-MD/rooms-go-api/websockets"
	"github.com/gin-gonic/gin"
)

type webSocketHandler struct {
	server websockets.IServer
}

func routeWebSocketHandler(router *gin.RouterGroup) {
	h := &webSocketHandler{
		server: websockets.GetServer(),
	}

	r := router.Group("/ws").Use(middleware.JwtAuth())
	r.GET("/", h.connect)
}

func (h *webSocketHandler) connect(c *gin.Context) {
	userID := c.GetString("user_id")

	err := h.server.ServeWS(c.Writer, c.Request, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
}
