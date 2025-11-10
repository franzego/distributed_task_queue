package auth

import (
	"context"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/franzego/distributed_task_queue/internal"
	"github.com/franzego/distributed_task_queue/internal/ratelimit"
	"github.com/franzego/distributed_task_queue/models"
	"github.com/gin-gonic/gin"
)

type Middleware struct {
	q           *internal.Repository
	rateLimiter *ratelimit.RateLimiter
}

func NewMiddlewareService(q *internal.Repository) *Middleware {
	if q == nil {
		return nil
	}
	// Initialize rate limiter with capacity of 100 requests and refill rate of 10 tokens per second
	rateLimiter := ratelimit.NewRateLimiterService(100, 10)
	return &Middleware{
		q:           q,
		rateLimiter: rateLimiter,
	}
}

// RateLimit middleware applies rate limiting based on api key
func (a *Middleware) RateLimit() gin.HandlerFunc {
	return func(c *gin.Context) {

		keyID, exists := c.Get("api_key_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, models.ErrorResponse{
				Message: "Api key does not exist",
			})
			c.Abort()
			return
		}
		if !a.rateLimiter.Allow(keyID.(string)) {
			c.JSON(http.StatusTooManyRequests, models.ErrorResponse{
				Message: "Rate limit exceeded. Please try again later.",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}

// Admin Authentication
func (a *Middleware) AdminAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, models.ErrorResponse{
				Message: "Missing Authorization Header",
			})
			c.Abort()
			return
		}
		bearerToken := strings.Split(authHeader, " ")
		if len(bearerToken) != 2 || bearerToken[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, models.ErrorResponse{
				Message: "Invalid Authorization format",
			})
			c.Abort()
			return
		}
		adminToken := os.Getenv("ADMIN_TOKEN")
		token := bearerToken[1]
		if token != adminToken {
			c.JSON(http.StatusUnauthorized, models.ErrorResponse{
				Message: "Invalid Access Token",
			})
			c.Abort()
			return
		}
		c.Set("is_admin", true)
		c.Next()
	}
}

func (a *Middleware) AuthMiddlerWare() gin.HandlerFunc {
	return func(c *gin.Context) {
		apiKey := c.GetHeader("X-API-Key")

		if apiKey == "" {
			c.JSON(http.StatusUnauthorized, models.ErrorResponse{
				Message: "Missing API key",
			})
			c.Abort()
			return
		}
		// hash the apikey
		hash := HashApiKeys(apiKey)
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		hashKey, err := a.q.GetAPIKeyByHash(ctx, hash)
		if err != nil {
			c.JSON(http.StatusUnauthorized, models.ErrorResponse{
				Message: "Invalid Api key",
				Error:   err.Error(),
			})
			c.Abort()
			return
		}
		if hashKey.ExpiresAt.Valid && hashKey.ExpiresAt.Time.Before(time.Now()) {
			c.JSON(http.StatusUnauthorized, models.ErrorResponse{
				Message: "Api Key Has Expired",
			})
			c.Abort()
			return
		}
		c.Set("api_key_id", hashKey.ID)
		c.Set("api_key_name", hashKey.Name)
		go func() {
			if err := a.q.UpdateLastUsed(ctx, hashKey.ID); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": "failed to update time",
				})
			}
		}()

		c.Next()

	}
}
