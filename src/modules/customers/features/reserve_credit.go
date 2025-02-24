package features

import (
	"context"
	"go-mma/modules/customers/exceptions"
	"go-mma/modules/customers/repository"
	"go-mma/shared/common/errs"
	"go-mma/shared/common/mediator"
	"go-mma/shared/contracts/customer_contracts/commands"
	"log"
)

type reserveCreditHandler struct {
	custRepo repository.CustomerRepository
}

func NewReserveCreditHandler(repo repository.CustomerRepository) *reserveCreditHandler {
	return &reserveCreditHandler{
		custRepo: repo,
	}
}

func (h *reserveCreditHandler) Handle(ctx context.Context, cmd *commands.ReserveCreditCommand) (*mediator.NoResponse, error) {
	customer, err := h.custRepo.FindByID(ctx, cmd.ID)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	if customer == nil {
		return nil, exceptions.ErrCustomerNotFound
	}

	if err := customer.ReserveCredit(cmd.CreditAmount); err != nil {
		log.Println(err)
		return nil, exceptions.ErrOrderTotalExceedsCreditLimit
	}

	if err := h.custRepo.UpdateCreditLimit(ctx, customer); err != nil {
		log.Println(err)
		return nil, errs.NewDatabaseFailureError(err.Error())
	}

	return nil, nil
}
