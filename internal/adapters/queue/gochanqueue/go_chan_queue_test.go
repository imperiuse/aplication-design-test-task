package gochanqueue

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"aplication-design-test-task/internal/adapters/queue"
	"aplication-design-test-task/internal/logger"
)

func TestChanQueueSynchronousPublish(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	q := NewChanQueue(logger.New())
	defer func() { _ = q.Close(ctx) }()

	topicName := queue.Topic("testTopic")
	message := "test message"

	require.NoError(t, q.CreateTopic(ctx, topicName), "could not create topic")

	done := make(chan struct{})
	go func() {
		defer close(done)
		assert.NoError(t, q.Publish(ctx, topicName, message), "Publish failed")
	}()

	select {
	case <-done:
		// Passed
	case <-time.After(1 * time.Second):
		assert.Fail(t, "Publish timed out")
	}

	msgChan, err := q.Subscribe(ctx, topicName)
	require.NoError(t, err, "Subscribe failed")
	receivedMsg := <-msgChan
	assert.Equal(t, message, receivedMsg, "Received message does not match expected")
}

func TestChanQueueSynchronousPublishMaxLenCapReached(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	q := NewChanQueue(logger.New())

	topicName := queue.Topic("testTopic")
	message := "test message"

	require.NoError(t, q.CreateTopic(ctx, topicName), "could not create topic")

	// This loop should not block since the channel has a buffer
	for i := 0; i < QueueMaxLength; i++ {
		assert.NoError(t, q.Publish(ctx, topicName, message), "Publish failed on buffered channel")
	}

	// This publish should block or fail, check it with a timeout to confirm it blocks
	done := make(chan bool)
	go func() {
		err := q.Publish(ctx, topicName, message) // this should block or fail due to full buffer
		if err != nil {
			done <- true // If it fails as expected, signal completion
		} else {
			done <- false // If it does not fail, signal an issue
		}
	}()

	// Test if the code blocks or not, expecting to timeout indicating full channel blocking
	select {
	case success := <-done:
		assert.False(t, success, "Publish should not succeed, the channel buffer should be full")
	case <-time.After(1 * time.Second):
		// Timeout expected, indicating that the publish operation is properly blocking
	}
}

func TestChanQueueAsynchronousPublish(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	q := NewChanQueue(logger.New())
	defer func() { _ = q.Close(ctx) }()

	topicName := queue.Topic("testTopic")
	message := "test message"

	require.NoError(t, q.CreateTopic(ctx, topicName), "could not create topic")

	msgChan, err := q.Subscribe(ctx, topicName)
	require.NoError(t, err, "Subscribe failed")

	err = q.AsyncPublish(ctx, topicName, message)
	assert.NoError(t, err, "AsyncPublish failed")

	select {
	case receivedMsg := <-msgChan:
		assert.Equal(t, message, receivedMsg, "Received message does not match expected")
	case <-time.After(10 * time.Second):
		assert.Fail(t, "Did not receive message from AsyncPublish")
	}
}

func TestChanQueueAsynchronousPublishMaxCapReached(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	q := NewChanQueue(logger.New())

	topicName := queue.Topic("testTopic")
	message := "test message"

	// Create topic
	require.NoError(t, q.CreateTopic(ctx, topicName), "Error while creating topic")

	// Fill the channel to its maximum capacity
	for i := 0; i < QueueMaxLength; i++ {
		require.NoError(t, q.AsyncPublish(ctx, topicName, message), "Error while asynchronously publishing to channel")
	}

	// Attempt to publish one more message which should not block, but return a buffer full error
	err := q.AsyncPublish(ctx, topicName, message)
	assert.Error(t, err, "Expected an error for buffer being full")

	// Test receiving message
	msgChan, err := q.Subscribe(ctx, topicName)
	require.NoError(t, err, "Subscribe failed")

	for i := 0; i < QueueMaxLength; i++ {
		select {
		case receivedMsg := <-msgChan:
			assert.Equal(t, message, receivedMsg, "Mismatch in message received")
		case <-time.After(1 * time.Second):
			t.Fatal("Timeout while trying to receive the message")
		}
	}

	// After processing all messages, no messages should be left and channel should be not blocked
	select {
	case extraMsg := <-msgChan:
		t.Errorf("Should not have received an extra message: %s", extraMsg)
	case <-time.After(50 * time.Millisecond):
		// No message should be received, this is expected
	}
}

func TestCreateAndDeleteTopic(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	q := NewChanQueue(logger.New())
	defer func() { _ = q.Close(ctx) }()

	name := queue.Topic("testTopic")

	assert.NoError(t, q.CreateTopic(ctx, name), "creating topic should not fail")
	assert.NoError(t, q.CreateTopic(ctx, name), "creating topic should not fail")
	assert.NoError(t, q.DeleteTopic(ctx, name), "deleting existing topic should not fail")
	assert.NoError(t, q.DeleteTopic(ctx, name), "deleting non-existent topic should not fail")
}

func TestPublishToNonExistentTopic(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	q := NewChanQueue(logger.New())
	defer func() { _ = q.Close(ctx) }()

	name := queue.Topic("nonExistentTopic")

	// Attempt to sync publish to a non-existent topic
	assert.Error(t, q.Publish(ctx, name, struct{}{}), "expected error when publishing to non-existent topic[sync]")

	// Attempt to async publish to a non-existent topic
	assert.Error(t, q.AsyncPublish(ctx, name, struct{}{}), "expected error when publishing to non-existent topic[async]")
}

func TestSubscribeToNonExistentTopic(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	q := NewChanQueue(logger.New())
	defer func() { _ = q.Close(ctx) }()

	name := queue.Topic("nonExistentTopic")

	// Attempt to subscribe to a non-existent topic
	_, err := q.Subscribe(ctx, name)
	assert.Error(t, err, "expected error when subscribing to non-existent topic")
}

func TestChanQueueWithCancelledContext(t *testing.T) {
	q := NewChanQueue(logger.New())

	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	defer q.Close(ctx) // Ensure cleanup

	// Simulate a delay before cancellation
	go func() {
		time.Sleep(100 * time.Millisecond) // Delay before cancel
		cancel()
	}()

	// Wait until the context is done to proceed with operations
	<-ctx.Done()

	topicName := queue.Topic("testTopic")
	message := "test message"

	// Test CreateTopic with cancelled context
	err := q.CreateTopic(ctx, topicName)
	assert.Error(t, err, "CreateTopic should fail with cancelled context")
	assert.Equal(t, context.Canceled, err, "Error should be context.Canceled")

	// Test DeleteTopic with cancelled context
	err = q.DeleteTopic(ctx, topicName)
	assert.Error(t, err, "DeleteTopic should fail with cancelled context")
	assert.Equal(t, context.Canceled, err, "Error should be context.Canceled")

	// Test Publish with cancelled context
	err = q.Publish(ctx, topicName, message)
	assert.Error(t, err, "Publish should fail with cancelled context")
	assert.Equal(t, context.Canceled, err, "Error should be context.Canceled")

	// Test AsyncPublish with cancelled context
	err = q.AsyncPublish(ctx, topicName, message)
	assert.Error(t, err, "AsyncPublish should fail with cancelled context")
	assert.Equal(t, context.Canceled, err, "Error should be context.Canceled")

	// Test Subscribe with cancelled context
	_, err = q.Subscribe(ctx, topicName)
	assert.Error(t, err, "Subscribe should fail with cancelled context")
	assert.Equal(t, context.Canceled, err, "Error should be context.Canceled")
}
