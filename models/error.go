package models

import "errors"

var (
	ErrTimeout            = errors.New("timeout")
	ErrReachedLimit       = errors.New("reached inflight limit")
	ErrReachedLimitPerKey = errors.New("reached inflight limit per key")
	ErrNilContext         = errors.New("context should not be empty")
	ErrListEmpty          = errors.New("input list cannot be empty")
)
