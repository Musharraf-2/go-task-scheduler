package worker

import (
	"context"
	"errors"
	"goq/internal/domain"
	"goq/internal/queue"
	"log"
	"math/rand"
	"sync"
	"time"
)

type Pool struct {
	q *queue.Queue
	handlers map[string]Handler
	workers int
	leaseDuration time.Duration
	queue string
}

type PoolConfig struct {
	Workers int
	LeaseDuration time.Duration
	QueueName string
}

func NewPool(q *queue.Queue, config PoolConfig) *Pool {
  if config.Workers <= 0 {
    config.Workers = 5
  }
  if config.LeaseDuration <= 0 {
    config.LeaseDuration = 2 * time.Minute
  }
  if config.QueueName == "" {
    config.QueueName = "default"
  }
  return &Pool{
    q: q,
    handlers: make(map[string]Handler),
    workers: config.Workers,
    queue:  config.QueueName,
    leaseDuration: config.LeaseDuration,
  }
}

func (p *Pool) Run(ctx context.Context) {
  var wg sync.WaitGroup
  for i := 0 ; i < p.workers; i++ {
    wg.Add(1)
    go func() {
      defer wg.Done()
      p.runWorker(ctx)
    }()
  }

  wg.Wait()
}

func (p *Pool) runWorker(ctx context.Context) {
 for {
  select {
     case <-ctx.Done():
      return
     default:
  }

  leaseUntil := time.Now().Add(p.leaseDuration)
  tasks, err := p.q.Lease(ctx, p.queue, 1, leaseUntil)
  if err != nil {
    log.Printf("worker lease error: %v", err)
  }

  if len(tasks) == 0 {
    time.Sleep(100 * time.Millisecond)
    continue
  }

  for _, task := range tasks {
    p.handleTask(ctx, task)
  }
 }
}

func (p *Pool) handleTask(ctx context.Context, task *domain.Task) {

  h, ok := p.handlers[task.Kind]

  if !ok {
    log.Printf("no handler found for kind %q, failing task kind: %s", task.Kind, task.ID)
    _ = p.q.Fail(ctx, task.ID, time.Time{}, true)
  }

  err := h(ctx, task)
	if err == nil {
		_ = p.q.Ack(ctx, task.ID)
		return
	}
	if errors.Is(err, domain.ErrFatal) {
		_ = p.q.Fail(ctx, task.ID, time.Time{}, true)
		return
	}
	task.Attempts++
	if task.Attempts >= task.MaxAttempts {
		_ = p.q.Fail(ctx, task.ID, time.Time{}, true)
		return
	}
	backoff := backoffWithJitter(task.Attempts)
	nextRun := time.Now().Add(backoff)
	_ = p.q.Retry(ctx, task.ID, nextRun)

}

func backoffWithJitter(attempt int) time.Duration {
	base := time.Duration(attempt*attempt) * time.Second
	jitter := time.Duration(rand.Intn(1000)) * time.Millisecond
	return base + jitter
}

func (p *Pool) Register(kind string, h Handler) {
	p.handlers[kind] = h
}