package middleware

import (
	"fmt"
	"net/http"
	"time"

	"taskmaster/internal/ports"

	"github.com/gin-gonic/gin"
)

func RateLimit(cache ports.CacheRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		key := fmt.Sprintf("ratelimit:%s:%s", c.ClientIP(), c.FullPath())
		count, _ := cache.Increment(c.Request.Context(), key)
		if count == 1 {
			_ = cache.Expire(c.Request.Context(), key, 60)
		}
		if count > 100 {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"error": "rate limit exceeded"})
			return
		}
		c.Next()
	}
}
