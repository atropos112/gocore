package pubsub

import (
	. "github.com/atropos112/atrogolib/types"
)

// Publisher is an interface for publishing messages to a topic of sorts.
type Publisher interface {
	Publish(topic string, event Event, metadata *map[string]string) error
}

type Subscriber interface {
	Subscribe(topic string) (<-chan Event, error)
}
