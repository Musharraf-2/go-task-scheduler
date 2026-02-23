package queue

import (
	"errors"
	"goq/internal/domain"
)

var ErrMissingD =  errors.New("task is is required")
var ErrMissingKind = errors.New("task kind is required")

func Validate(task *domain.Task) error {
	if task.ID == "" {
		return ErrMissingD 
	}
	if task.Kind == "" {
		return ErrMissingKind
	}
	return nil
}