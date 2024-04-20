package model

import (
	"time"

	"github.com/google/uuid"
)

type Status string

type OrderID = uuid.UUID

type Order struct {
	ID        OrderID
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	HotelID    int       `json:"hotel_id"`
	RoomTypeID int       `json:"room_type_id"`
	UserEmail  string    `json:"email"`
	From       time.Time `json:"from"`
	To         time.Time `json:"to"`

	Status Status `json:"status"`
}

const (
	New Status = "new"

	NoRooms = "no_rooms"

	Booked Status = "booked"
	Paid   Status = "paid"

	FailedBook Status = "failedBook"
	FailedPay  Status = "failedPay"
)

var allStatuses = [...]Status{New, NoRooms, Booked, Paid, FailedBook, FailedPay}
