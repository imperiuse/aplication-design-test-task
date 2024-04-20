package repository

import (
	"context"

	"aplication-design-test-task/internal/core/domain/model"
)

type (
	Order              = model.Order
	ReservationOrderID = model.OrderID
)

type OrderRepository struct {
	storage Storer[ReservationOrderID, Order]
}

func NewOrderRepository(store Storer[ReservationOrderID, Order]) *OrderRepository {
	return &OrderRepository{storage: store}
}

func (r *OrderRepository) StoreOrder(ctx context.Context, order Order) error {
	return r.storage.Create(ctx, order.ID, order)
}

func (r *OrderRepository) GetOrder(ctx context.Context, id ReservationOrderID) (Order, error) {
	return r.storage.Read(ctx, id)
}

func (r *OrderRepository) UpdateOrder(ctx context.Context, id ReservationOrderID, order Order) error {
	return r.storage.Update(ctx, id, order)
}

func (r *OrderRepository) GetListOrders(ctx context.Context) ([]Order, error) {
	return r.storage.List(ctx)
}
