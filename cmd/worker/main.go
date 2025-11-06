package main

import (
	"context"
	"log"
	"time"

	"github.com/franzego/distributed_task_queue/internal"
	"github.com/franzego/distributed_task_queue/worker"
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
	worker := worker.NewWorkerService(repository)
	log.Println("Starting worker...")
	if err := worker.WorkerFunction(); err != nil {
		log.Fatal(err)
	}
}
