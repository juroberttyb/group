/*
Package models defines all shared data models in this system
all models defined here should be Capitalized
*/
package models

type Result[T any] struct {
	Value T
	Err   error
}
