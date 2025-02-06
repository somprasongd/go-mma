package repository

import (
	"context"
	"database/sql"
	"fmt"
	"go-mma/models"
	"go-mma/util/errs"
	"go-mma/util/transactor"
	"time"
)

const queryTimeout = 20 * time.Second

type CustomerRepository interface {
	Create(ctx context.Context, customer *models.Customer) error
	FindByID(ctx context.Context, id int) (*models.Customer, error)
	UpdateCreditLimit(ctx context.Context, customer *models.Customer) error
}

type customerRepository struct {
	dbCtx transactor.DBContext
}

// NewCustomerRepository returns the CustomerRepository interface implementation.
func NewCustomerRepository(dbCtx transactor.DBContext) CustomerRepository {
	return &customerRepository{dbCtx: dbCtx}
}

func (r *customerRepository) Create(ctx context.Context, customer *models.Customer) error {
	query := `
		INSERT INTO public.customers (email, credit_limit)
		VALUES ($1, $2)
		RETURNING id, email, credit_limit
	`

	ctx, cancel := context.WithTimeout(ctx, queryTimeout)
	defer cancel()

	err := r.dbCtx(ctx).QueryRowxContext(ctx, query, customer.Email, customer.CreditLimit).
		Scan(&customer.ID, &customer.Email, &customer.CreditLimit)
	if err != nil {
		return errs.HandleDBError(err)
	}

	return nil
}

func (r *customerRepository) FindByID(ctx context.Context, id int) (*models.Customer, error) {
	query := `
		SELECT id, email, credit_limit 
		FROM public.customers 
		WHERE id = $1
	`

	ctx, cancel := context.WithTimeout(ctx, queryTimeout)
	defer cancel()

	var customer models.Customer
	err := r.dbCtx(ctx).GetContext(ctx, &customer, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, errs.HandleDBError(fmt.Errorf("failed to get customer by ID: %w", err))
	}

	return &customer, nil
}

func (r *customerRepository) UpdateCreditLimit(ctx context.Context, customer *models.Customer) error {
	query := `
		UPDATE public.customers
		SET credit_limit = $2
		WHERE id = $1
		RETURNING credit_limit
	`

	ctx, cancel := context.WithTimeout(ctx, queryTimeout)
	defer cancel()

	err := r.dbCtx(ctx).QueryRowxContext(ctx, query, customer.ID, customer.CreditLimit).
		Scan(&customer.CreditLimit)
	if err != nil {
		return errs.HandleDBError(fmt.Errorf("failed to update customer credit limit: %w", err))
	}

	return nil
}
