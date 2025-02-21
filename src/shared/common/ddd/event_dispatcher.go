package ddd

import (
	"context"
	"log"
	"sync"
)

// DomainEvent represents a domain event
type DomainEvent interface {
	EventName() string
}

// EventHandler defines how an event should be handled
type EventHandler interface {
	Handle(ctx context.Context, event DomainEvent) error
}

// EventDispatcher is the centralized event dispatcher
type EventDispatcher interface {
	Register(eventType DomainEvent, handler EventHandler)
	Dispatch(ctx context.Context, event DomainEvent)
}

// eventDispatcher manages event listeners
type eventDispatcher struct {
	listeners map[string][]EventHandler
	mu        sync.RWMutex
}

// NewEventDispatcher creates a new dispatcher
func NewEventDispatcher() EventDispatcher {
	return &eventDispatcher{listeners: make(map[string][]EventHandler)}
}

// Register adds an event listener for a specific event type
func (d *eventDispatcher) Register(eventType DomainEvent, handler EventHandler) {
	d.mu.Lock()
	defer d.mu.Unlock()

	d.listeners[eventType.EventName()] = append(d.listeners[eventType.EventName()], handler)
}

// Dispatch fires an event to all listeners of that event type
func (d *eventDispatcher) Dispatch(ctx context.Context, event DomainEvent) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	handlers, exists := d.listeners[event.EventName()]
	if !exists {
		return
	}

	for _, handler := range handlers {
		// // Run asynchronously
		func(h EventHandler) {
			err := h.Handle(ctx, event)
			if err != nil {
				log.Printf("Error handling event %s: %v", event.EventName(), err)
			}
		}(handler)
	}

}
