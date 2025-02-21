package module

import (
	"go-mma/shared/common/eventbus"
	"go-mma/shared/common/registry"
	"go-mma/shared/common/storage/db/transactor"

	"github.com/gin-gonic/gin"
)

type Module interface {
	// Init registers the module's services into the registry.
	Init(reg registry.ServiceRegistry, eventbus eventbus.EventBus) error
	// RegisterRoutes registers the module's routes to the Gin app.
	RegisterRoutes(r *gin.Engine)
}

type ModuleContext struct {
	Transactor transactor.Transactor
	DBCtx      transactor.DBContext
}

func NewModuleContext(transactor transactor.Transactor, dbCtx transactor.DBContext) *ModuleContext {
	return &ModuleContext{
		Transactor: transactor,
		DBCtx:      dbCtx,
	}
}
