package auth

import (
	"context"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/franzego/distributed_task_queue/internal"
	"github.com/franzego/distributed_task_queue/models"
	"github.com/gin-gonic/gin"
)

type Middleware struct {
	q *internal.Repository
}

func NewMiddlewareService(q *internal.Repository) *Middleware {
	if q == nil {
		return nil
	}
	return &Middleware{
		q: q,
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
		hash := c.GetHeader("X-API-Key")

		if hash == "" {
			c.JSON(http.StatusUnauthorized, models.ErrorResponse{
				Message: "Missing API key",
			})
			c.Abort()
			return
		}
		// hash the apikey
		// hash := HashApiKeys(apiKey)
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
		go a.q.UpdateLastUsed(ctx, hashKey.ID)
		c.Next()

	}
}
