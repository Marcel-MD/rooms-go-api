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
	r.POST("/register-otp", h.registerOtp)
	r.POST("/login", h.login)
	r.POST("/login-otp", h.loginOtp)
	r.POST("/send-otp", h.sendOtp)
	r.GET("/", h.findAll)
	r.GET("/all", h.findAll)
	r.GET("/email/:email", h.searchByEmail)
	r.GET("/:id", h.findOne)

	p := r.Use(middleware.JwtAuth())
	p.GET("/current", h.current)
	p.PUT("/update", h.update)
	p.PUT("/update-otp", h.updateOtp)

	p.POST("/:id/roles/:role", h.addRole)
	p.DELETE("/:id/roles/:role", h.removeRole)
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

func (h *userHandler) registerOtp(c *gin.Context) {

	var dto dto.RegisterOtpUser
	err := c.ShouldBindJSON(&dto)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.service.RegisterOtp(dto)
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
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

func (h *userHandler) loginOtp(c *gin.Context) {

	var dto dto.LoginOtpUser
	err := c.ShouldBindJSON(&dto)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := h.service.LoginOtp(dto)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

func (h *userHandler) sendOtp(c *gin.Context) {

	var dto dto.Email
	err := c.ShouldBindJSON(&dto)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.service.SendOtp(dto.Email)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "otp sent"})
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

func (j *userHandler) searchByEmail(c *gin.Context) {
	email := c.Param("email")

	users := j.service.SearchByEmail(email)
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

func (h *userHandler) updateOtp(c *gin.Context) {
	userID := c.GetString("user_id")

	var dto dto.UpdateOtpUser
	err := c.ShouldBindJSON(&dto)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.service.UpdateOtp(dto, userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h *userHandler) addRole(c *gin.Context) {
	userID := c.GetString("user_id")
	id := c.Param("id")
	role := c.Param("role")

	user, err := h.service.AddRole(id, role, userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h *userHandler) removeRole(c *gin.Context) {
	userID := c.GetString("user_id")
	id := c.Param("id")
	role := c.Param("role")

	user, err := h.service.RemoveRole(id, role, userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}
