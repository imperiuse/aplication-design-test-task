package gochanqueue

import (
	"context"
	"fmt"
	"sync"

	"aplication-design-test-task/internal/adapters/queue"
	"aplication-design-test-task/internal/logger"
)

const QueueMaxLength = 10 // magic number, Set for demonstration purposes.

type (
	topic = queue.Topic
	msg   = queue.Msg

	queueMap = map[topic]chan msg // for small size [] (instead of map) will be faster.
)

type ChanQueue struct {
	log logger.Logger
	m   sync.RWMutex
	q   queueMap
}

// NewChanQueue initializes a new channel-based message queue.
func NewChanQueue(log logger.Logger) *ChanQueue {
	return &ChanQueue{
		log: log,
		m:   sync.RWMutex{},
		q:   make(queueMap),
	}
}

// CreateTopic creates a new topic with a buffered channel.
func (c *ChanQueue) CreateTopic(ctx context.Context, name topic) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		c.m.Lock()
		defer c.m.Unlock()

		if _, found := c.q[name]; !found {
			c.q[name] = make(chan msg, QueueMaxLength)
			c.log.Info("topic `%s` is successfully created.", name)
		}
		return nil
	}
}

// DeleteTopic deletes a topic and closes its channel.
func (c *ChanQueue) DeleteTopic(ctx context.Context, name topic) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		c.m.Lock()
		defer c.m.Unlock()

		if ch, ok := c.q[name]; ok {
			close(ch)
			delete(c.q, name)
			c.log.Info("topic `%s` is successfully deleted", name)
		}
		return nil
	}
}

// Publish sends a message to the specified topic channel, blocking until the message is sent.
func (c *ChanQueue) Publish(ctx context.Context, name topic, m msg) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		c.m.RLock()
		ch, ok := c.q[name]
		c.m.RUnlock()

		if !ok {
			c.log.Error("topic `%s` not exists", name)
			return queue.TopicNotExists
		}

		select {
		case ch <- m:
			return nil
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

// AsyncPublish sends a message to the specified topic channel without blocking.
func (c *ChanQueue) AsyncPublish(ctx context.Context, name topic, m msg) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		c.m.RLock()
		ch, ok := c.q[name]
		c.m.RUnlock()

		if !ok {
			c.log.Error("topic `%s` not exists", name)
			return queue.TopicNotExists
		}

		select {
		case ch <- m:
			return nil
		case <-ctx.Done():
			return ctx.Err()
		default:
			c.log.Error("channel buffer is full. topic `%s`", name)
			return fmt.Errorf("channel buffer is full")
		}
	}
}

// Subscribe returns a channel from which messages can be received.
func (c *ChanQueue) Subscribe(ctx context.Context, name topic) (<-chan msg, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		c.m.RLock()
		ch, ok := c.q[name]
		c.m.RUnlock()

		if !ok {
			c.log.Error("topic `%s` not exists", name)
			return nil, queue.TopicNotExists
		}
		return ch, nil
	}
}

// Close shuts down the queue and closes all channels.
func (c *ChanQueue) Close(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		c.m.Lock()
		defer c.m.Unlock()

		for name, ch := range c.q {
			close(ch)
			delete(c.q, name)
			c.log.Info("topic `%s` successfully deleted", name)
		}
		c.log.Info("ChanQueue successfully closed")
		return nil
	}
}
