package model

import (
	"fmt"
	"go-mma/shared/common/ddd"
	"time"
)

type Customer struct {
	ID          int       `db:"id"`
	Email       string    `db:"email"`
	CreditLimit int       `db:"credit_limit"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
	ddd.Aggregate[int]
}

func NewCustomer(email string, creditLimit int) *Customer {
	return &Customer{
		Email:       email,
		CreditLimit: creditLimit,
	}
}

func (c *Customer) ReserveCredit(v int) error {
	newCreditLimit := c.CreditLimit - v
	if newCreditLimit < 0 {
		return fmt.Errorf("insufficient credit limit")
	}
	c.CreditLimit = newCreditLimit
	return nil
}

func (c *Customer) ReleaseCredit(v int) error {
	if c.CreditLimit <= 0 {
		c.CreditLimit = 0
	}
	c.CreditLimit = c.CreditLimit + v
	return nil
}
