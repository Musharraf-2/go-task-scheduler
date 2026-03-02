// Package scheduler promotes scheduled tasks and recovers expired leases.
// It runs a loop and does not depend on HTTP or any framework.
package scheduler

import (
	"context"
	"log"
	"time"

	"goq/internal/store"
)

// Scheduler runs PromoteScheduled and RequeueExpiredLeases periodically.
type Scheduler struct {
	store       store.Store
	interval    time.Duration
	leaseExpiry time.Duration
	clock       Clock
}

// Config for the scheduler.
type Config struct {
	Interval    time.Duration
	LeaseExpiry time.Duration
	Clock       Clock
}

// New builds a scheduler. If Config.Clock is nil, RealClock is used.
func New(s store.Store, cfg Config) *Scheduler {
	if cfg.Interval <= 0 {
		cfg.Interval = 5 * time.Second
	}
	if cfg.LeaseExpiry <= 0 {
		cfg.LeaseExpiry = 2 * time.Minute
	}
	if cfg.Clock == nil {
		cfg.Clock = RealClock{}
	}
	return &Scheduler{
		store:       s,
		interval:    cfg.Interval,
		leaseExpiry: cfg.LeaseExpiry,
		clock:       cfg.Clock,
	}
}

// Run blocks until ctx is cancelled. Each tick: promote scheduled, requeue expired leases.
func (s *Scheduler) Run(ctx context.Context) {
	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			now := s.clock.Now()
			if err := s.store.PromoteSchedules(ctx, now); err != nil {
				log.Printf("scheduler PromoteScheduled: %v", err)
			}
			until := now.Add(-s.leaseExpiry)
			if err := s.store.RequeueExpiredLeases(ctx, until); err != nil {
				log.Printf("scheduler RequeueExpiredLeases: %v", err)
			}
		}
	}
}