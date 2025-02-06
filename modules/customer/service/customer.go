package service

import (
	"context"
	"go-mma/modules/customer/dtos"
	"go-mma/modules/customer/model"
	"go-mma/modules/customer/repository"
	"go-mma/util/errs"
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

var (
	ErrEmailExists = errs.NewDuplicateEntryError("email already exists")
)

func (s *customerService) CreateCustomer(ctx context.Context, req *dtos.CreateCustomerRequest) (int, error) {
	// validate the request
	if err := req.Validate(); err != nil {
		return 0, errs.NewValidationError(err.Error())
	}

	// create model
	customer := model.NewCustomer(req.Email, req.CreditLimit)

	// save to database
	if err := s.custRepo.Create(ctx, customer); err != nil {
		log.Println(err)
		if errs.IsErrDuplicateEntry(err) {
			return 0, ErrEmailExists
		}
		return 0, err
	}
	return customer.ID, nil
}
