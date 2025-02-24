package features

import (
	"context"
	"go-mma/modules/orders/exceptions"
	"go-mma/modules/orders/model"
	"go-mma/modules/orders/repository"
	"go-mma/shared/common/errs"
	"go-mma/shared/common/mediator"
	"go-mma/shared/common/storage/db/transactor"
	"log"

	custCmd "go-mma/shared/contracts/customer_contracts/commands"
)

type CancelOrderCommand struct {
	ID int `json:"id"`
}

type cancelOrderCommandHanler struct {
	transactor transactor.Transactor
	orderRepo  repository.OrderRepository
}

func NewCancelOrderCommandHandler(transactor transactor.Transactor, orderRepo repository.OrderRepository) *cancelOrderCommandHanler {
	return &cancelOrderCommandHanler{
		transactor: transactor,
		orderRepo:  orderRepo,
	}
}

func (h *cancelOrderCommandHanler) Handle(ctx context.Context, cmd *CancelOrderCommand) (*mediator.NoResponse, error) {
	err := h.transactor.WithinTransaction(ctx, func(ctx context.Context) error {
		order, err := h.orderRepo.FindByID(ctx, cmd.ID)
		if err != nil {
			log.Println(err)
			return errs.NewDatabaseFailureError(err.Error())
		}

		if order == nil {
			return exceptions.ErrOrderNotFound
		}

		if err := h.orderRepo.Cancel(ctx, order.ID); err != nil {
			log.Println(err)
			return errs.NewDatabaseFailureError(err.Error())
		}

		_, err = mediator.Send[*custCmd.ReleaseCreditCommand, *mediator.NoResponse](
			ctx,
			newReleaseCreditCommand(order),
		)

		if err != nil {
			return err
		}

		return nil
	})

	return &mediator.NoResponse{}, err

}

func newReleaseCreditCommand(order *model.Order) *custCmd.ReleaseCreditCommand {
	return &custCmd.ReleaseCreditCommand{
		ID:           order.CustomerID,
		CreditAmount: order.OrderTotal,
	}
}
