package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/franzego/distributed_task_queue/internal"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
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
	handler := internal.NewHandlerService(service)

	r := gin.Default()
	r.GET("/health", func(c *gin.Context) {
		// Return JSON response
		c.JSON(http.StatusOK, gin.H{
			"message": "ok",
		})
	})
	r.POST("/jobs", handler.PostJob)
	r.GET("/jobs/status", handler.GetStatus)

	r.Run()
}
