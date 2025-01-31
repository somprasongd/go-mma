package repository

import (
	"context"
	"database/sql"
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
			email, credit_limit
	)
	VALUES ($1, $2)
	RETURNING *
	`

	ctx, cancel := context.WithTimeout(ctx, 20*time.Second)
	defer cancel()
	err := r.dbCtx.DB().QueryRowxContext(ctx, query, customer.Email, customer.CreditLimit).StructScan(customer)
	if err != nil {
		return fmt.Errorf("failed to create customer: %w", err)
	}
	return nil
}

func (r *CustomerRepository) FindByID(ctx context.Context, id int) (*models.Customer, error) {
	query := `
	SELECT *
	FROM public.customers
	WHERE id = $1
`
	ctx, cancel := context.WithTimeout(ctx, 20*time.Second)
	defer cancel()

	var customer models.Customer
	err := r.dbCtx.DB().QueryRowxContext(ctx, query, id).StructScan(&customer)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get customer by ID: %w", err)
	}

	return &customer, nil
}

func (r *CustomerRepository) UpdateCreditLimit(ctx context.Context, m *models.Customer) error {
	query := `
	UPDATE public.customers
	SET credit_limit = $2
	WHERE id = $1
	RETURNING *
`
	ctx, cancel := context.WithTimeout(ctx, 20*time.Second)
	defer cancel()

	err := r.dbCtx.DB().QueryRowxContext(ctx, query, m.ID, m.CreditLimit).StructScan(m)
	if err != nil {
		return fmt.Errorf("failed to update customer credit limit: %w", err)
	}
	return nil
}
