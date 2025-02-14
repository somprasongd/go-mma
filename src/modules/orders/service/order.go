package service

import (
	"context"
	"log"

	"go-mma/modules/orders/model"
	"go-mma/modules/orders/repository"
	"go-mma/shared/common/errs"
	"go-mma/shared/common/storage/db/transactor"

	customerContracts "go-mma/shared/contracts/customer_contracts"
	notificationContracts "go-mma/shared/contracts/notification_contracts"
	orderContracts "go-mma/shared/contracts/order_contracts"
)

type orderService struct {
	transactor  transactor.Transactor
	custServ    customerContracts.CreditManagement
	orderRepo   repository.OrderRepository
	notiService notificationContracts.NotificationService
}

func NewOrderService(transactor transactor.Transactor, custServ customerContracts.CreditManagement, orderRepo repository.OrderRepository, notiService notificationContracts.NotificationService) orderContracts.OrderService {
	return &orderService{
		transactor:  transactor,
		custServ:    custServ,
		orderRepo:   orderRepo,
		notiService: notiService,
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

		s.notiService.SendEmail(customer.Email, "Order Created", map[string]any{
			"order_id": order.ID,
			"total":    order.OrderTotal,
		})

		orderId = order.ID
		return nil
	})

	if err != nil {
		return 0, err
	}

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
