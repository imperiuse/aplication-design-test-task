package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	httpApi "aplication-design-test-task/internal/adapters/api/http"
	"aplication-design-test-task/internal/adapters/queue"
	"aplication-design-test-task/internal/adapters/queue/gochanqueue"
	"aplication-design-test-task/internal/adapters/storage/inmemory/storage"
	"aplication-design-test-task/internal/core/service/booking"
	"aplication-design-test-task/internal/logger"
	"aplication-design-test-task/migration"
)

// import _ "go.uber.org/automaxprocs" // todo good thing. Automatically set GOMAXPROCS to match Linux container CPU quota.

const (
	addr                      = "localhost:8080" // todo move this to env or config
	gracefullyShutdownTimeout = 5 * time.Second
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	log := logger.New()

	log.Info("App starting...")

	q := gochanqueue.NewChanQueue(log)
	log.Info("Queue is successfully created.")
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), gracefullyShutdownTimeout)
		defer cancel()

		if err := q.Close(ctx); err != nil {
			log.Error("Failed to close queue: %v ", err)
		}
	}()

	for _, topicName := range queue.AllTopics {
		if err := q.CreateTopic(ctx, topicName); err != nil {
			log.Error("Failed to create topics in Queue. err: %v ", err)
			os.Exit(1)
		}
	}

	store := storage.NewStorage()
	log.Info("Storage is successfully created.")
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), gracefullyShutdownTimeout)
		defer cancel()

		err := store.Close(ctx)
		if err != nil {
			log.Error("Failed to close Storage. err: %v ", err)
		} else {
			log.Info("Storage successfully closed")
		}
	}()

	err := migration.InitializeStorage(ctx, store)
	if err != nil {
		log.Error("Failed to init BookingService. err: %v ", err)
		os.Exit(2)
	}

	bookingService, err := booking.New(log, q, store)
	if err != nil {
		log.Error("Failed to init BookingService. err: %v ", err)
		os.Exit(3)
	}

	if err = bookingService.Run(ctx); err != nil {
		log.Error("Failed to Run BookingService. err: %v ", err)
		os.Exit(4)
	}

	httpServer := httpApi.NewServer(addr, log, q, bookingService)
	if err := httpServer.Run(ctx, gracefullyShutdownTimeout); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Error("Failed to run HTTP server: %v", err)
	}
	log.Info("Server exited gracefully.")

	log.Info("App finished.")
}
