package booking

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"aplication-design-test-task/internal/adapters/queue"
	"aplication-design-test-task/internal/adapters/storage"
	"aplication-design-test-task/internal/core/domain/model"
	"aplication-design-test-task/internal/core/port/events"
	"aplication-design-test-task/internal/core/service/booking/worker"
	"aplication-design-test-task/internal/core/util"
	"aplication-design-test-task/internal/logger"
)

const workerCnt = 1 // for now magic number

type (
	ID         int
	HotelID    = ID
	RoomTypeID = ID
	UUID       = uuid.UUID

	ReservationOrderID = model.OrderID
	ReservationOrder   = model.Order

	RoomAvailability = model.RoomAvailability

	bookingService struct {
		log     logger.Logger
		q       queue.Queue
		storage storage.Storage
		workers []bookingWorker
	}

	bookingWorker interface {
		Run(context.Context, <-chan queue.Msg)
	}
)

func New(log logger.Logger, q queue.Queue, s storage.Storage) (*bookingService, error) {
	service := &bookingService{
		log:     log,
		q:       q,
		storage: s,
		workers: make([]bookingWorker, 0, workerCnt),
	}

	for range workerCnt {
		service.workers = append(service.workers, worker.New(log, service))
	}

	return service, nil
}

func (s *bookingService) Run(ctx context.Context) error {
	for _, w := range s.workers {
		const topicName = queue.ReservedOrderRequest

		// TODO: If workerCnt > 1, additional logic is needed to support partitioning, queue groups, or similar features.
		ch, err := s.q.Subscribe(ctx, topicName)
		if err != nil {
			return fmt.Errorf("could not subscribe to topic %s. err: %v", topicName, err)
		}

		w.Run(ctx, ch)
	}

	return nil
}

