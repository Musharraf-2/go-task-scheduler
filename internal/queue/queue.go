package queue

import (
	"context"
	"goq/internal/domain"
	"goq/internal/store"
	"time"
)

type Queue struct {
	store store.Store
}

func New(store store.Store) *Queue {
  return &Queue{store: store}
}

func (q *Queue) Enqueue(ctx context.Context, task *domain.Task) error {
  if task.Queue == "" {
    task.Queue = "default"
  }
  if task.MaxAttempts <= 0 {
    task.MaxAttempts = 3
  }
  now := time.Now()
  if task.NextRunAt.IsZero() || task.NextRunAt.Before(now) {
    task.State =  domain.StatePending
  }else {
    task.State=  domain.StateScheduled
  }
  task.CreatedAt = now
  task.UpdatedAt = now
  return q.store.Enqueue(ctx, task)
}

func (q *Queue) Lease(ctx context.Context, queue string, limit int, leaseUntil time.Time) ([]*domain.Task, error) {
  return q.store.Lease(ctx, queue, limit, leaseUntil)
}

func (q *Queue) Ack(ctx context.Context, taskID string) error {
  return q.store.Ack(ctx, taskID)
}

func (q *Queue) Fail(ctx context.Context, taskID string, nextRunAt time.Time, permanent bool) error {
  return q.store.Fail(ctx, taskID, nextRunAt, permanent)
}

func (q *Queue) Retry(ctx context.Context, taskID string, nextRunAt time.Time) error {
  return q.store.Retry(ctx, taskID, nextRunAt)
}