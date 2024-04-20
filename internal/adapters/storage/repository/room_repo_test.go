package repository

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"aplication-design-test-task/internal/adapters/storage/repository/mock"
	"aplication-design-test-task/internal/core/domain/model"
	"aplication-design-test-task/internal/core/util"
)

func TestRoomRepository_StoreRoom(t *testing.T) {
	ctx := context.Background()
	room := model.RoomAvailability{ID: 1, HotelID: 101, RoomTypeID: 201, Quota: 5}
	mockStorer := new(mock.MockRoomStorer)
	repo := NewRoomRepository(mockStorer)

	mockStorer.On("Create", ctx, room.ID, room).Return(nil)

	err := repo.StoreRoom(ctx, room)

	assert.NoError(t, err)
	mockStorer.AssertExpectations(t)
}

func TestRoomRepository_GetRoom(t *testing.T) {
	ctx := context.Background()
	room := model.RoomAvailability{ID: 1, HotelID: 101, RoomTypeID: 201, Quota: 5}
	mockStorer := new(mock.MockRoomStorer)
	repo := NewRoomRepository(mockStorer)

	mockStorer.On("Read", ctx, room.ID).Return(room, nil)

	result, err := repo.GetRoom(ctx, room.ID)

	assert.NoError(t, err)
	assert.Equal(t, room, result)
	mockStorer.AssertExpectations(t)
}

func TestRoomRepository_GetAllRooms(t *testing.T) {
	ctx := context.Background()
	rooms := []model.RoomAvailability{
		{ID: 1, HotelID: 101, RoomTypeID: 201, Quota: 5},
		{ID: 2, HotelID: 102, RoomTypeID: 202, Quota: 3},
	}
	mockStorer := new(mock.MockRoomStorer)
	repo := NewRoomRepository(mockStorer)

	mockStorer.On("List", ctx).Return(rooms, nil)

	result, err := repo.GetAllRooms(ctx)

	assert.NoError(t, err)
	assert.Equal(t, rooms, result)
	mockStorer.AssertExpectations(t)
}

func TestRoomRepository_GetRoomsForHotelByRoomTypeAndDate(t *testing.T) {
	ctx := context.Background()
	rooms := []model.RoomAvailability{
		{ID: 1, HotelID: 101, RoomTypeID: 201, Date: time.Now(), Quota: 5},
	}
	mockStorer := new(mock.MockRoomStorer)
	repo := NewRoomRepository(mockStorer)
	today := time.Now()
	tomorrow := today.Add(24 * time.Hour)

	mockStorer.On("List", ctx).Return(rooms, nil)

	result, err := repo.GetRoomsForHotelByRoomTypeAndDate(ctx, 101, 201, today, tomorrow)

	require.NoError(t, err)
	require.Len(t, result, 1)
	assert.Equal(t, rooms[0], result[0])
	mockStorer.AssertExpectations(t)
}

func TestRoomRepository_GetRoomsForHotelByRoomTypeAndDate_ContextError(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	room := Room{
		ID:         1,
		HotelID:    101,
		RoomTypeID: 201,
		Date:       time.Now(),
		Quota:      5,
	}

	rooms := []Room{room}

	mockStorer := new(mock.MockRoomStorer)
	mockStorer.On("List", ctx).Return(rooms, nil)

	repo := NewRoomRepository(mockStorer)

	_, err := repo.GetRoomsForHotelByRoomTypeAndDate(ctx, 101, 201, time.Now(), time.Now().AddDate(0, 0, 1))

	assert.ErrorIs(t, err, context.Canceled)
}

func TestRoomRepository_GetRoomsForHotelByRoomTypeAndDate_ListError(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	storageErr := fmt.Errorf("storage err")
	mockStorer := new(mock.MockRoomStorer)
	mockStorer.On("List", ctx).Return([]Room{}, storageErr)

	repo := NewRoomRepository(mockStorer)

	r, err := repo.GetRoomsForHotelByRoomTypeAndDate(ctx, 101, 201, time.Now(), time.Now().AddDate(0, 0, 1))

	assert.ErrorIs(t, err, storageErr)
	assert.Nil(t, r)
}

func TestGetListRooms(t *testing.T) {
	ctx := context.Background()

	mockStorage := new(mock.MockRoomStorer)
	roomRepo := NewRoomRepository(mockStorage)

	expectedRooms := []model.RoomAvailability{
		{ID: 1},
		{ID: 2},
	}

	mockStorage.On("List", ctx).Return(expectedRooms, nil)

	rooms, err := roomRepo.GetListRooms(ctx)

	// Assert expectations
	mockStorage.AssertExpectations(t)
	assert.NoError(t, err, "Expected no error from GetListRooms")
	assert.Equal(t, expectedRooms, rooms, "Expected returned rooms to match the mock orders")
}

func TestUpdateRoom(t *testing.T) {
	mockStorage := new(mock.MockRoomStorer)
	roomRepo := NewRoomRepository(mockStorage)
	ctx := context.Background()

	const roomID = 1
	room := model.RoomAvailability{
		ID: 1,
	}

	mockStorage.On("Update", ctx, roomID, room).Return(nil)

	err := roomRepo.UpdateRoom(ctx, roomID, room)

	mockStorage.AssertExpectations(t)
	assert.NoError(t, err, "UpdateRoom should not return an error")
}

func TestGetRoomsForHotelByRoomTypeAndDate(t *testing.T) {
	mockStorage := new(mock.MockRoomStorer)
	roomRepo := NewRoomRepository(mockStorage)
	ctx := context.TODO()

	fromDate := util.NewDay(2024, 4, 1)
	toDate := util.NewDay(2024, 4, 7)

	mockRooms := []model.RoomAvailability{
		{ID: 1, HotelID: 1, RoomTypeID: 1, Date: fromDate},
		{ID: 2, HotelID: 1, RoomTypeID: 2, Date: fromDate.AddDate(0, 0, 1)}, // should be included
		{ID: 3, HotelID: 2, RoomTypeID: 1, Date: fromDate.AddDate(0, 0, 2)}, // different hotel ID
		{ID: 4, HotelID: 1, RoomTypeID: 2, Date: toDate.AddDate(0, 0, 1)},   // outside date range
		// Add more sample data if needed
	}

	mockStorage.On("List", ctx).Return(mockRooms, nil)

	expectedRooms := []model.RoomAvailability{
		mockRooms[1], // Only the second room matches all criteria
	}

	filteredRooms, err := roomRepo.GetRoomsForHotelByRoomTypeAndDate(ctx, 1, 2, fromDate, toDate)

	mockStorage.AssertExpectations(t)
	assert.NoError(t, err, "GetRoomsForHotelByRoomTypeAndDate should not return an error")
	assert.Equal(t, expectedRooms, filteredRooms, "GetRoomsForHotelByRoomTypeAndDate should return the correct filtered rooms")
}
