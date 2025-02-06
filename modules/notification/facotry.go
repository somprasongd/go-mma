package notification

import (
	"go-mma/modules"

	"github.com/gin-gonic/gin"
)

func NewModule(mCtx *modules.ModuleContext) modules.Module {
	return &module{mCtx}
}

type module struct {
	mCtx *modules.ModuleContext
}

func (m *module) RegisterRoutes(router *gin.Engine) {

}
