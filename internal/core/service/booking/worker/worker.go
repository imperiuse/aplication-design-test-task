package worker

import (
	"context"

	"github.com/google/uuid"

	"aplication-design-test-task/internal/adapters/queue"
	"aplication-design-test-task/internal/core/port/events"
	"aplication-design-test-task/internal/logger"
)

type eventHandler interface {
	ReservationOrderEventHandler(context.Context, events.ReservationOrderEvent)
	SuccessPaymentEventHandler(context.Context, events.SuccessPaymentEvent)
	FailedPaymentEventHandler(context.Context, events.FailedPaymentEvent)
}

type worker struct {
	log logger.Logger
	id  uuid.UUID
	eventHandler
}

func New(log logger.Logger, eh eventHandler) *worker {
	return &worker{id: uuid.New(), log: log, eventHandler: eh}
}

func (w *worker) Run(ctx context.Context, ch <-chan queue.Msg) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				w.log.Info("[bookingWorker: %v] ctx.Done(). finished", w.id)
				return

			case msg := <-ch:
				switch event := msg.(type) {
				case events.ReservationOrderEvent:
					w.log.Info("[bookingWorker: %v] received ReservationOrderEvent: %+v", w.id, event)
					w.ReservationOrderEventHandler(ctx, event)
				case events.SuccessPaymentEvent:
					w.log.Info("[bookingWorker: %v] received SuccessPaymentEvent: %+v", w.id, event)
					w.SuccessPaymentEventHandler(ctx, event)
				case events.FailedPaymentEvent:
					w.log.Info("[bookingWorker: %v] received FailedPaymentEvent: %+v", w.id, event)
					w.FailedPaymentEventHandler(ctx, event)
				case nil:
					continue
				default:
					w.log.Error("[bookingWorker: %v] received unknown msg: %+v", w.id, msg)
				}
			}
		}
	}()
}
