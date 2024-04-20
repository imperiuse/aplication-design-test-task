package mock

import (
	"context"

	"github.com/stretchr/testify/mock"

	"aplication-design-test-task/internal/core/domain/model"
)

// MockOrderStorer is a mock type for the Storer interface
type MockOrderStorer struct {
	mock.Mock
}

func (m *MockOrderStorer) Create(ctx context.Context, id model.OrderID, order model.Order) error {
	args := m.Called(ctx, id, order)
	return args.Error(0)
}

func (m *MockOrderStorer) Read(ctx context.Context, id model.OrderID) (model.Order, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(model.Order), args.Error(1)
}

func (m *MockOrderStorer) Update(ctx context.Context, id model.OrderID, order model.Order) error {
	args := m.Called(ctx, id, order)
	return args.Error(0)
}

func (m *MockOrderStorer) Delete(ctx context.Context, id model.OrderID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockOrderStorer) List(ctx context.Context) ([]model.Order, error) {
	args := m.Called(ctx)
	return args.Get(0).([]model.Order), args.Error(1)
}
