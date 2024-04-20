package model

import (
	"time"

	"github.com/google/uuid"
)

type (
	Payment = struct {
		ID        uuid.UUID `json:"id"`
		OrderID   OrderID   `json:"order_id"`
		CreatedAt time.Time `json:"createdAt"`
		PaidAt    time.Time `json:"paidAt"`
		IsPaid    bool      `json:"isPaid"`
		// other
	}

	SuccessPaymentEvent = struct{ _ [0]int }
	FailedPaymentEvent  = struct{}
)
