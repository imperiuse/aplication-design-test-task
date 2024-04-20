package repository

import (
	"context"
	"time"

	"aplication-design-test-task/internal/core/domain/model"
	"aplication-design-test-task/internal/core/util"
)

type Room = model.RoomAvailability

type RoomRepository struct {
	storage Storer[int, Room]
}

func NewRoomRepository(store Storer[int, Room]) *RoomRepository {
	return &RoomRepository{storage: store}
}

func (r *RoomRepository) StoreRoom(ctx context.Context, room Room) error {
	return r.storage.Create(ctx, room.ID, room)
}

func (r *RoomRepository) GetRoom(ctx context.Context, id int) (Room, error) {
	return r.storage.Read(ctx, id)
}

func (r *RoomRepository) UpdateRoom(ctx context.Context, id int, room Room) error {
	return r.storage.Update(ctx, id, room)
}

func (r *RoomRepository) GetAllRooms(ctx context.Context) ([]Room, error) {
	return r.storage.List(ctx)
}

func (r *RoomRepository) GetRoomsForHotelByRoomTypeAndDate(
	ctx context.Context,
	hotelID,
	roomTypeID int,
	fromDate time.Time,
	toDate time.Time,
) ([]Room, error) {
	allRooms, err := r.storage.List(ctx)
	if err != nil {
		return nil, err
	}

	var filteredRooms []Room
	for _, room := range allRooms {
		if err = ctx.Err(); err != nil {
			return nil, err
		}

		if room.HotelID == hotelID && room.RoomTypeID == roomTypeID && util.IsDayBetween(room.Date, fromDate, toDate) {
			filteredRooms = append(filteredRooms, room)
		}
	}
	return filteredRooms, nil
}

// GetListRooms retrieves all Rooms from the repository
func (r *RoomRepository) GetListRooms(ctx context.Context) ([]Room, error) {
	return r.storage.List(ctx)
}
