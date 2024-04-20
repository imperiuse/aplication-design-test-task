package storage

import "errors"

var (
	ErrNotFound            = errors.New("not found")
	ErrDuplicateConstraint = errors.New("already exists")
)
