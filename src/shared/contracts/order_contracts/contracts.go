package ordercontracts

import (
	"context"
	"go-mma/shared/common/registry"
)

const (
	OrderServiceKey registry.ServiceKey = "OrderService"
)

type OrderService interface {
	CreateOrder(ctx context.Context, req *CreateOrderRequest) (int, error)
	CancelOrder(ctx context.Context, id int) error
}
