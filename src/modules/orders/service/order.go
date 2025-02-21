package service

import (
	"context"
	"log"

	"go-mma/modules/orders/model"
	"go-mma/modules/orders/repository"
	"go-mma/shared/common/errs"
	"go-mma/shared/common/eventbus"
	"go-mma/shared/common/storage/db/transactor"
	"go-mma/shared/messaging"

	customerContracts "go-mma/shared/contracts/customer_contracts"
	orderContracts "go-mma/shared/contracts/order_contracts"
)

type orderService struct {
	transactor transactor.Transactor
	custServ   customerContracts.CreditManagement
	orderRepo  repository.OrderRepository
	eventbus   eventbus.EventBus
}

func NewOrderService(transactor transactor.Transactor, custServ customerContracts.CreditManagement, orderRepo repository.OrderRepository, eventbus eventbus.EventBus) orderContracts.OrderService {
	return &orderService{
		transactor: transactor,
		custServ:   custServ,
		orderRepo:  orderRepo,
		eventbus:   eventbus,
	}
}

var (
	ErrOrderNotFound = errs.NewResourceNotFoundError("the order with given id was not found")
)

func (s *orderService) CreateOrder(ctx context.Context, req *orderContracts.CreateOrderRequest) (int, error) {
	// validate request
	if err := req.Validate(); err != nil {
		return 0, errs.NewValidationError(err.Error())
	}

	var orderId int
	var orderCreatedIntegrationEvent *messaging.OrderCreatedIntegrationEvent
	err := s.transactor.WithinTransaction(ctx, func(ctx context.Context) error {
		err := s.custServ.ReserveCredit(ctx, req.CustomerID, req.OrderTotal)
		if err != nil {
			return err
		}

		order := model.NewOrder(req.CustomerID, req.OrderTotal)
		err = s.orderRepo.Create(ctx, order)
		if err != nil {
			log.Println(err)
			return errs.NewDatabaseFailureError(err.Error())
		}

		customer, err := s.custServ.GetCustomerByID(ctx, req.CustomerID)
		if err != nil {
			return err
		}

		orderCreatedIntegrationEvent = messaging.NewOrderCreatedIntegrationEvent(
			order.ID,
			order.CustomerID,
			order.OrderTotal,
			customer.Email,
		)

		orderId = order.ID
		return nil
	})

	if err != nil {
		return 0, err
	}

	// Publish the event
	s.eventbus.Publish(ctx, orderCreatedIntegrationEvent)

	return orderId, nil
}

func (s *orderService) CancelOrder(ctx context.Context, id int) error {
	return s.transactor.WithinTransaction(ctx, func(ctx context.Context) error {
		order, err := s.orderRepo.FindByID(ctx, id)
		if err != nil {
			log.Println(err)
			return errs.NewDatabaseFailureError(err.Error())
		}

		if order == nil {
			return ErrOrderNotFound
		}

		if err := s.orderRepo.Cancel(ctx, order.ID); err != nil {
			log.Println(err)
			return errs.NewDatabaseFailureError(err.Error())
		}

		if err := s.custServ.ReleaseCredit(ctx, order.CustomerID, order.OrderTotal); err != nil {
			return err
		}

		return nil
	})
}
