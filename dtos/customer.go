package dtos

import "errors"

type CreateCustomerRequest struct {
	Name        string `json:"name"`
	CreditLimit int    `json:"credit_limit"`
}

func (r *CreateCustomerRequest) Validate() error {
	if r.Name == "" {
		return errors.New("Name cannot be empty")
	}
	if r.CreditLimit <= 0 {
		return errors.New("Credit limit must be greater than 0")
	}
	return nil
}
