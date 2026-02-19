package domain

import "errors"

var ErrRetryable = errors.New("retryable error")

var ErrFatal = errors.New("fatal error")

