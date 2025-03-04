package application

import (
	"go-mma/config"
	"go-mma/shared/common/eventbus"
	"go-mma/shared/common/module"
	"go-mma/shared/common/registry"
	"log"
)

type Application struct {
	config config.Config

	httpServer      HTTPServer
	serviceRegistry registry.ServiceRegistry
	eventbus        eventbus.EventBus
}

func New(config config.Config) *Application {
	return &Application{
		config:          config,
		httpServer:      newHTTPServer(config),
		serviceRegistry: registry.NewServiceRegistry(),
		eventbus:        eventbus.NewInMemoryEventBus(),
	}
}

func (app *Application) Run() error {
	log.Printf("Starting server on port %d\n", app.config.HTTPPort)
	app.httpServer.Start()

	return nil
}

func (app *Application) RegisterModules(modules []module.Module) {
	for _, module := range modules {
		if err := module.Init(app.serviceRegistry, app.eventbus); err != nil {
			log.Fatalf("module initialization error: %v", err)
		}
		module.RegisterRoutes(app.httpServer.Router())
	}
}

func (app *Application) Shutdown() error {
	log.Println("Shutting down server")
	if err := app.httpServer.Shutdown(); err != nil {
		log.Printf("Server shutdown error: %v\n", err)
	}
	log.Println("Server stopped")
	return nil
}
