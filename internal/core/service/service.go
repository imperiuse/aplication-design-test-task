package service

import (
	"context"

	"aplication-design-test-task/internal/core/domain/model"
)

type (
	BookingService interface {
		Run(context.Context) error
		GetOrder(context.Context, model.OrderID) (model.Order, error)

		GetListOrders(ctx context.Context) ([]model.Order, error)
		GetListRooms(ctx context.Context) ([]model.RoomAvailability, error)
	}

	PaymentService interface {
		Run(context.Context) error
	}

	Notification interface {
		Run(context.Context) error
	}
)
