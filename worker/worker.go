package worker

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	db "github.com/franzego/distributed_task_queue/db/sqlc"
	"github.com/franzego/distributed_task_queue/internal"
	"github.com/jackc/pgx/v5/pgtype"
)

type Worker struct {
	r *internal.Repository
}

func NewWorkerService(r *internal.Repository) *Worker {
	if r == nil {
		return nil
	}
	return &Worker{
		r: r,
	}
}

func (w *Worker) WorkerFunction() error {
	log.Print("Worker has started")
	for {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		// we dequeue
		job, err := w.r.DequeueJob(ctx)
		// cancel the context for this iteration immediately after dequeue returns
		cancel()
		if errors.Is(err, sql.ErrNoRows) {
			log.Println("No Available Jobs. Waiting...")
			time.Sleep(20 * time.Second)
			continue
		}
		if err != nil {
			log.Printf("Error dequeuing: %v", err)
			time.Sleep(2 * time.Second)
			continue
		}
		log.Printf("Processing job %s of type %s", job.ID, job.Type)
		err = w.ProcessJobs(job)
		if err != nil {
			log.Printf("Job %s has failed: %v", job.ID, err)
			w.JobFailed(job, err)
		} else {
			log.Printf("Job %s has completed successfully", job.ID)
			w.CompletedJob(job)
		}

	}

}
func (w *Worker) ProcessJobs(job db.Job) error {
	switch job.Type {
	case "send_email":
		// handler for email
		return nil
	case "logs":
		// handler for logs
		return fmt.Errorf("invalid job type: %s", job.Type)
	default:
		//
		return fmt.Errorf("invalid job type: %s", job.Type)
	}
}

func (w *Worker) JobFailed(job db.Job, jobErr error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	arg := db.FailJobParams{
		ID:       job.ID,
		Status:   "failed",
		Attempts: job.Attempts + 1,
		ErrorMessage: pgtype.Text{
			String: jobErr.Error(),
		},
		ScheduledAt: pgtype.Timestamptz{Time: time.Now()},
	}
	err := w.r.FailJob(ctx, arg)
	if err != nil {
		log.Printf("Failed to mark %s as failed. %v", job.ID, err)
	}
}
func (w *Worker) CompletedJob(job db.Job) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := w.r.CompletedJob(ctx, job.ID); err != nil {
		return fmt.Errorf("error marking job %s as completed %v", job.ID, err)
	}
	return nil
}
