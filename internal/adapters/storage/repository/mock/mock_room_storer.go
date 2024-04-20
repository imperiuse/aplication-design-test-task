package mock

import (
	"context"

	"github.com/stretchr/testify/mock"

	"aplication-design-test-task/internal/core/domain/model"
)

type MockRoomStorer struct {
	mock.Mock
}

func (m *MockRoomStorer) Create(ctx context.Context, id int, room model.RoomAvailability) error {
	args := m.Called(ctx, id, room)
	return args.Error(0)
}

func (m *MockRoomStorer) Read(ctx context.Context, id int) (model.RoomAvailability, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(model.RoomAvailability), args.Error(1)
}

func (m *MockRoomStorer) Update(ctx context.Context, id int, room model.RoomAvailability) error {
	args := m.Called(ctx, id, room)
	return args.Error(0)
}

func (m *MockRoomStorer) Delete(ctx context.Context, id int) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockRoomStorer) List(ctx context.Context) ([]model.RoomAvailability, error) {
	args := m.Called(ctx)
	return args.Get(0).([]model.RoomAvailability), args.Error(1)
}
