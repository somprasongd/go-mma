package application

import (
	"context"
	"fmt"
	"go-mma/config"
	"go-mma/handlers"
	"log"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Application struct {
	config config.Config

	httpRouter *gin.Engine
	httpServer *http.Server
}

func New(config config.Config) *Application {
	httpRouter := gin.Default()
	return &Application{
		config:     config,
		httpRouter: httpRouter,
		httpServer: &http.Server{
			Addr:    fmt.Sprintf(":%d", config.HTTPPort),
			Handler: httpRouter,
		},
	}
}

func (app *Application) Run() error {
	app.startHTTPServer()

	return nil
}

func (app *Application) startHTTPServer() {
	go func() {
		log.Printf("Starting server on port %d", app.config.HTTPPort)
		if err := app.httpServer.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("Error starting server: %v", err)
		}
	}()
}

func (app *Application) RegisterRoutes() {
	// global middleware
	app.httpRouter.Use(gin.Logger())
	// Recovery middleware recovers from any panics and writes a 500 if there was one.
	app.httpRouter.Use(gin.Recovery())
	app.httpRouter.Use(cors.Default()) // allows all origins

	// route handler
	app.httpRouter.GET("/", func(c *gin.Context) {
		// time.Sleep(10 * time.Second)
		c.JSON(http.StatusOK, gin.H{
			"message": "Hello World",
		})
	})
	v1 := app.httpRouter.Group("/api/v1")

	rCustomer := v1.Group("/customers")
	{
		hdl := handlers.NewCustomerHandler()
		rCustomer.POST("", hdl.CreateCustomer)
	}

	rOrder := v1.Group("/orders")
	{
		hdl := handlers.NewOrderHandler()
		rOrder.POST("", hdl.CreateOrder)
		rOrder.DELETE("", hdl.CancelOrder)
	}

}

func (app *Application) Shutdown() error {
	log.Println("Shutting down server")
	if err := app.httpServer.Shutdown(context.Background()); err != nil {
		log.Printf("Server shutdown error: %v", err)
	}
	log.Println("Server stopped")
	return nil
}
