package storage

import (
	"context"

	"aplication-design-test-task/internal/adapters/storage/repository"
)

type (
	Storage interface {
		BeginTx(ctx context.Context) (Transaction, error)

		GetOrderRepo() *repository.OrderRepository
		GetRoomRepo() *repository.RoomRepository

		// Repo[T any]()T // todo wait in future in Golang =)
		//  see more Repository pattern with Go generics -> github.com/imperiuse/golib/db/db.go

		Close(context.Context) error
	}

	Transaction interface {
		Commit() error
		Rollback() error

		Execute(op func() error, rollbackFunc func() error)
	}
)
