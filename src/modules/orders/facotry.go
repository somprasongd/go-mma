package orders

import (
	"go-mma/shared/common/module"
	"go-mma/shared/common/registry"

	"github.com/gin-gonic/gin"

	custModule "go-mma/modules/customers"
	custService "go-mma/modules/customers/service"
	notiModule "go-mma/modules/notifications"
	notiService "go-mma/modules/notifications/service"
	"go-mma/modules/orders/handler"
	"go-mma/modules/orders/repository"
	"go-mma/modules/orders/service"
)

const (
	OrderServiceKey registry.ServiceKey = "OrderService"
)

type mod struct {
	mCtx     *module.ModuleContext
	orderSvc service.OrderService
}

func NewModule(mCtx *module.ModuleContext) module.Module {
	return &mod{mCtx: mCtx}
}

func (m *mod) Init(reg registry.ServiceRegistry) error {
	// Resolve CustomerService from the registry
	custSvc, err := registry.ResolveAs[custService.CustomerService](reg, custModule.CustomerServiceKey)
	if err != nil {
		return err
	}

	// Resolve NotificationService from the registry
	notiSvc, err := registry.ResolveAs[notiService.NotificationService](reg, notiModule.NotificationServiceKey)
	if err != nil {
		return err
	}

	repo := repository.NewOrderRepository(m.mCtx.DBCtx)
	m.orderSvc = service.NewOrderService(m.mCtx.Transactor, custSvc, repo, notiSvc)

	reg.Register(OrderServiceKey, m.orderSvc)
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
