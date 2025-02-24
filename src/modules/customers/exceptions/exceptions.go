package exceptions

import "go-mma/shared/common/errs"

var (
	ErrEmailExists                  = errs.NewDuplicateEntryError("email already exists")
	ErrCustomerNotFound             = errs.NewResourceNotFoundError("the customer with given id was not found")
	ErrOrderTotalExceedsCreditLimit = errs.NewBusinessLogicError("order total exceeds credit limit")
	ErrReleaseCreditFailed          = errs.NewBusinessLogicError("release credit failed")
)
