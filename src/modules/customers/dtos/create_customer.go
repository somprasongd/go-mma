package dtos

type CreateCustomerRequest struct {
	Email       string `json:"email"`
	CreditLimit int    `json:"credit_limit"`
}

type CreateCustomerResponse struct {
	ID int `json:"id"`
}
