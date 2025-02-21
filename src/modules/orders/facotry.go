package orders

import (
	"go-mma/shared/common/ddd"
	"go-mma/shared/common/eventbus"
	"go-mma/shared/common/module"
	"go-mma/shared/common/registry"

	"github.com/gin-gonic/gin"

	eventhandlers "go-mma/modules/orders/event-handlers"
	"go-mma/modules/orders/events"
	"go-mma/modules/orders/handler"
	"go-mma/modules/orders/repository"
	"go-mma/modules/orders/service"

	customerContracts "go-mma/shared/contracts/customer_contracts"
	orderContracts "go-mma/shared/contracts/order_contracts"
)

type mod struct {
	mCtx     *module.ModuleContext
	orderSvc orderContracts.OrderService
}

func NewModule(mCtx *module.ModuleContext) module.Module {
	return &mod{mCtx: mCtx}
}

func (m *mod) Init(reg registry.ServiceRegistry, eventbus eventbus.EventBus) error {
	// Resolve CustomerService from the registry
	custSvc, err := registry.ResolveAs[customerContracts.CreditManagement](reg, customerContracts.CustomerServiceKey)
	if err != nil {
		return err
	}

	// Register domain event handler
	dispatcher := ddd.NewEventDispatcher()
	dispatcher.Register(&events.OrderPlacedEvent{}, eventhandlers.NewOrderPlacedEventHandler(custSvc))

	repo := repository.NewOrderRepository(m.mCtx.DBCtx)
	m.orderSvc = service.NewOrderService(m.mCtx.Transactor, custSvc, repo, eventbus)

	reg.Register(orderContracts.OrderServiceKey, m.orderSvc)

	return nil
}

func (m *mod) RegisterRoutes(router *gin.Engine) {
	hdl := handler.NewOrderHandler(m.orderSvc)

	rOrder := router.Group("/api/v1/orders")
	{
		rOrder.POST("", hdl.CreateOrder)
		rOrder.DELETE("/:id", hdl.CancelOrder)
	}
}
