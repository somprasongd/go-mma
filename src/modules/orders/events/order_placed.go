package events

import "go-mma/modules/orders/model"

type OrderPlacedEvent struct {
	model.Order
}

func (e *OrderPlacedEvent) EventName() string {
	return "order_placed"
}

func NewOrderPlacedEvent(o model.Order) *OrderPlacedEvent {
	return &OrderPlacedEvent{o}
}
