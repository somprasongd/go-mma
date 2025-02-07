package orders

import (
	"go-mma/shared/common/module"

	"github.com/gin-gonic/gin"

	custRepo "go-mma/modules/customers/repository"
	custServ "go-mma/modules/customers/service"
	notiServ "go-mma/modules/notifications/service"
	"go-mma/modules/orders/handler"
	"go-mma/modules/orders/repository"
	"go-mma/modules/orders/service"
)

func NewModule(mCtx *module.ModuleContext) module.Module {
	return &mod{mCtx}
}

type mod struct {
	mCtx *module.ModuleContext
}

func (m *mod) RegisterRoutes(router *gin.Engine) {
	servNoti := notiServ.NewNotificationService()

	custRepo := custRepo.NewCustomerRepository(m.mCtx.DBCtx)
	servCust := custServ.NewCustomerService(custRepo)

	repoOrder := repository.NewOrderRepository(m.mCtx.DBCtx)
	serv := service.NewOrderService(m.mCtx.Transactor, servCust, repoOrder, servNoti)
	hdl := handler.NewOrderHandler(serv)

	rOrder := router.Group("/api/v1/orders")
	{
		rOrder.POST("", hdl.CreateOrder)
		rOrder.DELETE("/:id", hdl.CancelOrder)
	}
}
