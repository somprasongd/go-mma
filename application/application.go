package application

import (
	"go-mma/config"
	"go-mma/data/sqldb"
	"log"
)

type Application struct {
	config config.Config

	db         sqldb.DBContext
	httpServer HTTPServer
}

func New(config config.Config, db sqldb.DBContext) *Application {
	return &Application{
		config:     config,
		db:         db,
		httpServer: newHTTPServer(config),
	}
}

func (app *Application) Run() error {
	app.httpServer.Start()

	return nil
}

func (app *Application) RegisterRoutes() {
	log.Printf("Starting server on port %d", app.config.HTTPPort)
	registerRoutes(app.httpServer.Router(), app.db)
}

func (app *Application) Shutdown() error {
	log.Println("Shutting down server")
	if err := app.httpServer.Shutdown(); err != nil {
		log.Printf("Server shutdown error: %v", err)
	}
	log.Println("Server stopped")
	return nil
}
