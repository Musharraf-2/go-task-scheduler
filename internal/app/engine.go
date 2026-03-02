// Package app wires the queue, worker, and scheduler for a headless engine (Phase 1).
// No HTTP here; that comes in Phase 2.
package app

import (
	"context"
	"sync"

	"goq/internal/queue"
	"goq/internal/scheduler"
	"goq/internal/store/memory"
	"goq/internal/worker"
)

// Engine holds the core components and runs them until shutdown.
type Engine struct {
	Store     *memory.Store
	Queue     *queue.Queue
	Worker    *worker.Pool
	Scheduler *scheduler.Scheduler
}

// NewEngine creates store, queue, worker pool, and scheduler. Handlers must be registered on Engine.Worker before Run.
func NewEngine(cfg EngineConfig) *Engine {
	st := memory.New()
	q := queue.New(st)
	wp := worker.NewPool(q, cfg.Worker)
	sched := scheduler.New(st, cfg.Scheduler)
	return &Engine{
		Store:     st,
		Queue:     q,
		Worker:    wp,
		Scheduler: sched,
	}
}

// EngineConfig holds config for worker and scheduler.
type EngineConfig struct {
	Worker   worker.PoolConfig
	Scheduler scheduler.Config
}

// Run starts the worker pool and scheduler in the background. Returns when ctx is cancelled; drains gracefully.
func (e *Engine) Run(ctx context.Context) {
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		e.Worker.Run(ctx)
	}()
	go func() {
		defer wg.Done()
		e.Scheduler.Run(ctx)
	}()
	wg.Wait()
}