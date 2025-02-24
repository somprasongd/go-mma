package dtos

type CustomerInfo struct {
	ID          int    `json:"id"`
	Email       string `json:"email"`
	CreditLimit int    `json:"credit_limit"`
}
