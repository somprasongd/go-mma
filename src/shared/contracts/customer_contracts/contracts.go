package customers_contracts

import (
	"context"
	"go-mma/shared/common/registry"
)

const (
	CustomerServiceKey registry.ServiceKey = "CustomerService"
)

type CustomerService interface {
	CreateCustomer(ctx context.Context, req *CreateCustomerRequest) (int, error)
	ReserveCredit(ctx context.Context, id int, amount int) error
	ReleaseCredit(ctx context.Context, id int, amount int) error
	GetCustomerByID(ctx context.Context, id int) (*CustomerInfo, error)
}
