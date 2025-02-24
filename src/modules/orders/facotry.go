package orders

import (
	"go-mma/shared/common/eventbus"
	"go-mma/shared/common/mediator"
	"go-mma/shared/common/module"
	"go-mma/shared/common/registry"

	"github.com/gin-gonic/gin"

	"go-mma/modules/orders/features"
	"go-mma/modules/orders/handler"
	"go-mma/modules/orders/repository"
)

type mod struct {
	mCtx *module.ModuleContext
}

func NewModule(mCtx *module.ModuleContext) module.Module {
	return &mod{mCtx: mCtx}
}

func (m *mod) Init(reg registry.ServiceRegistry, eventbus eventbus.EventBus) error {
	repo := repository.NewOrderRepository(m.mCtx.DBCtx)

	mediator.Register(features.NewCreateOrderCommand(m.mCtx.Transactor, repo, eventbus))
	mediator.Register(features.NewCancelOrderCommandHandler(m.mCtx.Transactor, repo))

	return nil
}

func (m *mod) RegisterRoutes(router *gin.Engine) {
	hdl := handler.NewOrderHandler()

	rOrder := router.Group("/api/v1/orders")
	{
		rOrder.POST("", hdl.CreateOrder)
		rOrder.DELETE("/:id", hdl.CancelOrder)
	}
}
