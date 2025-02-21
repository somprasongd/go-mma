package eventhandlers

import (
	"context"
	"errors"
	"go-mma/modules/notifications/service"
	"go-mma/shared/common/eventbus"
	"go-mma/shared/messaging"
)

type orderCreatedIntegrationEventHandler struct {
	svc service.NotificationService
}

func NewOrderCreatedIntegrationEventHandler(svc service.NotificationService) *orderCreatedIntegrationEventHandler {
	return &orderCreatedIntegrationEventHandler{
		svc: svc,
	}
}

func (h *orderCreatedIntegrationEventHandler) Handle(ctx context.Context, event eventbus.Event) error {
	odrCreatedEvent, ok := event.(*messaging.OrderCreatedIntegrationEvent)

	if !ok {
		return errors.New("Invalid event type: OrderCreatedIntegrationEvent")
	}

	to := "customer@example.com"
	subject := "Order Created Notification"
	payload := map[string]any{
		"orderId": odrCreatedEvent.ID,
		"amount":  odrCreatedEvent.OrderTotal,
	}

	h.svc.SendEmail(to, subject, payload)

	return nil
}
