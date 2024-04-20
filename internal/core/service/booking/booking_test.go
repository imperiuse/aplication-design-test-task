package booking

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"

	"aplication-design-test-task/internal/adapters/queue"
	"aplication-design-test-task/internal/adapters/queue/gochanqueue"
	"aplication-design-test-task/internal/adapters/storage"
	instorage "aplication-design-test-task/internal/adapters/storage/inmemory/storage"
	"aplication-design-test-task/internal/core/port/events"
	"aplication-design-test-task/internal/core/service"
	"aplication-design-test-task/internal/core/util"
	"aplication-design-test-task/internal/logger"
	"aplication-design-test-task/migration"
)

type BookingServiceSuite struct {
	suite.Suite
	Service     service.BookingService
	ServiceImpl *bookingService
	Context     context.Context
	Logger      logger.Logger
	Queue       queue.Queue
	Storage     storage.Storage
}

func (suite *BookingServiceSuite) SetupSuite() {
	suite.Logger = logger.New()
	suite.Context = context.Background()

	suite.Queue = gochanqueue.NewChanQueue(suite.Logger)

	for _, topicName := range queue.AllTopics {
		suite.NoError(suite.Queue.CreateTopic(suite.Context, topicName))
	}

	suite.Storage = instorage.NewStorage()

	var err error
	suite.ServiceImpl, err = New(suite.Logger, suite.Queue, suite.Storage) // for white box testing
	suite.Service = suite.ServiceImpl                                      // blackBox testing
	suite.NoError(err)

	suite.NoError(suite.Service.Run(suite.Context))
}

func (suite *BookingServiceSuite) TearDownSuite() {
	suite.NoError(suite.Queue.Close(suite.Context))
	suite.NoError(suite.Storage.Close(suite.Context))
}

func (suite *BookingServiceSuite) BeforeTest(suiteName, testName string) {
	suite.Storage.Close(suite.Context)
	suite.Storage = instorage.NewStorage()
	suite.NoError(migration.InitializeStorage(suite.Context, suite.Storage))
}

func (suite *BookingServiceSuite) AfterTest(suiteName, testName string) {
}

func TestBookingServiceSuite(t *testing.T) {
	suite.Run(t, new(BookingServiceSuite))
}

func (suite *BookingServiceSuite) TestBookingService_ReservationOrderEventHandler() {
	suite.NotNil(suite.Service)

	orderID := uuid.New()

	baseOrderEvent := ReservationOrder{
		ID:         orderID,
		CreatedAt:  util.NewDay(2024, 04, 01),
		UpdatedAt:  util.NewDay(2024, 04, 01),
		HotelID:    1,
		RoomTypeID: 1,
		UserEmail:  "ars-saz@ya.ru",
		From:       util.NewDay(2024, 04, 01),
		To:         util.NewDay(2024, 04, 07),
		Status:     "",
	}

	testCases := []struct {
		name  string
		event events.ReservationOrderEvent
	}{
		{
			name:  "One day order",
			event: baseOrderEvent,
		},
		// todo
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			suite.ServiceImpl.ReservationOrderEventHandler(suite.Context, tc.event)
			time.Sleep(time.Second)
		},
		)
	}
}
