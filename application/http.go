package application

import (
	"context"
	"fmt"
	"go-mma/config"
	"go-mma/data/sqldb"
	"go-mma/handlers"
	"log"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type HTTPServer interface {
	Start()
	Shutdown() error
	Router() *gin.Engine
}

type httpServer struct {
	config config.Config
	router *gin.Engine
	server *http.Server
}

func newHTTPServer(config config.Config) HTTPServer {
	httpRouter := newRouter()
	return &httpServer{
		router: httpRouter,
		server: &http.Server{
			Addr:    fmt.Sprintf(":%d", config.HTTPPort),
			Handler: httpRouter,
		},
	}
}

func newRouter() *gin.Engine {
	router := gin.Default()

	// add default middlewares here
	router.Use(gin.Logger())
	// Recovery middleware recovers from any panics and writes a 500 if there was one.
	router.Use(gin.Recovery())
	router.Use(cors.Default()) // allows all origins

	return router
}

func (s *httpServer) Start() {
	go func() {
		if err := s.server.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("Error starting server: %v", err)
		}
	}()
}

func (s *httpServer) Shutdown() error {
	return s.server.Shutdown(context.Background())
}

func (s *httpServer) Router() *gin.Engine {
	return s.router
}

func registerRoutes(r *gin.Engine, dbCtx sqldb.DBContext) {
	r.GET("/", func(c *gin.Context) {
		// time.Sleep(10 * time.Second)
		c.JSON(http.StatusOK, gin.H{
			"message": "Hello World",
		})
	})

	v1 := r.Group("/api/v1")

	rCustomer := v1.Group("/customers")
	{
		hdl := handlers.NewCustomerHandler()
		rCustomer.POST("", hdl.CreateCustomer)
	}

	rOrder := v1.Group("/orders")
	{
		hdl := handlers.NewOrderHandler()
		rOrder.POST("", hdl.CreateOrder)
		rOrder.DELETE("/:id", hdl.CancelOrder)
	}
}
