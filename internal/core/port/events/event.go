package events

import "aplication-design-test-task/internal/core/domain/model"

type (
	ReservationOrderEvent = model.Order // todo can be dived to separate struct Event <-> DTO

	PaymentRequest = model.Payment

	// todo
	SuccessPaymentEvent = model.SuccessPaymentEvent
	FailedPaymentEvent  = model.FailedPaymentEvent
)
