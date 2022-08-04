package handlers

import (
	"net/http"

	"github.com/Marcel-MD/rooms-go-api/dto"
	"github.com/Marcel-MD/rooms-go-api/middleware"
	"github.com/Marcel-MD/rooms-go-api/services"
	"github.com/gin-gonic/gin"
)

type userHandler struct {
	service services.IUserService
}

func routeUserHandler(router *gin.RouterGroup) {
	h := &userHandler{
		service: services.GetUserService(),
	}

	r := router.Group("/users")
	r.POST("/register", h.register)
	r.POST("/login", h.login)
	r.GET("/", h.findAll)
	r.GET("/:id", h.findOne)

	p := r.Use(middleware.JwtAuth())
	p.GET("/current", h.current)
	p.PUT("/", h.update)
}

func (h *userHandler) register(c *gin.Context) {

	var dto dto.RegisterUser
	err := c.ShouldBindJSON(&dto)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.service.Register(dto)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h *userHandler) login(c *gin.Context) {

	var dto dto.LoginUser
	err := c.ShouldBindJSON(&dto)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := h.service.Login(dto)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "username or password is incorrect."})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

func (h *userHandler) current(c *gin.Context) {
	id := c.GetString("user_id")

	user, err := h.service.FindOne(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "record not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h *userHandler) findAll(c *gin.Context) {
	users := h.service.FindAll()
	c.JSON(http.StatusOK, users)
}

func (h *userHandler) findOne(c *gin.Context) {
	id := c.Param("id")

	user, err := h.service.FindOne(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "record not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h *userHandler) update(c *gin.Context) {
	userID := c.GetString("user_id")

	var dto dto.UpdateUser
	err := c.ShouldBindJSON(&dto)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.service.Update(dto, userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}
