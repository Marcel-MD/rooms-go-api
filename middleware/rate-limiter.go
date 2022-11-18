package middleware

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/Marcel-MD/rooms-go-api/rdb"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v9"
)

func RateLimiter() gin.HandlerFunc {

	limitStr := os.Getenv("RATE_LIMIT")
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		limit = 50
	}

	windowStr := os.Getenv("RATE_WINDOW")
	window, err := time.ParseDuration(windowStr)
	if err != nil {
		window = 1 * time.Second
	}

	client, ctx := rdb.GetRDB()

	return func(c *gin.Context) {
		now := time.Now().UnixNano()
		ipAddr := c.ClientIP()

		client.ZRemRangeByScore(ctx, ipAddr, "0", fmt.Sprint(now-(window.Nanoseconds()))).Result()

		reqs, err := client.ZRange(ctx, ipAddr, 0, -1).Result()
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		if len(reqs) >= limit {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error": "too many request",
			})
			return
		}

		c.Next()

		client.ZAddNX(ctx, ipAddr, redis.Z{Score: float64(now), Member: float64(now)})
		client.Expire(ctx, ipAddr, window)
	}
}
