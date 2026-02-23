package queue

import (
	"goq/internal/domain"
)



func Validate(task *domain.Task) error {
	if task.ID == "" {
		return domain.ErrMissingID 
	}
	if task.Kind == "" {
		return domain.ErrMissingKind
	}
	return nil
}