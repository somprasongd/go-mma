package notifications

import (
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

}
