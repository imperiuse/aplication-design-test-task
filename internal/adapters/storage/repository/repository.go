package repository

import "context"

// Storer is a generic interface for basic CRUD operations on storage where ID must be comparable.
type Storer[ID comparable, T any] interface {
	Create(context.Context, ID, T) error
	Read(context.Context, ID) (T, error)
	Update(context.Context, ID, T) error
	Delete(context.Context, ID) error
	List(context.Context) ([]T, error)
}
