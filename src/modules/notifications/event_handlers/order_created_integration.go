package eventhandlers

import (
	"context"
	"errors"
	"go-mma/modules/notifications/features"
	"go-mma/shared/common/eventbus"
	"go-mma/shared/common/mediator"
	"go-mma/shared/messaging"
)

type orderCreatedIntegrationEventHandler struct {
}

func NewOrderCreatedIntegrationEventHandler() *orderCreatedIntegrationEventHandler {
	return &orderCreatedIntegrationEventHandler{}
}

func (h *orderCreatedIntegrationEventHandler) Handle(ctx context.Context, event eventbus.Event) error {
	odrCreatedEvent, ok := event.(*messaging.OrderCreatedIntegrationEvent)

	if !ok {
		return errors.New("Invalid event type: OrderCreatedIntegrationEvent")
	}

	to := odrCreatedEvent.Email
	subject := "Order Created Notification"
	payload := map[string]any{
		"orderId": odrCreatedEvent.ID,
		"amount":  odrCreatedEvent.OrderTotal,
	}

	mediator.Send[*features.SendEmailCommand, *mediator.NoResponse](ctx, &features.SendEmailCommand{
		To:      to,
		Subject: subject,
		Payload: payload,
	})

	return nil
}
