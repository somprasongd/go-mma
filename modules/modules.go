package modules

import (
	"go-mma/util/transactor"

	"github.com/gin-gonic/gin"
)

type Module interface {
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
