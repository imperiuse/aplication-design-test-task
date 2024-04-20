package repository

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	m "github.com/stretchr/testify/mock"

	"aplication-design-test-task/internal/adapters/storage/repository/mock"
	"aplication-design-test-task/internal/core/domain/model"
)

func TestOrderRepository_StoreOrder(t *testing.T) {
	ctx := context.Background()
	order := model.Order{ID: uuid.UUID{}, RoomTypeID: 123}
	mockStorer := new(mock.MockOrderStorer)
	repo := NewOrderRepository(mockStorer)

	mockStorer.On("Create", ctx, order.ID, order).Return(nil)

	err := repo.StoreOrder(ctx, order)

	assert.NoError(t, err)
	mockStorer.AssertExpectations(t)
}

func TestOrderRepository_GetOrder(t *testing.T) {
	ctx := context.Background()
	order := model.Order{ID: uuid.UUID{}, RoomTypeID: 123}
	mockStorer := new(mock.MockOrderStorer)
	repo := NewOrderRepository(mockStorer)

	mockStorer.On("Read", ctx, order.ID).Return(order, nil)

	result, err := repo.GetOrder(ctx, order.ID)

	assert.NoError(t, err)
	assert.Equal(t, order, result)
	mockStorer.AssertExpectations(t)
}

func TestGetListOrders(t *testing.T) {
	ctx := context.Background()

	mockStorage := new(mock.MockOrderStorer)
	orderRepo := NewOrderRepository(mockStorage)

	// Define the expected slice of orders
	expectedOrders := []model.Order{
		{ID: uuid.New()},
		{ID: uuid.New()},
	}

	// Setup expectations
	mockStorage.On("List", ctx).Return(expectedOrders, nil)

	orders, err := orderRepo.GetListOrders(ctx)

	// Assert expectations
	mockStorage.AssertExpectations(t)
	assert.NoError(t, err, "Expected no error from GetListOrders")
	assert.Equal(t, expectedOrders, orders, "Expected returned orders to match the mock orders")
}

func TestUpdateOrder(t *testing.T) {
	mockStorage := new(mock.MockOrderStorer)
	orderRepo := NewOrderRepository(mockStorage)
	ctx := context.Background()

	orderID := uuid.New()
	order := model.Order{
		ID:         orderID,
		HotelID:    1,
		RoomTypeID: 2,
		UserEmail:  "test@example.com",
		From:       time.Now(),
		To:         time.Now().Add(24 * time.Hour),
		Status:     model.Booked,
	}

	testCases := []struct {
		name          string
		orderID       uuid.UUID
		order         model.Order
		mockReturnErr error
		expectedErr   error
	}{
		{
			name:          "successful update",
			orderID:       orderID,
			order:         order,
			mockReturnErr: nil,
			expectedErr:   nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockStorage.On("Update", m.Anything, tc.orderID, tc.order).Return(tc.mockReturnErr)

			err := orderRepo.UpdateOrder(ctx, tc.orderID, tc.order)

			mockStorage.AssertExpectations(t)

			if tc.expectedErr != nil {
				assert.Error(t, err)
				assert.EqualError(t, err, tc.expectedErr.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
