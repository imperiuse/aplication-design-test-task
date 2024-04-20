package migration

import (
	"context"

	"aplication-design-test-task/internal/adapters/storage"
	"aplication-design-test-task/internal/core/domain/model"
	"aplication-design-test-task/internal/core/util"
)

const (
	firstHotelID  = 1
	secondHotelID = 2
)

// InitializeStorage initializes the storage with predefined data
func InitializeStorage(ctx context.Context, store storage.Storage) error {

	const quotaTen = 10
	roomsHotelOne := []model.RoomAvailability{
		// ONE WEEK - ONE HOTEL WITH 1 types rooms and 10 quota
		{HotelID: firstHotelID, RoomTypeID: 1, Date: util.NewDay(2024, 4, 1), Quota: quotaTen},
		{HotelID: firstHotelID, RoomTypeID: 2, Date: util.NewDay(2024, 4, 1), Quota: quotaTen},
		{HotelID: firstHotelID, RoomTypeID: 3, Date: util.NewDay(2024, 4, 1), Quota: quotaTen},
		{HotelID: firstHotelID, RoomTypeID: 1, Date: util.NewDay(2024, 4, 2), Quota: quotaTen},
		{HotelID: firstHotelID, RoomTypeID: 2, Date: util.NewDay(2024, 4, 2), Quota: quotaTen},
		{HotelID: firstHotelID, RoomTypeID: 3, Date: util.NewDay(2024, 4, 2), Quota: quotaTen},
		{HotelID: firstHotelID, RoomTypeID: 1, Date: util.NewDay(2024, 4, 3), Quota: quotaTen},
		{HotelID: firstHotelID, RoomTypeID: 2, Date: util.NewDay(2024, 4, 3), Quota: quotaTen},
		{HotelID: firstHotelID, RoomTypeID: 3, Date: util.NewDay(2024, 4, 3), Quota: quotaTen},
		{HotelID: firstHotelID, RoomTypeID: 1, Date: util.NewDay(2024, 4, 4), Quota: quotaTen},
		{HotelID: firstHotelID, RoomTypeID: 2, Date: util.NewDay(2024, 4, 4), Quota: quotaTen},
		{HotelID: firstHotelID, RoomTypeID: 3, Date: util.NewDay(2024, 4, 4), Quota: quotaTen},
		{HotelID: firstHotelID, RoomTypeID: 1, Date: util.NewDay(2024, 4, 5), Quota: quotaTen},
		{HotelID: firstHotelID, RoomTypeID: 2, Date: util.NewDay(2024, 4, 5), Quota: quotaTen},
		{HotelID: firstHotelID, RoomTypeID: 3, Date: util.NewDay(2024, 4, 5), Quota: quotaTen},
		{HotelID: firstHotelID, RoomTypeID: 1, Date: util.NewDay(2024, 4, 6), Quota: quotaTen},
		{HotelID: firstHotelID, RoomTypeID: 2, Date: util.NewDay(2024, 4, 6), Quota: quotaTen},
		{HotelID: firstHotelID, RoomTypeID: 3, Date: util.NewDay(2024, 4, 6), Quota: quotaTen},
		{HotelID: firstHotelID, RoomTypeID: 1, Date: util.NewDay(2024, 4, 7), Quota: quotaTen},
		{HotelID: firstHotelID, RoomTypeID: 2, Date: util.NewDay(2024, 4, 7), Quota: quotaTen},
		{HotelID: firstHotelID, RoomTypeID: 3, Date: util.NewDay(2024, 4, 7), Quota: quotaTen},
	}

	roomsHotelTwo := []model.RoomAvailability{
		// ONE WEEK - ONE HOTEL WITH 1 types rooms and 10 quota
		{HotelID: secondHotelID, RoomTypeID: 1, Date: util.NewDay(2024, 4, 1), Quota: quotaTen},
		{HotelID: secondHotelID, RoomTypeID: 2, Date: util.NewDay(2024, 4, 1), Quota: quotaTen},
		{HotelID: secondHotelID, RoomTypeID: 3, Date: util.NewDay(2024, 4, 1), Quota: quotaTen},
		{HotelID: secondHotelID, RoomTypeID: 1, Date: util.NewDay(2024, 4, 2), Quota: quotaTen},
		{HotelID: secondHotelID, RoomTypeID: 2, Date: util.NewDay(2024, 4, 2), Quota: quotaTen},
		{HotelID: secondHotelID, RoomTypeID: 3, Date: util.NewDay(2024, 4, 2), Quota: quotaTen},
		{HotelID: secondHotelID, RoomTypeID: 1, Date: util.NewDay(2024, 4, 3), Quota: quotaTen},
		{HotelID: secondHotelID, RoomTypeID: 2, Date: util.NewDay(2024, 4, 3), Quota: quotaTen},
		{HotelID: secondHotelID, RoomTypeID: 3, Date: util.NewDay(2024, 4, 3), Quota: quotaTen},
		{HotelID: secondHotelID, RoomTypeID: 1, Date: util.NewDay(2024, 4, 4), Quota: quotaTen},
		{HotelID: secondHotelID, RoomTypeID: 2, Date: util.NewDay(2024, 4, 4), Quota: quotaTen},
		{HotelID: secondHotelID, RoomTypeID: 3, Date: util.NewDay(2024, 4, 4), Quota: quotaTen},
		{HotelID: secondHotelID, RoomTypeID: 1, Date: util.NewDay(2024, 4, 5), Quota: quotaTen},
		{HotelID: secondHotelID, RoomTypeID: 2, Date: util.NewDay(2024, 4, 5), Quota: quotaTen},
		{HotelID: secondHotelID, RoomTypeID: 3, Date: util.NewDay(2024, 4, 5), Quota: quotaTen},
		{HotelID: secondHotelID, RoomTypeID: 1, Date: util.NewDay(2024, 4, 6), Quota: quotaTen},
		{HotelID: secondHotelID, RoomTypeID: 2, Date: util.NewDay(2024, 4, 6), Quota: quotaTen},
		{HotelID: secondHotelID, RoomTypeID: 3, Date: util.NewDay(2024, 4, 6), Quota: quotaTen},
		{HotelID: secondHotelID, RoomTypeID: 1, Date: util.NewDay(2024, 4, 7), Quota: quotaTen},
		{HotelID: secondHotelID, RoomTypeID: 2, Date: util.NewDay(2024, 4, 7), Quota: quotaTen},
		{HotelID: secondHotelID, RoomTypeID: 3, Date: util.NewDay(2024, 4, 7), Quota: quotaTen},
	}

	for id, room := range append(roomsHotelOne, roomsHotelTwo...) {
		room.ID = id
		if err := store.GetRoomRepo().StoreRoom(ctx, room); err != nil {
			return err
		}
	}

	return nil
}
