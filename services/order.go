package services

import (
	"context"
	"fmt"
	"go-mma/dtos"
	"go-mma/models"
	"go-mma/repository"
	"go-mma/util/transactor"
	"log"
)

type OrderService struct {
	transactor  transactor.Transactor
	custRepo    *repository.CustomerRepository
	orderRepo   *repository.OrderRepository
	notiService *NotificationService
}

func NewOrderService(transactor transactor.Transactor, custRepo *repository.CustomerRepository, orderRepo *repository.OrderRepository, notiService *NotificationService) *OrderService {
	return &OrderService{
		transactor:  transactor,
		custRepo:    custRepo,
		orderRepo:   orderRepo,
		notiService: notiService,
	}
}

func (s *OrderService) CreateOrder(ctx context.Context, req *dtos.CreateOrderRequest) (int, error) {
	// validate request
	if err := req.Validate(); err != nil {
		return 0, err
	}

	var orderId int
	err := s.transactor.WithinTransaction(ctx, func(ctx context.Context) error {
		customer, err := s.custRepo.FindByID(ctx, req.CustomerID)
		if err != nil {
			log.Println(err)
			return fmt.Errorf("failed to get customer: %w", err)
		}

		if customer == nil {
			return fmt.Errorf("customer not found")
		}

		if err := customer.ReserveCredit(req.OrderTotal); err != nil {
			log.Println(err)
			return err
		}

		if err := s.custRepo.UpdateCreditLimit(ctx, customer); err != nil {
			log.Println(err)
			return err
		}

		order := models.NewOrder(req.CustomerID, req.OrderTotal)
		err = s.orderRepo.Create(ctx, order)
		if err != nil {
			log.Println(err)
			return fmt.Errorf("failed to create order: %w", err)
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

func (s *OrderService) CancelOrder(ctx context.Context, id int) error {
	order, err := s.orderRepo.FindByID(ctx, id)
	if err != nil {
		log.Println(err)
		return fmt.Errorf("failed to get order: %w", err)
	}

	if order == nil {
		return fmt.Errorf("order not found")
	}

	if err := s.orderRepo.Cancel(ctx, order.ID); err != nil {
		log.Println(err)
		return fmt.Errorf("failed to cancel order: %w", err)
	}

	customer, err := s.custRepo.FindByID(ctx, order.CustomerID)
	if err != nil {
		log.Println(err)
		return fmt.Errorf("failed to get customer: %w", err)
	}
	if err := customer.ReleaseCredit(order.OrderTotal); err != nil {
		log.Println(err)
		return err
	}

	if err := s.custRepo.UpdateCreditLimit(ctx, customer); err != nil {
		log.Println(err)
		return err
	}

	return nil
}
