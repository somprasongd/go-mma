package order

import (
	"go-mma/modules"

	"github.com/gin-gonic/gin"

	custRepo "go-mma/modules/customer/repository"
	custServ "go-mma/modules/customer/service"
	notiServ "go-mma/modules/notification/service"
	"go-mma/modules/order/handler"
	"go-mma/modules/order/repository"
	"go-mma/modules/order/service"
)

func NewModule(mCtx *modules.ModuleContext) modules.Module {
	return &module{mCtx}
}

type module struct {
	mCtx *modules.ModuleContext
}

func (m *module) RegisterRoutes(router *gin.Engine) {
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
