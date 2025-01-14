package middlewares

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/aminul-i-abid/url-shortener/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

var (
	rateLimit = 50
	duration  = 24 * time.Hour
	client    *redis.Client
)

func init() {
	client = redis.NewClient(&redis.Options{
		Addr: "redis:6379",
	})

	// Verify Redis connection
	pong, err := client.Ping(context.Background()).Result()
	if err != nil {
		log.Fatalf("Redis connection failed: %v", err)
	}
	log.Println("Redis connection successful:", pong)
}

func RateLimiterMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		userIP := c.ClientIP()
		key := "rate_limit:" + userIP

		// Check if the user has exceeded the rate limit
		requestCount, err := client.Get(context.Background(), key).Int()
		if err != nil && err != redis.Nil {
			utils.WriteJSON(c.Writer, http.StatusInternalServerError, "Internal server error", nil)
			c.Abort()
			return
		}

		// If the rate limit exceeded, reject the request
		if requestCount >= rateLimit {
			utils.WriteJSON(c.Writer, http.StatusTooManyRequests, "Rate limit exceeded, try again tomorrow", nil)
			c.Abort()
			return
		}

		// Increment request count and set expiration atomically
		_, err = client.Pipelined(ctx, func(pipe redis.Pipeliner) error {
			pipe.Incr(ctx, key)
			pipe.Expire(ctx, key, duration)
			return nil
		})
		if err != nil {
			utils.WriteJSON(c.Writer, http.StatusInternalServerError, "Internal server error", nil)
			c.Abort()
			return
		}

		c.Next()
	}
}
