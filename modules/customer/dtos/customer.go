package dtos

import "errors"

type CreateCustomerRequest struct {
	Email       string `json:"email"`
	CreditLimit int    `json:"credit_limit"`
}

func (r *CreateCustomerRequest) Validate() error {
	if r.Email == "" {
		return errors.New("Email cannot be empty")
	}
	if r.CreditLimit <= 0 {
		return errors.New("Credit limit must be greater than 0")
	}
	return nil
}
