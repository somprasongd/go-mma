package queries

import "go-mma/shared/contracts/customer_contracts/dtos"

type GetCustomerByIDQuery struct {
	ID int `json:"id"`
}

type GetCustomerByIDResult struct {
	*dtos.CustomerInfo
}
