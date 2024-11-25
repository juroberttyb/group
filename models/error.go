package models

import "errors"

var (
	ErrTimeout      = errors.New("timeout")
	ErrReachedLimit = errors.New("reached inflight limit")
	ErrNilContext   = errors.New("context should not be empty")
	ErrNilInput     = errors.New("input cannot be empty")
)
