package services

import (
	"context"
	"go-mma/dtos"
	"go-mma/models"
	"go-mma/repository"
	"log"
)

type CustomerService interface {
	CreateCustomer(ctx context.Context, req *dtos.CreateCustomerRequest) (int, error)
}

type customerService struct {
	custRepo repository.CustomerRepository
}

func NewCustomerService(custRepo repository.CustomerRepository) CustomerService {
	return &customerService{
		custRepo: custRepo,
	}
}

func (s *customerService) CreateCustomer(ctx context.Context, req *dtos.CreateCustomerRequest) (int, error) {
	// validate the request
	if err := req.Validate(); err != nil {
		return 0, err
	}

	// create model
	customer := models.NewCustomer(req.Email, req.CreditLimit)

	// save to database
	if err := s.custRepo.Create(ctx, customer); err != nil {
		log.Println(err)
		return 0, err
	}
	return customer.ID, nil
}
