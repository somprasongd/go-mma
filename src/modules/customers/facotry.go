package customers

import (
	"go-mma/modules/customers/handler"
	"go-mma/modules/customers/repository"
	"go-mma/modules/customers/service"
	"go-mma/shared/common/module"

	"github.com/gin-gonic/gin"
)

func NewModule(mCtx *module.ModuleContext) module.Module {
	return &mod{mCtx}
}

type mod struct {
	mCtx *module.ModuleContext
}

func (m *mod) RegisterRoutes(router *gin.Engine) {
	repo := repository.NewCustomerRepository(m.mCtx.DBCtx)
	serv := service.NewCustomerService(repo)
	hdl := handler.NewCustomerHandler(serv)

	rCustomer := router.Group("/api/v1/customers")
	{
		rCustomer.POST("", hdl.CreateCustomer)
	}
}
