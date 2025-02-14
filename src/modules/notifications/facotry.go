package notifications

import (
	"go-mma/modules/notifications/service"
	"go-mma/shared/common/module"
	"go-mma/shared/common/registry"
	notificationsContracts "go-mma/shared/contracts/notification_contracts"

	"github.com/gin-gonic/gin"
)

type mod struct {
	mCtx    *module.ModuleContext
	notiSvc notificationsContracts.NotificationService
}

func NewModule(mCtx *module.ModuleContext) module.Module {
	return &mod{mCtx: mCtx}
}

func (m *mod) Init(reg registry.ServiceRegistry) error {
	m.notiSvc = service.NewNotificationService()

	reg.Register(notificationsContracts.NotificationServiceKey, m.notiSvc)

	return nil
}

func (m *mod) RegisterRoutes(router *gin.Engine) {

}
