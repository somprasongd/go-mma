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
	ReserveCredit(ctx context.Context, id int, amount int) error
	ReleaseCredit(ctx context.Context, id int, amount int) error
	GetCustomerByID(ctx context.Context, id int) (*dtos.Customer, error)
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
	ErrEmailExists                  = errs.NewDuplicateEntryError("email already exists")
	ErrCustomerNotFound             = errs.NewResourceNotFoundError("the customer with given id was not found")
	ErrOrderTotalExceedsCreditLimit = errs.NewBusinessLogicError("order total exceeds credit limit")
	ErrReleaseCreditFailed          = errs.NewBusinessLogicError("release credit failed")
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

func (s *customerService) GetCustomerByID(ctx context.Context, id int) (*dtos.Customer, error) {
	customer, err := s.custRepo.FindByID(ctx, id)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	if customer == nil {
		return nil, ErrCustomerNotFound
	}

	return &dtos.Customer{
		ID:          customer.ID,
		Email:       customer.Email,
		CreditLimit: customer.CreditLimit,
	}, nil
}

func (s *customerService) ReserveCredit(ctx context.Context, id int, amount int) error {
	customer, err := s.custRepo.FindByID(ctx, id)
	if err != nil {
		log.Println(err)
		return err
	}

	if customer == nil {
		return ErrCustomerNotFound
	}

	if err := customer.ReserveCredit(amount); err != nil {
		log.Println(err)
		return ErrOrderTotalExceedsCreditLimit
	}

	if err := s.custRepo.UpdateCreditLimit(ctx, customer); err != nil {
		log.Println(err)
		return errs.NewDatabaseFailureError(err.Error())
	}

	return nil
}

func (s *customerService) ReleaseCredit(ctx context.Context, id int, amount int) error {
	customer, err := s.custRepo.FindByID(ctx, id)
	if err != nil {
		log.Println(err)
		return err
	}

	if customer == nil {
		return ErrCustomerNotFound
	}

	if err := customer.ReleaseCredit(amount); err != nil {
		log.Println(err)
		return ErrReleaseCreditFailed
	}

	if err := s.custRepo.UpdateCreditLimit(ctx, customer); err != nil {
		log.Println(err)
		return errs.NewDatabaseFailureError(err.Error())
	}

	return nil
}
