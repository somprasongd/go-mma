package messaging

const (
	OrderCreatedIntegrationEventName = "order_created"
)

type OrderCreatedIntegrationEvent struct {
	ID         int    `json:"id"`
	CustomerID int    `json:"customer_id"`
	OrderTotal int    `json:"order_total"`
	Email      string `json:"email"`
}

func NewOrderCreatedIntegrationEvent(id int, customerID int, orderTotal int, email string) *OrderCreatedIntegrationEvent {
	return &OrderCreatedIntegrationEvent{
		ID:         id,
		CustomerID: customerID,
		OrderTotal: orderTotal,
		Email:      email,
	}
}

func (e *OrderCreatedIntegrationEvent) EventName() string {
	return OrderCreatedIntegrationEventName
}
