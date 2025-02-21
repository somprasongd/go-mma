package eventbus

import "context"

// Event represents a domain or integration event
type Event interface {
	EventName() string
}

// EventHandler defines how an event should be handled
type EventHandler interface {
	Handle(ctx context.Context, event Event) error
}

// EventBus is the centralized event dispatcher
type EventBus interface {
	Publish(ctx context.Context, event Event) error
	Subscribe(eventName string, handler EventHandler)
}
