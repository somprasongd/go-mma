package models

import "time"

type Customer struct {
	ID          int       `db:"id"`
	Name        string    `db:"name"`
	CreditLimit int       `db:"credit_limit"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}

func NewCustomer(name string, creditLimit int) *Customer {
	return &Customer{
		Name:        name,
		CreditLimit: creditLimit,
	}
}
