package notifications

import (
	eventhandlers "go-mma/modules/notifications/event_handlers"
	"go-mma/modules/notifications/service"
	"go-mma/shared/common/eventbus"
	"go-mma/shared/common/module"
	"go-mma/shared/common/registry"
	"go-mma/shared/messaging"

	"github.com/gin-gonic/gin"
)

type mod struct {
	mCtx    *module.ModuleContext
	notiSvc service.NotificationService
}

func NewModule(mCtx *module.ModuleContext) module.Module {
	return &mod{mCtx: mCtx}
}

func (m *mod) Init(reg registry.ServiceRegistry, eventbus eventbus.EventBus) error {
	m.notiSvc = service.NewNotificationService()

	eventbus.Subscribe(messaging.OrderCreatedIntegrationEventName, eventhandlers.NewOrderCreatedIntegrationEventHandler(m.notiSvc))

	return nil
}

func (m *mod) RegisterRoutes(router *gin.Engine) {

}
