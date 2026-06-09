package store

import "errors"

var (
	ErrNotFound     = errors.New("not found")
	ErrInsufficient = errors.New("insufficient stock")
)
