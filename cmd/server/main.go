package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/franzego/distributed_task_queue/auth"
	"github.com/franzego/distributed_task_queue/internal"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("there was a problem loading .env file")
	}

	dburl := "postgresql://franz:Franzego%401@localhost:5433/jobsdb?sslmode=disable"
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	dbConn, err := pgxpool.New(ctx, dburl)
	if err != nil {
		log.Printf("Error connecting to database: %v", err)
		return
	}
	repository := internal.NewRepositoryService(dbConn)
	service := internal.NewService(repository)
	middlewareAuth := auth.NewMiddlewareService(repository)
	handler := internal.NewHandlerService(service)

	r := gin.Default()
	// This is a public path
	r.GET("/health", func(c *gin.Context) {
		// Return JSON response
		c.JSON(http.StatusOK, gin.H{
			"message": "ok",
		})
	})

	// This is a protected path for admin endpoint
	admin := r.Group("/admin")
	admin.Use(middlewareAuth.AdminAuth())
	{
		admin.POST("/api-keys", handler.PostAdminApiKey)
		admin.GET("/api-keys", handler.GetApiKeys)
	}

	// This is a protected path for jobs endpint
	api := r.Group("/")
	api.Use(middlewareAuth.AuthMiddlerWare())
	api.Use(middlewareAuth.RateLimit())
	{
		api.POST("/jobs", handler.PostJob)
		api.GET("/jobs/:id", handler.GetStatus)
	}

	r.Run()
}
