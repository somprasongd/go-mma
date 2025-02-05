package repository

import (
	"context"
	"database/sql"
	"fmt"
	"go-mma/models"
	"go-mma/util/transactor"
	"time"
)

type OrderRepository struct {
	dbCtx transactor.DBContext
}

func NewOrderRepository(dbCtx transactor.DBContext) *OrderRepository {
	return &OrderRepository{
		dbCtx: dbCtx,
	}
}

func (r *OrderRepository) Create(ctx context.Context, m *models.Order) error {
	query := `
	INSERT INTO public.orders (
			customer_id, order_total
	)
	VALUES ($1, $2)
	RETURNING *
	`

	ctx, cancel := context.WithTimeout(ctx, 20*time.Second)
	defer cancel()

	err := r.dbCtx(ctx).QueryRowxContext(ctx, query, m.CustomerID, m.OrderTotal).StructScan(m)
	if err != nil {
		return fmt.Errorf("failed to create order: %w", err)
	}
	return nil
}

func (r *OrderRepository) FindByID(ctx context.Context, id int) (*models.Order, error) {
	query := `
	SELECT *
	FROM public.orders
	WHERE id = $1
	AND canceled_at IS NULL
`
	ctx, cancel := context.WithTimeout(ctx, 20*time.Second)
	defer cancel()

	var order models.Order
	err := r.dbCtx(ctx).QueryRowxContext(ctx, query, id).StructScan(&order)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get order by ID: %w", err)
	}
	return &order, nil
}

func (r *OrderRepository) Cancel(ctx context.Context, id int) error {
	query := `
	UPDATE public.orders
	SET canceled_at = current_timestamp
	WHERE id = $1
`
	ctx, cancel := context.WithTimeout(ctx, 20*time.Second)
	defer cancel()
	_, err := r.dbCtx(ctx).ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to cancel order: %w", err)
	}
	return nil
}
