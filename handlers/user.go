package handlers

import (
	"net/http"

	"github.com/Marcel-MD/rooms-go-api/dto"
	"github.com/Marcel-MD/rooms-go-api/middleware"
	"github.com/Marcel-MD/rooms-go-api/services"
	"github.com/Marcel-MD/rooms-go-api/token"
	"github.com/gin-gonic/gin"
)

type IUserHandler interface {
	Route(r *gin.RouterGroup)
}

type UserHandler struct {
	Service services.IUserService
}

func NewUserHandler() IUserHandler {
	return &UserHandler{
		Service: services.GetUserService(),
	}
}

func (h *UserHandler) Route(router *gin.RouterGroup) {

	r := router.Group("/users")
	r.POST("/register", h.Register)
	r.POST("/login", h.Login)
	r.GET("/", h.FindAll)
	r.GET("/:id", h.FindOne)

	p := r.Use(middleware.JwtAuth())
	p.GET("/current", h.Current)
}

func (h *UserHandler) Register(c *gin.Context) {
	var dto dto.RegisterUser

	err := c.ShouldBindJSON(&dto)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.Service.Register(dto)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h *UserHandler) Login(c *gin.Context) {
	var dto dto.LoginUser

	err := c.ShouldBindJSON(&dto)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := h.Service.Login(dto)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "username or password is incorrect."})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})

}

func (h *UserHandler) Current(c *gin.Context) {

	id, err := token.ExtractID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	user, err := h.Service.FindOne(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "record not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h *UserHandler) FindAll(c *gin.Context) {
	users := h.Service.FindAll()
	c.JSON(http.StatusOK, users)
}

func (h *UserHandler) FindOne(c *gin.Context) {
	user, err := h.Service.FindOne(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "record not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}
