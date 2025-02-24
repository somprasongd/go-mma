package features

import (
	"context"
	"errors"
	"go-mma/modules/customers/dtos"
	"go-mma/modules/customers/exceptions"
	"go-mma/modules/customers/model"
	"go-mma/modules/customers/repository"
	"go-mma/shared/common/errs"
	"log"
)

type CreateCustomerCommand struct {
	*dtos.CreateCustomerRequest
}

func (r *CreateCustomerCommand) Validate() error {
	if r.Email == "" {
		return errors.New("Email cannot be empty")
	}
	if r.CreditLimit <= 0 {
		return errors.New("Credit limit must be greater than 0")
	}
	return nil
}

type createCustomerHandler struct {
	custRepo repository.CustomerRepository
}

func NewCreateCustomerHandler(custRepo repository.CustomerRepository) *createCustomerHandler {
	return &createCustomerHandler{
		custRepo: custRepo,
	}
}

func (h *createCustomerHandler) Handle(ctx context.Context, cmd *CreateCustomerCommand) (*dtos.CreateCustomerResponse, error) {
	// validate the request
	if err := cmd.Validate(); err != nil {
		return nil, errs.NewValidationError(err.Error())
	}

	// create model
	customer := model.NewCustomer(cmd.Email, cmd.CreditLimit)

	// save to database
	if err := h.custRepo.Create(ctx, customer); err != nil {
		log.Println(err)
		if errs.IsErrDuplicateEntry(err) {
			return nil, exceptions.ErrEmailExists
		}
		return nil, err
	}
	return newCreateCustomerResponse(customer.ID), nil
}

func newCreateCustomerResponse(id int) *dtos.CreateCustomerResponse {
	return &dtos.CreateCustomerResponse{ID: id}
}
