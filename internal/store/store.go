package store

import (
	"context"
	"goq/internal/domain"
	"time"
)

type Store interface {
  Enqueue(ctx context.Context, task *domain.Task) error
  Ack(ctx context.Context, taskID string) error
  Fail(ctx context.Context, taskID string, nextRunAt time.Time, parmanent bool) error
  Retry(ctx context.Context, taskID string, nextRunAt time.Time) error
  PromoteSehedules(ctx context.Context, now time.Time) error
  RequeueExpiredLeases(ctx context.Context, until time.Time) error
  Lease(ctx context.Context, queue string, limit int, leaseUntil time.Time) ([]*domain.Task, error)
}