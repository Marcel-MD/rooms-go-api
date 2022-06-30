package handlers

import (
	"net/http"

	"github.com/Marcel-MD/rooms-go-api/dto"
	"github.com/Marcel-MD/rooms-go-api/middleware"
	"github.com/Marcel-MD/rooms-go-api/services"
	"github.com/Marcel-MD/rooms-go-api/token"
	"github.com/gin-gonic/gin"
)

type IRoomHandler interface {
	Route(r *gin.RouterGroup)
}

type RoomHandler struct {
	Service services.IRoomService
}

func NewRoomHandler() IRoomHandler {
	return &RoomHandler{
		Service: services.NewRoomService(),
	}
}

func (h *RoomHandler) Route(router *gin.RouterGroup) {
	r := router.Group("/rooms")
	r.GET("/", h.FindAll)
	r.GET("/:id", h.FindOne)

	p := r.Use(middleware.JwtAuth())
	p.POST("/", h.Create)
	p.PUT("/:id", h.Update)
	p.DELETE("/:id", h.Delete)
	p.POST("/:id/users/:email", h.AddUser)
	p.DELETE("/:id/users/:email", h.RemoveUser)
}

func (h *RoomHandler) FindAll(c *gin.Context) {
	rooms := h.Service.FindAll()
	c.JSON(http.StatusOK, rooms)
}

func (h *RoomHandler) FindOne(c *gin.Context) {
	id := c.Param("id")

	room, err := h.Service.FindOne(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "room not found"})
		return
	}

	c.JSON(http.StatusOK, room)
}

func (h *RoomHandler) Create(c *gin.Context) {

	userID, err := token.ExtractID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var dto dto.CreateRoom
	err = c.ShouldBindJSON(&dto)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	room, err := h.Service.Create(dto, userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, room)
}

func (h *RoomHandler) Update(c *gin.Context) {
	id := c.Param("id")

	userID, err := token.ExtractID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var dto dto.UpdateRoom
	err = c.ShouldBindJSON(&dto)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	room, err := h.Service.Update(id, dto, userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, room)
}

func (h *RoomHandler) Delete(c *gin.Context) {
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

	c.JSON(http.StatusOK, gin.H{"message": "room deleted"})
}

func (h *RoomHandler) AddUser(c *gin.Context) {
	id := c.Param("id")
	email := c.Param("email")

	userID, err := token.ExtractID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	err = h.Service.AddUser(id, email, userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "user added"})
}

func (h *RoomHandler) RemoveUser(c *gin.Context) {
	id := c.Param("id")
	email := c.Param("email")

	userID, err := token.ExtractID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	err = h.Service.RemoveUser(id, email, userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "user removed"})
}
