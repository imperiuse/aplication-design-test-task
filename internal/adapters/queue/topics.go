package queue

// Success flow.
const (
	ReservedOrderRequest  Topic = "ReservedOrderRequest"
	PaymentRequest        Topic = "PaymentRequest"
	NotificationRequest   Topic = "NotificationRequest"
	SuccessPaymentProcess Topic = "SuccessPaymentProcess"
)

// Error flow.
const (
	FailedOrder          Topic = "FailedOrder"
	FailedPayment        Topic = "FailedPayment"
	FailedPaymentProcess Topic = "FailedPaymentProcess"
)

var AllTopics = [...]Topic{
	ReservedOrderRequest,
	FailedOrder,
	PaymentRequest,
	FailedPayment,
	NotificationRequest,
	FailedPaymentProcess,
	SuccessPaymentProcess,
}
