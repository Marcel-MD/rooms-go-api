package handlers

import (
	"net/http"

	"github.com/Marcel-MD/rooms-go-api/dto"
	"github.com/Marcel-MD/rooms-go-api/middleware"
	"github.com/Marcel-MD/rooms-go-api/services"
	"github.com/gin-gonic/gin"
)

type roomHandler struct {
	service services.IRoomService
}

func newRoomHandler() handler {
	return &roomHandler{
		service: services.GetRoomService(),
	}
}

func (h *roomHandler) route(router *gin.RouterGroup) {
	r := router.Group("/rooms")
	r.GET("/", h.findAll)
	r.GET("/:id", h.findOne)

	p := r.Use(middleware.JwtAuth())
	p.POST("/", h.create)
	p.PUT("/:id", h.update)
	p.DELETE("/:id", h.delete)
	p.POST("/:id/users/:email", h.addUser)
	p.DELETE("/:id/users/:email", h.removeUser)
}

func (h *roomHandler) findAll(c *gin.Context) {
	rooms := h.service.FindAll()
	c.JSON(http.StatusOK, rooms)
}

func (h *roomHandler) findOne(c *gin.Context) {
	id := c.Param("id")

	room, err := h.service.FindOne(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "room not found"})
		return
	}

	c.JSON(http.StatusOK, room)
}

func (h *roomHandler) create(c *gin.Context) {

	userID := c.GetString("user_id")

	var dto dto.CreateRoom
	err := c.ShouldBindJSON(&dto)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	room, err := h.service.Create(dto, userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, room)
}

func (h *roomHandler) update(c *gin.Context) {
	id := c.Param("id")

	userID := c.GetString("user_id")

	var dto dto.UpdateRoom
	err := c.ShouldBindJSON(&dto)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	room, err := h.service.Update(id, dto, userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, room)
}

func (h *roomHandler) delete(c *gin.Context) {
	id := c.Param("id")

	userID := c.GetString("user_id")

	err := h.service.Delete(id, userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "room deleted"})
}

func (h *roomHandler) addUser(c *gin.Context) {
	id := c.Param("id")
	email := c.Param("email")

	userID := c.GetString("user_id")

	err := h.service.AddUser(id, email, userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "user added"})
}

func (h *roomHandler) removeUser(c *gin.Context) {
	id := c.Param("id")
	email := c.Param("email")

	userID := c.GetString("user_id")

	err := h.service.RemoveUser(id, email, userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "user removed"})
}
