package domain

import "time"

const (
	StatePending = "pending"
	StateScheduled = "scheduled"
	StateRunning = "running"
	StateDone = "done"
	StateFailed = "failed"
	StateDead = "dead"
)

type Task struct {
  ID string
  Queue string
  Kind string
  Payload map[string]interface{}
  RunAt time.Time
  NextRunAt time.Time
  Attempts int
  MaxAttempts int
  CreatedAt time.Time
  UpdatedAt time.Time
}