package models

import "errors"

var (
	ErrTimeout            = errors.New("timeout")
	ErrReachedLimit       = errors.New("reached inflight limit")
	ErrReachedLimitPerKey = errors.New("reached inflight limit per key")
	ErrNilContext         = errors.New("context should not be empty")
	ErrNilInput           = errors.New("input cannot be empty")
	ErrMakingHash         = errors.New("failed making hash for function parameters")
)
