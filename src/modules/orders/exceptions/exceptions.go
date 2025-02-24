package exceptions

import "go-mma/shared/common/errs"

var (
	ErrOrderNotFound = errs.NewResourceNotFoundError("the order with given id was not found")
)
