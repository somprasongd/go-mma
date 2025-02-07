package application

import (
	"go-mma/config"
	"go-mma/shared/common/module"
	"log"
)

type Application struct {
	config config.Config

	httpServer HTTPServer
}

func New(config config.Config) *Application {
	return &Application{
		config:     config,
		httpServer: newHTTPServer(config),
	}
}

func (app *Application) Run() error {
	log.Printf("Starting server on port %d\n", app.config.HTTPPort)
	app.httpServer.Start()

	return nil
}

func (app *Application) RegisterModules(modules []module.Module) {
	for _, module := range modules {
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
