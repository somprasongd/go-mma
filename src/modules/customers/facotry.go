package customers

import (
	"go-mma/modules/customers/handler"
	"go-mma/modules/customers/repository"
	"go-mma/modules/customers/service"
	"go-mma/shared/common/module"
	"go-mma/shared/common/registry"

	"github.com/gin-gonic/gin"

	customerContracts "go-mma/shared/contracts/customer_contracts"
)

type mod struct {
	mCtx    *module.ModuleContext
	custSvc customerContracts.CustomerService
}

func NewModule(mCtx *module.ModuleContext) module.Module {
	return &mod{mCtx: mCtx}
}

func (m *mod) Init(reg registry.ServiceRegistry) error {
	repo := repository.NewCustomerRepository(m.mCtx.DBCtx)
	m.custSvc = service.NewCustomerService(repo)

	reg.Register(customerContracts.CustomerServiceKey, m.custSvc)

	return nil
}

func (m *mod) RegisterRoutes(router *gin.Engine) {
	hdl := handler.NewCustomerHandler(m.custSvc)

	rCustomer := router.Group("/api/v1/customers")
	{
		rCustomer.POST("", hdl.CreateCustomer)
	}
}
