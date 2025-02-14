package customers_contracts

import (
	"context"
	"go-mma/shared/common/registry"
)

const (
	CustomerServiceKey registry.ServiceKey = "CustomerService"
)

type CustomerCommon interface {
	GetCustomerByID(ctx context.Context, id int) (*CustomerInfo, error)
}

type CustomerFactory interface {
	CreateCustomer(ctx context.Context, req *CreateCustomerRequest) (int, error)
	CustomerCommon
}

type CreditManagement interface {
	ReserveCredit(ctx context.Context, id int, amount int) error
	ReleaseCredit(ctx context.Context, id int, amount int) error
	CustomerCommon
}
