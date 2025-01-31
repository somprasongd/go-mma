package models

import "time"

type Order struct {
	ID         int
	CustomerID int
	OrderTotal int
	CreatedAt  time.Time
	CanceledAt *time.Time
}
