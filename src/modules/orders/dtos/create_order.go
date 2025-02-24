package dtos

type CreateOrderRequest struct {
	CustomerID int `json:"customer_id"`
	OrderTotal int `json:"order_total"`
}

type CreateOrderResponse struct {
	ID int `json:"id"`
}
