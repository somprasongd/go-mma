package features

import (
	"context"
	"go-mma/modules/customers/exceptions"
	"go-mma/modules/customers/model"
	"go-mma/modules/customers/repository"
	"go-mma/shared/contracts/customer_contracts/dtos"
	"go-mma/shared/contracts/customer_contracts/queries"
)

type getCustomerByIDQuery struct {
	custRepo repository.CustomerRepository
}

func NewGetCustomerByIDQuery(custRepo repository.CustomerRepository) *getCustomerByIDQuery {
	return &getCustomerByIDQuery{
		custRepo: custRepo,
	}
}

func (h *getCustomerByIDQuery) Handle(ctx context.Context, query *queries.GetCustomerByIDQuery) (*queries.GetCustomerByIDResult, error) {
	customer, err := h.custRepo.FindByID(ctx, query.ID)
	if err != nil {
		return nil, err
	}
	if customer == nil {
		return nil, exceptions.ErrCustomerNotFound
	}
	return newGetCustomerByIDResult(customer), nil
}

func newGetCustomerByIDResult(customer *model.Customer) *queries.GetCustomerByIDResult {
	return &queries.GetCustomerByIDResult{
		CustomerInfo: &dtos.CustomerInfo{
			ID:          customer.ID,
			Email:       customer.Email,
			CreditLimit: customer.CreditLimit,
		},
	}
}