// ReservationOrderEventHandler - provide CORE logic of Booking service!
func (s *bookingService) ReservationOrderEventHandler(ctx context.Context, event events.ReservationOrderEvent) {
	// TODO !!!DISCLAIMER!!!
	// All actions with the database must be within a transaction, and all events must be sent only at the "right" point, which will be suitable.
	// We should understand that at any line of code, we might encounter a failure, so we must maintain a consistent state in the database for our orders and hotels.
	// IN IDEAL situation use
	//tx, err := s.storage.BeginTx(ctx)
	//	if err != nil {
	//		s.log.Error("[bookingService.ReservationOrderEventHandler] Failed to start transaction: %v", err)
	//		return
	//	}
	//	defer func() {
	//		if p := recover(); p != nil {
	//			_ = tx.Rollback()
	//			// panic(p)  // do not like use panic mechanism in prod system.
	//		} else if err != nil {
	//			_ = tx.Rollback()
	//		} else {
	//			err = tx.Commit()
	//		}
	//	}()
	//newOrder, err := s.createOrderFromEvent(ctx, event)
	//if err != nil {
	//	return err // newOrder creation failed, return the error
	//}
	//
	//if err := s.storeNewOrder(ctx, tx, newOrder); err != nil {
	//	return err // storing new order failed, return the error
	//}
	//
	//if err := s.processRoomAvailability(ctx, tx, newOrder, event); err != nil {
	//	return err // processing room availability failed, return the error
	//}
	//
	//if err := s.publishPaymentRequestEvent(ctx, newOrder); err != nil {
	//	return err // publishing payment request failed, return the error
	//}

	// FOR NOW, FOR CUSTOM IN MEMORY STORAGE I use own super simple transaction mechanism.

	newOrder := ReservationOrder{
		ID:         event.ID,
		CreatedAt:  event.CreatedAt,
		UpdatedAt:  time.Now().UTC(),
		HotelID:    event.HotelID,
		RoomTypeID: event.RoomTypeID,
		UserEmail:  event.UserEmail,
		From:       event.From,
		To:         event.To,
		Status:     model.New,
	}

	// todo properly handle db error (re-try or other policy...)
	if err := s.storage.GetOrderRepo().StoreOrder(ctx, newOrder); err != nil { // optionally. can catch error.Is(err, storage.ErrDuplicateConstraint) for Upsert goal in the future....
		s.log.Error("[bookingService.ReservationOrderEventHandler] Failed to store new order: %v", err)
	} else {
		s.log.Info("[bookingService.ReservationOrderEventHandler] Stored new order: %v", newOrder)
	}

	// Wrap transaction stuff in separate func
	isBookedSuccessfully := func() bool {
		tx, err := s.storage.BeginTx(ctx)
		if err != nil {
			s.log.Error("[bookingService.ReservationOrderEventHandler] Failed to start transaction: %v", err)
			return false
		}

		defer func() {
			if err := tx.Commit(); err != nil {
				s.log.Error("[bookingService.ReservationOrderEventHandler] Failed to commit transaction: %v", err)
				return
			}

			s.log.Info("[bookingService.ReservationOrderEventHandler] Success commit transaction")
		}()

		processedOrder := newOrder
		processedOrder.Status = model.NoRooms

		defer func() {
			tx.Execute(
				func() error {
					err := s.storage.GetOrderRepo().UpdateOrder(ctx, processedOrder.ID, processedOrder)
					if err != nil {
						s.log.Error("[bookingService.ReservationOrderEventHandler] Failed to update processed order: %v", err)
					} else {
						s.log.Info("[bookingService.ReservationOrderEventHandler] Updated processed order: %v", processedOrder)
					}

					return err
				},

				func() error {
					s.log.Info("[bookingService.ReservationOrderEventHandler] Rollback processed order to new: %v", newOrder)
					return s.storage.GetOrderRepo().UpdateOrder(ctx, newOrder.ID, newOrder)
				},
			)
		}()

		rooms, err := s.storage.GetRoomRepo().GetRoomsForHotelByRoomTypeAndDate(ctx, event.HotelID, event.RoomTypeID, event.From, event.To)
		if err != nil {
			s.log.Error("[bookingService.ReservationOrderEventHandler] Failed to retrieve rooms information: %v", err)

			processedOrder.Status = model.FailedBook
			return false // early exit (defer for Order update status to model.NoRooms)
		}

		if len(rooms) == 0 || len(rooms) < len(util.DaysBetween(event.From, event.To)) {
			s.log.Error("[bookingService.ReservationOrderEventHandler] No rooms for period")

			return false // early exit (defer for Order update status to model.NoRooms)
		}

		s.log.Info("[bookingService.ReservationOrderEventHandler] Retrieve rooms information for order."+
			" event.HotelID: %d, event.RoomTypeID: %d, event.From: %v, event.To: %v.  rooms: %s",
			event.HotelID, event.RoomTypeID, event.From, event.To, rooms)

		for _, room := range rooms {
			oldRoomInfo := room

			if room.Quota > 0 { // todo if user will can book more the one room, must be change this place
				room.Quota-- // change here too: e.g. room.Quota -= some

				tx.Execute(
					func() error {
						err := s.storage.GetRoomRepo().UpdateRoom(ctx, room.ID, room)
						if err != nil {
							s.log.Error("[bookingService.ReservationOrderEventHandler] Failed to update room quota: %v", err)
						} else {
							s.log.Info("[bookingService.ReservationOrderEventHandler] Updated room quota: %v", room)
						}

						return err
					},

					func() error {
						processedOrder.Status = model.FailedBook
						return s.storage.GetRoomRepo().UpdateRoom(ctx, room.ID, oldRoomInfo)
					})
			} else {
				s.log.Info("[bookingService.ReservationOrderEventHandler] "+
					"No room quota event.HotelID: %d, event.RoomTypeID: %d for date: %v. Booking process stopped!",
					room.HotelID, room.RoomTypeID, room.Date)

				return false // early exit (defer for Order update status to model.NoRooms)
			}
		}

		s.log.Info("[bookingService.ReservationOrderEventHandler] Room available for all days. " +
			"Try to set up processOrder.Status = model.Booked")

		processedOrder.Status = model.Booked

		return true
	}()

	if !isBookedSuccessfully {
		return // not send paymentRequestMsg if not successfully booke
	}

	paymentRequestMsg := events.PaymentRequest{
		ID:        uuid.New(),
		OrderID:   newOrder.ID,
		CreatedAt: time.Now().UTC(),
		PaidAt:    time.Time{},
		IsPaid:    false,
	}
	err := s.q.AsyncPublish(ctx, queue.PaymentRequest, paymentRequestMsg)
	if err != nil {
		s.log.Error("[bookingService.ReservationOrderEventHandler] Failed to publish PaymentRequest msg: %v", err)
	}

	s.log.Info("[bookingService.ReservationOrderEventHandler] Published PaymentRequest msg: %v", paymentRequestMsg)
}

func (s *bookingService) SuccessPaymentEventHandler(ctx context.Context, event events.SuccessPaymentEvent) {
	//TODO implement me
	panic("implement me")
}

func (s *bookingService) FailedPaymentEventHandler(ctx context.Context, event events.FailedPaymentEvent) {
	//TODO implement me
	panic("implement me")
}

func (s *bookingService) GetOrder(ctx context.Context, id ReservationOrderID) (ReservationOrder, error) {
	return s.storage.GetOrderRepo().GetOrder(ctx, id)
}

func (s *bookingService) GetListOrders(ctx context.Context) ([]ReservationOrder, error) {
	return s.storage.GetOrderRepo().GetListOrders(ctx)
}

func (s *bookingService) GetListRooms(ctx context.Context) ([]RoomAvailability, error) {
	return s.storage.GetRoomRepo().GetListRooms(ctx)
}
