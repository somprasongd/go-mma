package notifications

import (
	"go-mma/modules/notifications/service"
	"go-mma/shared/common/module"
	"go-mma/shared/common/registry"

	"github.com/gin-gonic/gin"
)

const (
	NotificationServiceKey registry.ServiceKey = "NotificationService"
)

type mod struct {
	mCtx    *module.ModuleContext
	notiSvc service.NotificationService
}

func NewModule(mCtx *module.ModuleContext) module.Module {
	return &mod{mCtx: mCtx}
}

func (m *mod) Init(reg registry.ServiceRegistry) error {
	m.notiSvc = service.NewNotificationService()

	reg.Register(NotificationServiceKey, m.notiSvc)

	return nil
}

func (m *mod) RegisterRoutes(router *gin.Engine) {

}
