package eventbus

import (
	"context"
	"log"
	"sync"
)

// InMemoryEventBus is a simple event bus
type InMemoryEventBus struct {
	subscribers map[string][]EventHandler
	mu          sync.RWMutex
}

// NewInMemoryEventBus creates an event bus instance
func NewInMemoryEventBus() *InMemoryEventBus {
	return &InMemoryEventBus{
		subscribers: make(map[string][]EventHandler),
	}
}

// Publish sends an event to all subscribers
func (eb *InMemoryEventBus) Publish(ctx context.Context, event Event) error {
	eb.mu.RLock()
	defer eb.mu.RUnlock()

	handlers, ok := eb.subscribers[event.EventName()]
	if !ok {
		return nil
	}

	pCtx := context.WithValue(context.Background(), "name", "context in event bus")
	for _, handler := range handlers {
		go func(h EventHandler) {
			err := h.Handle(pCtx, event)
			if err != nil {
				log.Printf("Error handling event %s: %v", event.EventName(), err)
			}
		}(handler)
	}
	return nil
}

// Subscribe registers a handler for a specific event
func (eb *InMemoryEventBus) Subscribe(eventName string, handler EventHandler) {
	eb.mu.Lock()
	defer eb.mu.Unlock()

	eb.subscribers[eventName] = append(eb.subscribers[eventName], handler)
}
