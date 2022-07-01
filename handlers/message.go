package handlers

import (
	"net/http"
	"strconv"

	"github.com/Marcel-MD/rooms-go-api/dto"
	"github.com/Marcel-MD/rooms-go-api/middleware"
	"github.com/Marcel-MD/rooms-go-api/services"
	"github.com/Marcel-MD/rooms-go-api/token"
	"github.com/gin-gonic/gin"
)

type IMessageHandler interface {
	Route(r *gin.RouterGroup)
}

type MessageHandler struct {
	Service services.IMessageService
}

func NewMessageHandler() IMessageHandler {
	return &MessageHandler{
		Service: services.NewMessageService(),
	}
}

func (h *MessageHandler) Route(router *gin.RouterGroup) {

	r := router.Group("/messages").Use(middleware.JwtAuth())

	r.GET("/:room_id", h.Find)
	r.POST("/:room_id", h.Create)
	r.PUT("/:id", h.Update)
	r.DELETE("/:id", h.Delete)
}

func (h *MessageHandler) Find(c *gin.Context) {
	roomID := c.Param("room_id")
	var err error
	params := dto.MessageQueryParams{}
	params.Page, err = strconv.Atoi(c.Query("page"))
	params.Size, err = strconv.Atoi(c.Query("size"))

	if err != nil {
		params.Page = 1
		params.Size = 20
	}

	userID, err := token.ExtractID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	messages, err := h.Service.FindByRoomID(roomID, userID, params)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, messages)
}

func (h *MessageHandler) Create(c *gin.Context) {
	roomID := c.Param("room_id")

	var dto dto.CreateMessage
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, err := token.ExtractID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	message, err := h.Service.Create(dto, roomID, userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, message)
}

func (h *MessageHandler) Update(c *gin.Context) {
	id := c.Param("id")

	var dto dto.UpdateMessage
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, err := token.ExtractID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	message, err := h.Service.Update(id, dto, userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, message)
}

func (h *MessageHandler) Delete(c *gin.Context) {
	id := c.Param("id")

	userID, err := token.ExtractID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	err = h.Service.Delete(id, userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "message deleted"})
}
