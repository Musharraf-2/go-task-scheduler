package app

import (
	"context"
	"testing"
	"time"

	"goq/internal/domain"
)

func TestEngine_EnqueueAndProcess(t *testing.T) {
	cfg := DefaultEngineConfig()
	cfg.Scheduler.Interval = 200 * time.Millisecond
	cfg.Scheduler.LeaseExpiry = 5 * time.Second
	cfg.Worker.Workers = 1
	cfg.Worker.LeaseDuration = 1 * time.Minute

	eng := NewEngine(cfg)
	done := make(chan struct{})
	eng.Worker.Register("Echo", func(ctx context.Context, task *domain.Task) error {
		close(done)
		return nil
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go eng.Run(ctx)

	task := &domain.Task{
		ID:      "test-1",
		Queue:   "default",
		Kind:    "Echo",
		Payload: map[string]interface{}{"x": 1},
		MaxAttempts: 3,
	}
	if err := eng.Queue.Enqueue(ctx, task); err != nil {
		t.Fatal(err)
	}

	select {
	case <-done:
		// handler ran
	case <-time.After(5 * time.Second):
		t.Fatal("handler did not run in time")
	}
	cancel()
	time.Sleep(300 * time.Millisecond)

	// Task should be Done in store
	st := eng.Store
	tt := st.GetTask("test-1")
	if tt == nil || tt.State != domain.StateDone {
		t.Errorf("expected task state Done, got %v", tt)
	}
}