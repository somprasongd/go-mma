package eventhandlers

import (
	"context"
	"errors"
	"fmt"
	"go-mma/modules/orders/events"
	"go-mma/shared/common/ddd"
	"go-mma/shared/common/eventbus"
	"log"

	customerContracts "go-mma/shared/contracts/customer_contracts"
)

type orderPlacedEventHandler struct {
	custSrv  customerContracts.CreditManagement
	eventbus eventbus.EventBus
}

func NewOrderPlacedEventHandler(custSrv customerContracts.CreditManagement) *orderPlacedEventHandler {
	return &orderPlacedEventHandler{custSrv: custSrv}
}

func (h *orderPlacedEventHandler) Handle(ctx context.Context, event ddd.DomainEvent) error {
	// Implement the logic to handle the order placed event
	orderPlacedEvent, ok := event.(*events.OrderPlacedEvent)
	if !ok {
		return errors.New("invalid order placed event")
	}

	custInfo, err := h.custSrv.GetCustomerByID(ctx, orderPlacedEvent.CustomerID)
	if err != nil {
		return fmt.Errorf("failed to get customer information: %w", err)
	}

	// Example: Log the order details
	log.Printf("Sending email: %s, Order %d placed by customer %d with total %d\n", custInfo.Email, orderPlacedEvent.ID, orderPlacedEvent.CustomerID, orderPlacedEvent.OrderTotal)

	// orderPlacedIntegrationEvent := &events.OrderPlacedIntegrationEvent{
	// 	ID:         orderPlacedEvent.ID,
	// 	CustomerID: orderPlacedEvent.CustomerID,
	// 	OrderTotal: orderPlacedEvent.OrderTotal,
	// 	Email:      custInfo.Email,
	// }

	// h.eventbus.Publish(ctx, orderPlacedIntegrationEvent)

	return nil
}
