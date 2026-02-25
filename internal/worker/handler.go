package worker

import (
	"context"
	"goq/internal/domain"
)

type Handler func(ctx context.Context, task *domain.Task) error