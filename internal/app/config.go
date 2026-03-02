package app

import (
	"time"

	"goq/internal/scheduler"
	"goq/internal/worker"
)

// DefaultEngineConfig returns a config suitable for dev/tests.
func DefaultEngineConfig() EngineConfig {
	return EngineConfig{
		Worker: worker.PoolConfig{
			Workers:   2,
			LeaseDuration: 1 * time.Minute,
			QueueName: "default",
		},
		Scheduler: scheduler.Config{
			Interval:    2 * time.Second,
			LeaseExpiry: 1 * time.Minute,
		},
	}
}