package repository

import (
	"context"
	"fmt"
	"go-mma/data/sqldb"
	"go-mma/models"
	"time"
)

type CustomerRepository struct {
	dbCtx sqldb.DBContext
}

func NewCustomerRepository(dbCtx sqldb.DBContext) *CustomerRepository {
	return &CustomerRepository{
		dbCtx: dbCtx,
	}
}

func (r *CustomerRepository) Create(ctx context.Context, customer *models.Customer) error {
	query := `
INSERT INTO public.customers (
    name, credit_limit
)
VALUES ($1, $2)
RETURNING *
`

	ctx, cancel := context.WithTimeout(ctx, 20*time.Second)
	defer cancel()
	err := r.dbCtx.DB().QueryRowxContext(ctx, query, customer.Name, customer.CreditLimit).StructScan(customer)
	if err != nil {
		return fmt.Errorf("failed to create customer: %w", err)
	}
	return nil
}
