package mock

import (
	"context"

	"github.com/stretchr/testify/mock"

	"aplication-design-test-task/internal/adapters/queue"
)

type MockQueue struct {
	mock.Mock
}

func (m *MockQueue) CreateTopic(ctx context.Context, topic queue.Topic) error {
	args := m.Called(ctx, topic)
	return args.Error(0)
}

func (m *MockQueue) DeleteTopic(ctx context.Context, topic queue.Topic) error {
	args := m.Called(ctx, topic)
	return args.Error(0)
}

func (m *MockQueue) Publish(ctx context.Context, topic queue.Topic, message interface{}) error {
	args := m.Called(ctx, topic, message)
	return args.Error(0)
}

func (m *MockQueue) AsyncPublish(ctx context.Context, topic queue.Topic, message interface{}) error {
	args := m.Called(ctx, topic, message)
	return args.Error(0)
}

func (m *MockQueue) Subscribe(ctx context.Context, topic queue.Topic) (<-chan queue.Msg, error) {
	args := m.Called(ctx, topic)
	return args.Get(0).(<-chan queue.Msg), args.Error(1)
}

func (m *MockQueue) Close(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}
