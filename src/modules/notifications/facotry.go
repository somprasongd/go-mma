package notifications

import (
	eventhandlers "go-mma/modules/notifications/event_handlers"
	"go-mma/modules/notifications/features"
	"go-mma/shared/common/eventbus"
	"go-mma/shared/common/mediator"
	"go-mma/shared/common/module"
	"go-mma/shared/common/registry"
	"go-mma/shared/messaging"

	"github.com/gin-gonic/gin"
)

type mod struct {
	mCtx *module.ModuleContext
}

func NewModule(mCtx *module.ModuleContext) module.Module {
	return &mod{mCtx: mCtx}
}

func (m *mod) Init(reg registry.ServiceRegistry, eventbus eventbus.EventBus) error {
	mediator.Register[*features.SendEmailCommand, *mediator.NoResponse](features.NewSendEmailHandler())

	eventbus.Subscribe(messaging.OrderCreatedIntegrationEventName, eventhandlers.NewOrderCreatedIntegrationEventHandler())

	return nil
}

func (m *mod) RegisterRoutes(router *gin.Engine) {

}
