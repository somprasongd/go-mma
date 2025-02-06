package service

import (
	"context"
	"go-mma/modules/order/dtos"
	"go-mma/modules/order/model"
	"go-mma/modules/order/repository"
	"go-mma/util/errs"
	"go-mma/util/transactor"
	"log"

	custRepo "go-mma/modules/customer/repository"
	notiServ "go-mma/modules/notification/service"
)

type OrderService interface {
	CreateOrder(ctx context.Context, req *dtos.CreateOrderRequest) (int, error)
	CancelOrder(ctx context.Context, id int) error
}

type orderService struct {
	transactor  transactor.Transactor
	custRepo    custRepo.CustomerRepository
	orderRepo   repository.OrderRepository
	notiService notiServ.NotificationService
}

func NewOrderService(transactor transactor.Transactor, custRepo custRepo.CustomerRepository, orderRepo repository.OrderRepository, notiService notiServ.NotificationService) OrderService {
	return &orderService{
		transactor:  transactor,
		custRepo:    custRepo,
		orderRepo:   orderRepo,
		notiService: notiService,
	}
}

var (
	ErrCustomerNotFound             = errs.NewResourceNotFoundError("the customer with given id was not found")
	ErrOrderTotalExceedsCreditLimit = errs.NewBusinessLogicError("order total exceeds credit limit")
	ErrOrderNotFound                = errs.NewResourceNotFoundError("the order with given id was not found")
	ErrReleaseCreditFailed          = errs.NewBusinessLogicError("release credit failed")
)

func (s *orderService) CreateOrder(ctx context.Context, req *dtos.CreateOrderRequest) (int, error) {
	// validate request
	if err := req.Validate(); err != nil {
		return 0, errs.NewValidationError(err.Error())
	}

	var orderId int
	err := s.transactor.WithinTransaction(ctx, func(ctx context.Context) error {
		customer, err := s.custRepo.FindByID(ctx, req.CustomerID)
		if err != nil {
			log.Println(err)
			return errs.NewDatabaseFailureError(err.Error())
		}

		if customer == nil {
			return ErrCustomerNotFound
		}

		if err := customer.ReserveCredit(req.OrderTotal); err != nil {
			log.Println(err)
			return ErrOrderTotalExceedsCreditLimit
		}

		if err := s.custRepo.UpdateCreditLimit(ctx, customer); err != nil {
			log.Println(err)
			return errs.NewDatabaseFailureError(err.Error())
		}

		order := model.NewOrder(req.CustomerID, req.OrderTotal)
		err = s.orderRepo.Create(ctx, order)
		if err != nil {
			log.Println(err)
			return errs.NewDatabaseFailureError(err.Error())
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

	customer, err := s.custRepo.FindByID(ctx, order.CustomerID)
	if err != nil {
		log.Println(err)
		return ErrCustomerNotFound
	}
	if err := customer.ReleaseCredit(order.OrderTotal); err != nil {
		log.Println(err)
		return ErrReleaseCreditFailed
	}

	if err := s.custRepo.UpdateCreditLimit(ctx, customer); err != nil {
		log.Println(err)
		return errs.NewDatabaseFailureError(err.Error())
	}

	return nil
}
