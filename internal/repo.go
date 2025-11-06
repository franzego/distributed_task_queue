package internal

import (
	"context"
	"fmt"

	db "github.com/franzego/distributed_task_queue/db/sqlc"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	q      db.Queries
	dbconn *pgxpool.Pool
}

func NewRepositoryService(dbconn *pgxpool.Pool) *Repository {
	if dbconn == nil {
		return nil
	}
	return &Repository{
		q:      *db.New(dbconn),
		dbconn: dbconn,
	}
}

func (r *Repository) CreateJob(ctx context.Context, arg db.CreateJobParams) (db.Job, error) {
	job, err := r.q.CreateJob(ctx, arg)
	if err != nil {
		return db.Job{}, fmt.Errorf("could not create job in db: %w", err)
	}
	return job, nil
}
func (r *Repository) GetJob(ctx context.Context, id string) (db.Job, error) {
	job, err := r.q.GetJob(ctx, id)
	if err != nil {
		return db.Job{}, fmt.Errorf("could not get job of id %s in db", id)
	}
	return job, nil
}
func (r *Repository) DequeueJob(ctx context.Context) (db.Job, error) {
	job, err := r.q.DequeueJob(ctx)
	if err != nil {
		return db.Job{}, fmt.Errorf("could not dequeue job: %w", err)
	}
	return job, nil
}
func (r *Repository) FailJob(ctx context.Context, arg db.FailJobParams) error {
	return r.q.FailJob(ctx, arg)
}
func (r *Repository) CompletedJob(ctx context.Context, id string) error {
	return r.q.CompleteJob(ctx, id)
}
