package http

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/mail"
	"strconv"
	"time"

	"github.com/google/uuid"

	"aplication-design-test-task/internal/adapters/queue"
	"aplication-design-test-task/internal/adapters/storage"
	"aplication-design-test-task/internal/core/domain/model"
	"aplication-design-test-task/internal/core/port/events"
	"aplication-design-test-task/internal/core/service"
	"aplication-design-test-task/internal/logger"
)

type orderReservationRequest struct {
	ID         uuid.UUID `json:"-"`
	HotelID    int       `json:"hotel_id"`
	RoomTypeID int       `json:"room_type_id"`
	UserEmail  string    `json:"email"`
	From       time.Time `json:"from"`
	To         time.Time `json:"to"`
	PromoCode  string    `json:"promo_code"`
	// todo other options...
}

type orderReservationResponse struct {
	ID         uuid.UUID `json:"id"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	Status     string    `json:"status"`
	HotelID    int       `json:"hotel_id"`
	RoomTypeID int       `json:"room_type_id"`
	UserEmail  string    `json:"email"`
	From       time.Time `json:"from"`
	To         time.Time `json:"to"`
}

func validateEmail(email string) bool {
	_, err := mail.ParseAddress(email) // https://pkg.go.dev/net/mail
	return err == nil
}

func postReservationOrderHandler(log logger.Logger, q queue.Queue) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Info("postReservationOrderHandler")

		var orderRequest orderReservationRequest

		// Validation stuff
		err := json.NewDecoder(r.Body).Decode(&orderRequest)
		if err != nil {
			log.Error("Invalid request body: %v", err)
			http.Error(w, "Invalid request body: please ensure JSON is properly formatted", http.StatusBadRequest)
			return
		}

		if !validateEmail(orderRequest.UserEmail) {
			log.Error("Invalid email format: %s", orderRequest.UserEmail)
			http.Error(w, "Invalid email format: please provide a valid email", http.StatusBadRequest)
			return
		}

		if orderRequest.HotelID <= 0 || orderRequest.RoomTypeID <= 0 {
			log.Error("Invalid hotel or room type ID")
			http.Error(w, "Invalid hotel or room type ID: IDs must be positive integers", http.StatusBadRequest)
			return
		}

		if !orderRequest.From.Before(orderRequest.To) {
			log.Error("From date must be before To date")
			http.Error(w, "From date must be before To date", http.StatusBadRequest)
			return
		}

		orderRequest.ID = uuid.New() // https://en.wikipedia.org/w/index.php?title=Universally_unique_identifier&oldid=755882275#Random_UUID_probability_of_duplicates

		orderReservationEvent := events.ReservationOrderEvent{
			ID:         orderRequest.ID,
			CreatedAt:  time.Now().UTC(),
			UpdatedAt:  time.Now().UTC(),
			HotelID:    orderRequest.HotelID,
			RoomTypeID: orderRequest.RoomTypeID,
			UserEmail:  orderRequest.UserEmail,
			From:       orderRequest.From,
			To:         orderRequest.To,
		}

		err = q.Publish(r.Context(), queue.ReservedOrderRequest, orderReservationEvent)
		if err != nil {
			log.Error("Failed to publish the order request: %v", err)
			http.Error(w, "Failed to publish the order request: internal server error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		err = json.NewEncoder(w).Encode(map[string]string{
			"order_id": orderRequest.ID.String(),
			"status":   "received",
		})
		if err != nil {
			log.Error("Failed to encode the response: %v", err)
			http.Error(w, "Failed to encode the response", http.StatusInternalServerError)
		}
	}
}

func getReservationOrderHandler(log logger.Logger, bookingService service.BookingService) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Info("getReservationOrderHandler")

		orderIDStr := r.PathValue("id")
		if orderIDStr == "" {
			log.Error("Order ID is required")
			http.Error(w, "Order ID is required", http.StatusBadRequest)
			return
		}

		orderID, err := uuid.Parse(orderIDStr)
		if err != nil {
			log.Error("Invalid Order ID format")
			http.Error(w, "Invalid Order ID format", http.StatusBadRequest)
			return
		}

		order, err := bookingService.GetOrder(r.Context(), orderID)
		if err != nil {
			if errors.Is(err, storage.ErrNotFound) {
				log.Error("Order with id: `%s` not found", orderID)

				http.Error(w, "Order not found. Probably need wait a little bit... ", http.StatusNotFound)
				return
			}

			log.Error("Failed to retrieve order")
			http.Error(w, "Failed to retrieve order", http.StatusInternalServerError)
			return
		}

		response := orderReservationResponse{
			ID:         order.ID,
			CreatedAt:  order.CreatedAt,
			UpdatedAt:  order.UpdatedAt,
			Status:     string(order.Status),
			HotelID:    order.HotelID,
			RoomTypeID: order.RoomTypeID,
			UserEmail:  order.UserEmail,
			From:       order.From,
			To:         order.To,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		err = json.NewEncoder(w).Encode(response)
		if err != nil {
			log.Error("Failed to encode the response: %v", err)
			http.Error(w, "Failed to send the response", http.StatusInternalServerError)
		}
	}
}

// debug handler, not for production, in real life use checking env.
func registerDebugHandlers(mux *http.ServeMux, bookingService service.BookingService) {

	mux.HandleFunc("GET /api/v1/order", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		orders, _ := bookingService.GetListOrders(r.Context())
		_ = json.NewEncoder(w).Encode(orders)
	})

	mux.HandleFunc("GET /api/v1/room", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		rooms, _ := bookingService.GetListRooms(r.Context())
		_ = json.NewEncoder(w).Encode(rooms)
	})

	mux.HandleFunc("GET /api/v1/room/{hotel_id}/{room_type_id}", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		var (
			hotelID    = 0
			roomTypeID = 0
		)

		hotelIDstr := r.PathValue("hotel_id")
		if hotelIDstr != "" {
			hotelID = atoi(hotelIDstr)
		}

		roomTypeIDstr := r.PathValue("room_type_id")
		if roomTypeIDstr != "" {
			roomTypeID = atoi(roomTypeIDstr)
		}

		allRooms, _ := bookingService.GetListRooms(r.Context())

		rooms := make([]model.RoomAvailability, 0, len(allRooms))

		for _, room := range allRooms {
			if hotelID != 0 {
				if room.HotelID != hotelID {
					continue
				}
			}

			if roomTypeID != 0 {
				if room.RoomTypeID != roomTypeID {
					continue
				}
			}

			rooms = append(rooms, room)
		}

		_ = json.NewEncoder(w).Encode(rooms)
	})
}

func atoi(s string) int {
	val, _ := strconv.Atoi(s)
	return val
}
