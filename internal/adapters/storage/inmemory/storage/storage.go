package storage

import (
	"context"

	s "aplication-design-test-task/internal/adapters/storage"
	"aplication-design-test-task/internal/adapters/storage/inmemory"
	"aplication-design-test-task/internal/adapters/storage/repository"
	"aplication-design-test-task/internal/adapters/storage/transaction"
	"aplication-design-test-task/internal/core/domain/model"
)

type storage struct {
	orderRepo *repository.OrderRepository
	roomRepo  *repository.RoomRepository
}

func NewStorage() *storage {
	innMemStoreForReservationOrders := inmemory.NewInMemoryStorage[model.OrderID, model.Order]()
	innMemStoreForRoomAvailability := inmemory.NewInMemoryStorage[model.RoomAvailabilityID, model.RoomAvailability]()

	return &storage{
		orderRepo: repository.NewOrderRepository(innMemStoreForReservationOrders),
		roomRepo:  repository.NewRoomRepository(innMemStoreForRoomAvailability),
	}
}

func (s *storage) BeginTx(ctx context.Context) (s.Transaction, error) {
	return transaction.New(ctx), nil
}

func (s *storage) GetOrderRepo() *repository.OrderRepository {
	return s.orderRepo
}

func (s *storage) GetRoomRepo() *repository.RoomRepository {
	return s.roomRepo
}

func (s *storage) Close(_ context.Context) error {
	return nil
}
