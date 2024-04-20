package http

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"aplication-design-test-task/internal/adapters/queue"
	"aplication-design-test-task/internal/core/service"
	"aplication-design-test-task/internal/logger"
)

type server struct {
	addr   string
	log    logger.Logger
	mux    http.Handler
	server *http.Server
}

func NewServer(addr string, log logger.Logger, q queue.Queue, bookingService service.BookingService) *server {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /api/v1/ping", func(writer http.ResponseWriter, request *http.Request) {
		_, _ = fmt.Fprintf(writer, "pong")
	})

	mux.HandleFunc("GET /api/v1/order/{id}", getReservationOrderHandler(log, bookingService))
	mux.HandleFunc("POST /api/v1/order/", postReservationOrderHandler(log, q))
	// TODO: payments "ping-back" handlers

	registerDebugHandlers(mux, bookingService)

	return &server{
		addr: addr,
		log:  log,
		mux:  mux,
		server: &http.Server{
			Addr:    addr,
			Handler: mux,
		},
	}
}

func (s *server) Run(ctx context.Context, gracefullyShutdownTimeout time.Duration) error {
	go func() {
		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), gracefullyShutdownTimeout)
		defer cancel()

		err := s.server.Shutdown(shutdownCtx) // Gracefully shuts down the server
		if err != nil {
			s.log.Error("shutdown http server: %v", err)
		}
	}()

	s.log.Info("Start http server: " + s.addr)
	return s.server.ListenAndServe()
}
