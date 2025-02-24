package customers

import (
	"go-mma/modules/customers/features"
	"go-mma/modules/customers/handler"
	"go-mma/modules/customers/repository"
	"go-mma/shared/common/eventbus"
	"go-mma/shared/common/mediator"
	"go-mma/shared/common/module"
	"go-mma/shared/common/registry"

	"github.com/gin-gonic/gin"
)

type mod struct {
	mCtx *module.ModuleContext
}

func NewModule(mCtx *module.ModuleContext) module.Module {
	return &mod{mCtx: mCtx}
}

func (m *mod) Init(reg registry.ServiceRegistry, eventbus eventbus.EventBus) error {
	repo := repository.NewCustomerRepository(m.mCtx.DBCtx)

	mediator.Register(features.NewCreateCustomerHandler(repo))
	mediator.Register(features.NewGetCustomerByIDQuery(repo))
	mediator.Register(features.NewReserveCreditHandler(repo))
	mediator.Register(features.NewReleaseCreditHandler(repo))

	return nil
}

func (m *mod) RegisterRoutes(router *gin.Engine) {
	hdl := handler.NewCustomerHandler()

	rCustomer := router.Group("/api/v1/customers")
	{
		rCustomer.POST("", hdl.CreateCustomer)
	}
}
