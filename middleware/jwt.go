package middleware

import (
	"net/http"

	"github.com/Marcel-MD/rooms-go-api/token"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func JwtAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := token.ExtractID(c)
		if err != nil {
			log.Err(err).Msg("Invalid token")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			c.Abort()
			return
		}
		c.Set("user_id", id)
		c.Next()
	}
}
