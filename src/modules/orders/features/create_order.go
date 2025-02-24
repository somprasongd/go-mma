package features

import (
	"context"
	"fmt"
	"go-mma/modules/orders/dtos"
	"go-mma/modules/orders/model"
	"go-mma/modules/orders/repository"
	"go-mma/shared/common/errs"
	"go-mma/shared/common/eventbus"
	"go-mma/shared/common/mediator"
	"go-mma/shared/common/storage/db/transactor"
	custCmd "go-mma/shared/contracts/customer_contracts/commands"
	custQuery "go-mma/shared/contracts/customer_contracts/queries"
	"go-mma/shared/messaging"
	"log"
)

type CreateOrderCommand struct {
	*dtos.CreateOrderRequest
}

func (r *CreateOrderCommand) Validate() error {
	if r.CustomerID <= 0 {
		return fmt.Errorf("customer ID must be greater than 0")
	}
	if r.OrderTotal <= 0 {
		return fmt.Errorf("order total must be greater than 0")
	}
	return nil
}

type CreateOrderResult struct {
	*dtos.CreateOrderResponse
}

type createOrderHandler struct {
	transactor transactor.Transactor
	orderRepo  repository.OrderRepository
	eventbus   eventbus.EventBus
}

func NewCreateOrderCommand(transactor transactor.Transactor, orderRepo repository.OrderRepository, eventbus eventbus.EventBus) *createOrderHandler {
	return &createOrderHandler{
		transactor: transactor,
		orderRepo:  orderRepo,
		eventbus:   eventbus,
	}
}

func (h *createOrderHandler) Handle(ctx context.Context, cmd *CreateOrderCommand) (*CreateOrderResult, error) {
	// validate request
	if err := cmd.Validate(); err != nil {
		return nil, errs.NewValidationError(err.Error())
	}

	var orderId int
	var orderCreatedIntegrationEvent *messaging.OrderCreatedIntegrationEvent
	err := h.transactor.WithinTransaction(ctx, func(ctx context.Context) error {
		_, err := mediator.Send[*custCmd.ReserveCreditCommand, *mediator.NoResponse](
			ctx,
			newReserveCreditCommand(cmd),
		)
		if err != nil {
			return err
		}

		order := model.NewOrder(cmd.CustomerID, cmd.OrderTotal)
		err = h.orderRepo.Create(ctx, order)
		if err != nil {
			log.Println(err)
			return errs.NewDatabaseFailureError(err.Error())
		}

		customer, err := mediator.Send[*custQuery.GetCustomerByIDQuery, *custQuery.GetCustomerByIDResult](
			ctx,
			&custQuery.GetCustomerByIDQuery{ID: cmd.CustomerID},
		)
		if err != nil {
			return err
		}

		fmt.Println(customer)

		orderCreatedIntegrationEvent = messaging.NewOrderCreatedIntegrationEvent(
			order.ID,
			order.CustomerID,
			order.OrderTotal,
			customer.Email,
		)

		fmt.Println(orderCreatedIntegrationEvent)

		orderId = order.ID
		return nil
	})

	if err != nil {
		return nil, err
	}

	// Publish the event
	h.eventbus.Publish(ctx, orderCreatedIntegrationEvent)

	return &CreateOrderResult{CreateOrderResponse: &dtos.CreateOrderResponse{ID: orderId}}, nil
}

func newReserveCreditCommand(cmd *CreateOrderCommand) *custCmd.ReserveCreditCommand {
	return &custCmd.ReserveCreditCommand{
		ID:           cmd.CustomerID,
		CreditAmount: cmd.OrderTotal,
	}
}
