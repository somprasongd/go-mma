package customer

import (
	"go-mma/modules"
	"go-mma/modules/customer/handler"
	"go-mma/modules/customer/repository"
	"go-mma/modules/customer/service"

	"github.com/gin-gonic/gin"
)

func NewModule(mCtx *modules.ModuleContext) modules.Module {
	return &module{mCtx}
}

type module struct {
	mCtx *modules.ModuleContext
}

func (m *module) RegisterRoutes(router *gin.Engine) {
	repo := repository.NewCustomerRepository(m.mCtx.DBCtx)
	serv := service.NewCustomerService(repo)
	hdl := handler.NewCustomerHandler(serv)

	rCustomer := router.Group("/api/v1/customers")
	{
		rCustomer.POST("", hdl.CreateCustomer)
	}
}
