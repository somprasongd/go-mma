package ddd

import "time"

type Entity[TId any] struct {
	ID        TId
	CreatedAt time.Time
	CreatedBy string
	UpdatedAt time.Time
	UpdatedBy string
}
