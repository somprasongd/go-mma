package dtos

import "fmt"

type CreateOrderRequest struct {
	CustomerID int `json:"customer_id"`
	OrderTotal int `json:"order_total"`
}

func (r *CreateOrderRequest) Validate() error {
	if r.CustomerID <= 0 {
		return fmt.Errorf("customer ID must be greater than 0")
	}
	if r.OrderTotal <= 0 {
		return fmt.Errorf("order total must be greater than 0")
	}
	return nil
}
