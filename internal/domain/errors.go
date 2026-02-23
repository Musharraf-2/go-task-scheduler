package domain

import "errors"

var ErrRetryable = errors.New("retryable error")

var ErrFatal = errors.New("fatal error")

var ErrMissingID =  errors.New("task is is required")

var ErrMissingKind = errors.New("task kind is required")