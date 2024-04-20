package queue

import "context"

type (
	Topic string
	Msg   = any
)

// Queue represents an interface for a message queue system, which allows
// creation of topics, publishing, and subscription to messages.
type Queue interface {
	CreateTopic(context.Context, Topic) error
	DeleteTopic(context.Context, Topic) error

	Publish(context.Context, Topic, Msg) error
	AsyncPublish(context.Context, Topic, Msg) error
	Subscribe(context.Context, Topic) (<-chan Msg, error)

	Close(context.Context) error
}

// todo future
// Msg = Serializable
//type Serializable interface {
//	json.Marshaler
//	json.Unmarshaler
//}
