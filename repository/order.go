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

const orderQueryTimeout = 20 * time.Second

type OrderRepository interface {
	Create(ctx context.Context, order *models.Order) error
	FindByID(ctx context.Context, id int) (*models.Order, error)
	Cancel(ctx context.Context, id int) error
}

type orderRepository struct {
	dbCtx transactor.DBContext
}

// NewOrderRepository returns an OrderRepository interface implementation.
func NewOrderRepository(dbCtx transactor.DBContext) OrderRepository {
	return &orderRepository{dbCtx: dbCtx}
}

func (r *orderRepository) Create(ctx context.Context, order *models.Order) error {
	query := `
		INSERT INTO public.orders (customer_id, order_total)
		VALUES ($1, $2)
		RETURNING id, customer_id, order_total, created_at
	`

	ctx, cancel := context.WithTimeout(ctx, orderQueryTimeout)
	defer cancel()

	err := r.dbCtx(ctx).QueryRowxContext(ctx, query, order.CustomerID, order.OrderTotal).
		Scan(&order.ID, &order.CustomerID, &order.OrderTotal, &order.CreatedAt)
	if err != nil {
		return errs.HandleDBError(fmt.Errorf("failed to create order: %w", err))
	}
	return nil
}

func (r *orderRepository) FindByID(ctx context.Context, id int) (*models.Order, error) {
	query := `
		SELECT id, customer_id, order_total, created_at, canceled_at
		FROM public.orders
		WHERE id = $1 AND canceled_at IS NULL
	`

	ctx, cancel := context.WithTimeout(ctx, orderQueryTimeout)
	defer cancel()

	var order models.Order
	err := r.dbCtx(ctx).GetContext(ctx, &order, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, errs.HandleDBError(fmt.Errorf("failed to get order by ID: %w", err))
	}
	return &order, nil
}

func (r *orderRepository) Cancel(ctx context.Context, id int) error {
	query := `
		UPDATE public.orders
		SET canceled_at = current_timestamp
		WHERE id = $1
		RETURNING canceled_at
	`

	ctx, cancel := context.WithTimeout(ctx, orderQueryTimeout)
	defer cancel()

	var canceledAt sql.NullTime
	err := r.dbCtx(ctx).QueryRowxContext(ctx, query, id).Scan(&canceledAt)
	if err != nil {
		return errs.HandleDBError(fmt.Errorf("failed to cancel order: %w", err))
	}

	return nil
}
