package mock

import (
	"context"

	"github.com/stretchr/testify/mock"

	"aplication-design-test-task/internal/core/domain/model"
)

type MockBookingService struct {
	mock.Mock
}

func (m *MockBookingService) Run(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockBookingService) StoreOrder(ctx context.Context, order model.Order) error {
	args := m.Called(ctx, order)
	return args.Error(0)
}

func (m *MockBookingService) GetOrder(ctx context.Context, id model.OrderID) (model.Order, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(model.Order), args.Error(1)
}

func (m *MockBookingService) GetListOrders(ctx context.Context) ([]model.Order, error) {
	args := m.Called(ctx)
	return args.Get(0).([]model.Order), args.Error(1)
}

func (m *MockBookingService) GetListRooms(ctx context.Context) ([]model.RoomAvailability, error) {
	args := m.Called(ctx)
	return args.Get(0).([]model.RoomAvailability), args.Error(1)
}
