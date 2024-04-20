package inmemory

import (
	"context"
	"sync"

	"aplication-design-test-task/internal/adapters/storage"
)

// InMemoryStorage is an in-memory implementation of the Storer interface using a map.
type InMemoryStorage[ID comparable, T any] struct {
	sync.RWMutex
	store map[ID]T
}

// NewInMemoryStorage creates a new instance of InMemoryStorage.
func NewInMemoryStorage[ID comparable, T any]() *InMemoryStorage[ID, T] {
	return &InMemoryStorage[ID, T]{
		RWMutex: sync.RWMutex{},
		store:   make(map[ID]T),
	}
}

func (m *InMemoryStorage[ID, T]) Create(ctx context.Context, id ID, item T) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	m.Lock()
	defer m.Unlock()
	if _, exists := m.store[id]; exists {
		return storage.ErrDuplicateConstraint
	}
	m.store[id] = item
	return nil
}

func (m *InMemoryStorage[ID, T]) Read(ctx context.Context, id ID) (T, error) {
	select {
	case <-ctx.Done():
		return *new(T), ctx.Err()
	default:
	}

	m.RLock()
	defer m.RUnlock()
	item, exists := m.store[id]
	if !exists {
		return *new(T), storage.ErrNotFound
	}
	return item, nil
}

func (m *InMemoryStorage[ID, T]) Update(ctx context.Context, id ID, item T) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	m.Lock()
	defer m.Unlock()
	if _, exists := m.store[id]; !exists {
		return storage.ErrNotFound
	}
	m.store[id] = item
	return nil
}

func (m *InMemoryStorage[ID, T]) Delete(ctx context.Context, id ID) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	m.Lock()
	defer m.Unlock()
	if _, exists := m.store[id]; !exists {
		return storage.ErrNotFound
	}
	delete(m.store, id)
	return nil
}

func (m *InMemoryStorage[ID, T]) List(ctx context.Context) ([]T, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	m.RLock()
	defer m.RUnlock()
	items := make([]T, 0, len(m.store))
	for _, item := range m.store {
		items = append(items, item)
	}
	return items, nil
}
