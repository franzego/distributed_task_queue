package internal

import (
	"context"

	db "github.com/franzego/distributed_task_queue/db/sqlc"
)

type Queue interface {
	Enqueue(ctx context.Context, job db.Job) error
	Dequeue(ctx context.Context) (*db.Job, error)
	GetJob(ctx context.Context, id string) (db.Job, error)
	CreateAPIKey(ctx context.Context, arg db.CreateAPIKeyParams) (db.ApiKey, error)
	ListAPIKeys(ctx context.Context) ([]db.ApiKey, error)
}

type Service struct {
	// Q Queue
	r *Repository
}

func NewService(r *Repository) *Service {
	if r == nil {
		return nil
	}
	return &Service{
		r: r,
	}
}

// The enqueue function is the one that actually creates a job in the queue(db).
func (s *Service) Enqueue(ctx context.Context, job db.Job) error {
	arg := db.CreateJobParams{
		ID:          job.ID,
		Type:        job.Type,
		Payload:     job.Payload,
		Status:      job.Status,
		MaxAttempts: job.MaxAttempts,
	}
	_, err := s.r.CreateJob(ctx, arg)
	if err != nil {
		return err
	}
	return nil
}
func (s *Service) Dequeue(ctx context.Context) (*db.Job, error) {
	job, err := s.r.DequeueJob(ctx)
	if err != nil {
		return nil, err
	}
	return &job, err
}
func (s *Service) GetJob(ctx context.Context, id string) (db.Job, error) {
	job, err := s.r.GetJob(ctx, id)
	if err != nil {
		return db.Job{}, err
	}
	return job, nil
}
func (s *Service) CreateAPIKey(ctx context.Context, arg db.CreateAPIKeyParams) (db.ApiKey, error) {
	key, err := s.r.CreateAPIKey(ctx, arg)
	if err != nil {
		return db.ApiKey{}, err
	}
	return key, nil
}
func (s *Service) ListAPIKeys(ctx context.Context) ([]db.ApiKey, error) {
	keys, err := s.r.ListAPIKeys(ctx)
	if err != nil {
		return nil, err
	}
	return keys, nil
}
