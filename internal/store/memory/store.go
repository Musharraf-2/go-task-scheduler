package memory

import (
	"context"
	"goq/internal/domain"
	"goq/internal/store"
	"sync"
	"time"
)

type Store struct {
	mu    sync.Mutex
	tasks map[string]*domain.Task
}

func (s *Store) Ack(ctx context.Context, taskID string) error {
	s.mu.Lock()
  defer s.mu.Unlock()

  if t, ok := s.tasks[taskID]; ok {
    t.UpdatedAt = time.Now()
    t.State = domain.StateDone
  }
  return nil
}

func (s *Store) Enqueue(ctx context.Context, task *domain.Task) error {
	s.mu.Lock()
  defer s.mu.Unlock()

  copiedTask := copyTask(task)
  now := time.Now()
  copiedTask.CreatedAt = now
  copiedTask.UpdatedAt = now
  s.tasks[copiedTask.ID] = copiedTask
  return nil
}

func (s *Store) Fail(ctx context.Context, taskID string, nextRunAt time.Time, permanent bool) error {
	s.mu.Lock()
  defer s.mu.Unlock()

  if t, ok := s.tasks[taskID]; ok {
    if permanent {
      t.State =  domain.StateDead
    } else {
      t.State = domain.StateFailed
      t.NextRunAt = nextRunAt
    }
    t.UpdatedAt = time.Now()
  }
  
  return nil
}

func (s *Store) Lease(ctx context.Context, queue string, limit int, leaseUntil time.Time) ([]*domain.Task, error) {
  s.mu.Lock()
  defer s.mu.Unlock()

  now := time.Now()
  var out []*domain.Task
  for _, task := range s.tasks {
    if len(out) >= limit {
      break
    }

    if task.Queue != queue || task.State != domain.StatePending {
      continue
    }

    if !task.RunAt.IsZero() &&  task.RunAt.After(now) {
      continue
    }

    t := copyTask(task)
    t.State = domain.StateRunning
    t.UpdatedAt = leaseUntil
    s.tasks[t.ID] = t
    out = append(out,t)
  }
  return out, nil
}

func (s *Store) PromoteSchedules(ctx context.Context, now time.Time) error {
	s.mu.Lock()
  defer s.mu.Unlock()

  for _, task := range s.tasks {
    if task.State == domain.StateScheduled && !task.RunAt.After(now) {
      task.State =  domain.StatePending
      task.UpdatedAt = now
    }
  }
  return nil
}

func (s *Store) RequeueExpiredLeases(ctx context.Context, until time.Time) error {
  s.mu.Lock()
  defer s.mu.Unlock()

  for _, task := range s.tasks {
    if task.State == domain.StateRunning && !task.UpdatedAt.Before(until) {
      task.State = domain.StatePending
      task.UpdatedAt = time.Now()
    }
  }
  return nil
}

func (s *Store) Retry(ctx context.Context, taskID string, nextRunAt time.Time) error {
	s.mu.Lock()
  defer s.mu.Unlock()

  if t, ok := s.tasks[taskID]; ok {
    t.State = domain.StatePending
    t.NextRunAt = nextRunAt
    t.UpdatedAt = time.Now()
  }
  return nil
}

func copyTask(task *domain.Task) * domain.Task {
	newTask := *task
	payload := make(map[string]interface{}, len(task.Payload))
	for k,v := range task.Payload {
		payload[k] = v
	}
	newTask.Payload =  payload
	return &newTask
}

func New() *Store {
	return &Store{
		tasks: make(map[string]*domain.Task),
	}
}

var _ store.Store = (*Store)(nil)
