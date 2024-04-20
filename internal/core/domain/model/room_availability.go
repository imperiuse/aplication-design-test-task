package model

import "time"

type RoomAvailabilityID = int

type RoomAvailability struct {
	ID         RoomAvailabilityID
	HotelID    int       `json:"hotel_id"`
	RoomTypeID int       `json:"room_type_id"`
	Date       time.Time `json:"date"`
	Quota      int       `json:"quota"`
}
