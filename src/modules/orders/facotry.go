package orders

import (
	"go-mma/shared/common/module"
	"go-mma/shared/common/registry"

	"github.com/gin-gonic/gin"

	"go-mma/modules/orders/handler"
	"go-mma/modules/orders/repository"
	"go-mma/modules/orders/service"

	customerContracts "go-mma/shared/contracts/customer_contracts"
	notificationContracts "go-mma/shared/contracts/notification_contracts"
	orderContracts "go-mma/shared/contracts/order_contracts"
)

type mod struct {
	mCtx     *module.ModuleContext
	orderSvc orderContracts.OrderService
}

func NewModule(mCtx *module.ModuleContext) module.Module {
	return &mod{mCtx: mCtx}
}

func (m *mod) Init(reg registry.ServiceRegistry) error {
	// Resolve CustomerService from the registry
	custSvc, err := registry.ResolveAs[customerContracts.CustomerService](reg, customerContracts.CustomerServiceKey)
	if err != nil {
		return err
	}

	// Resolve NotificationService from the registry
	notiSvc, err := registry.ResolveAs[notificationContracts.NotificationService](reg, notificationContracts.NotificationServiceKey)
	if err != nil {
		return err
	}

	repo := repository.NewOrderRepository(m.mCtx.DBCtx)
	m.orderSvc = service.NewOrderService(m.mCtx.Transactor, custSvc, repo, notiSvc)

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
