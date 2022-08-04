package handlers

import (
	"net/http"
	"strconv"

	"github.com/Marcel-MD/rooms-go-api/dto"
	"github.com/Marcel-MD/rooms-go-api/middleware"
	"github.com/Marcel-MD/rooms-go-api/services"
	"github.com/gin-gonic/gin"
)

type messageHandler struct {
	service services.IMessageService
}

func routeMessageHandler(router *gin.RouterGroup) {
	h := &messageHandler{
		service: services.GetMessageService(),
	}

	r := router.Group("/messages").Use(middleware.JwtAuth())
	r.GET("/:room_id", h.find)
	r.POST("/:room_id", h.create)
	r.PUT("/:id", h.update)
	r.DELETE("/:id", h.delete)
}

func (h *messageHandler) find(c *gin.Context) {
	roomID := c.Param("room_id")
	var err error
	params := dto.MessageQueryParams{}

	params.Page, err = strconv.Atoi(c.Query("page"))
	if err != nil {
		params.Page = 1
	}

	params.Size, err = strconv.Atoi(c.Query("size"))
	if err != nil {
		params.Size = 20
	}

	userID := c.GetString("user_id")

	messages, err := h.service.FindByRoomID(roomID, userID, params)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, messages)
}

func (h *messageHandler) create(c *gin.Context) {
	roomID := c.Param("room_id")

	var dto dto.CreateMessage
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.GetString("user_id")

	message, err := h.service.Create(roomID, userID, dto)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, message)
}

func (h *messageHandler) update(c *gin.Context) {
	id := c.Param("id")

	var dto dto.UpdateMessage
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.GetString("user_id")

	message, err := h.service.Update(id, userID, dto)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, message)
}

func (h *messageHandler) delete(c *gin.Context) {
	id := c.Param("id")

	userID := c.GetString("user_id")

	err := h.service.Delete(id, userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "message deleted"})
}
