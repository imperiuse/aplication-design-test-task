package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	m "github.com/stretchr/testify/mock"

	"aplication-design-test-task/internal/adapters/api/http/mock"
	"aplication-design-test-task/internal/adapters/queue"
	"aplication-design-test-task/internal/adapters/storage"
	"aplication-design-test-task/internal/core/domain/model"
	"aplication-design-test-task/internal/core/util"
	"aplication-design-test-task/internal/logger"
)

func TestPingHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/ping", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		_, _ = fmt.Fprintf(writer, "pong")
	})

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code, "handler returned wrong status code")

	expected := "pong"
	assert.Equal(t, expected, rr.Body.String(), "handler returned unexpected body")
}

func TestPostReservationOrderHandler(t *testing.T) {
	log := logger.New()
	queueMock := new(mock.MockQueue)
	handler := postReservationOrderHandler(log, queueMock)

	validOrderRequest := orderReservationRequest{
		HotelID:    1,
		RoomTypeID: 1,
		UserEmail:  "test@example.com",
		From:       util.NewDay(2024, 4, 1),
		To:         util.NewDay(2024, 4, 7),
	}

	testCases := []struct {
		name           string
		requestBody    interface{}
		prepareMock    func()
		expectedStatus int
		expectedBody   map[string]string
		expectedError  string
	}{
		{
			name:        "Valid Request",
			requestBody: validOrderRequest,
			prepareMock: func() {
				queueMock.On("Publish", m.Anything, queue.ReservedOrderRequest, m.Anything).Return(nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   map[string]string{"order_id": "some", "status": "received"},
		},
		{
			name: "Invalid Email",
			requestBody: orderReservationRequest{
				HotelID:    1,
				RoomTypeID: 1,
				UserEmail:  "invalid-email",
				From:       time.Now(),
				To:         time.Now().Add(24 * time.Hour),
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Invalid email format",
		},
		{
			name: "Invalid Date Range",
			requestBody: orderReservationRequest{
				HotelID:    1,
				RoomTypeID: 1,
				UserEmail:  "test@example.com",
				From:       time.Now().Add(24 * time.Hour),
				To:         time.Now(),
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "From date must be before To date",
		},
		// More test cases...
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			bodyBytes, _ := json.Marshal(tc.requestBody)
			r := httptest.NewRequest("POST", "/api/v1/order", bytes.NewBuffer(bodyBytes))
			r.Header.Add("Content-Type", "application/json")
			w := httptest.NewRecorder()

			if tc.prepareMock != nil {
				tc.prepareMock()
			}

			handler(w, r)

			res := w.Result()
			defer res.Body.Close()

			assert.Equal(t, tc.expectedStatus, res.StatusCode)

			if tc.expectedBody != nil {
				var response map[string]string
				err := json.NewDecoder(res.Body).Decode(&response)
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedBody["status"], response["status"])
				assert.NotNil(t, response["order_id"])
			}

			if tc.expectedError != "" {
				responseBody, _ := io.ReadAll(res.Body)
				assert.Contains(t, string(responseBody), tc.expectedError)
			}

			queueMock.AssertExpectations(t)
		})
	}
}

func TestGetReservationOrderHandler2(t *testing.T) {
	log := logger.New()
	validUUID := uuid.New()
	notFoundUUID := uuid.Nil

	order := model.Order{
		ID:         validUUID,
		HotelID:    1,
		RoomTypeID: 2,
		UserEmail:  "test@example.com",
		From:       util.NewDay(2024, 4, 1),
		To:         util.NewDay(2024, 4, 7),
		Status:     "new",
	}

	bookingServiceMock := new(mock.MockBookingService)

	tests := []struct {
		name             string
		orderID          string
		prepareMock      func()
		expectedStatus   int
		expectedResponse *orderReservationResponse
		expectedError    string
	}{
		{
			name:    "Valid Order",
			orderID: validUUID.String(),
			prepareMock: func() {
				bookingServiceMock.On("GetOrder", m.Anything, validUUID).Return(order, nil)
			},
			expectedStatus: http.StatusOK,
			expectedResponse: &orderReservationResponse{ID: validUUID,
				HotelID:    order.HotelID,
				RoomTypeID: order.RoomTypeID,
				UserEmail:  order.UserEmail,
				From:       order.From,
				To:         order.To,
				Status:     string(order.Status)},
		},
		{
			name:           "Order ID Missing",
			orderID:        "",
			expectedStatus: http.StatusNotFound,
			expectedError:  "",
		},
		{
			name:           "Invalid Order ID",
			orderID:        "invalid-uuid",
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Invalid Order ID format",
		},
		{
			name:    "Order Not Found",
			orderID: notFoundUUID.String(),
			prepareMock: func() {
				bookingServiceMock.On("GetOrder", m.Anything, notFoundUUID).Return(model.Order{}, storage.ErrNotFound)
			},
			expectedStatus: http.StatusNotFound,
			expectedError:  "Order not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.prepareMock != nil {
				tt.prepareMock()
			}

			// Set up chi router with the route // standard http router did not support params.....
			r := chi.NewRouter()
			r.Get("/api/v1/order/{id}", getReservationOrderHandler(log, bookingServiceMock))

			// Create an HTTP request with the order ID in the URL path
			req := httptest.NewRequest("GET", "/api/v1/order/"+tt.orderID, nil)
			w := httptest.NewRecorder()

			// Serve the HTTP request using our chi router
			r.ServeHTTP(w, req)

			// Check the HTTP response
			res := w.Result()
			defer res.Body.Close()

			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedResponse != nil {
				var response orderReservationResponse
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, *tt.expectedResponse, response)
			}

			if tt.expectedError != "" {
				assert.Contains(t, w.Body.String(), tt.expectedError)
			}
		})
	}
}
