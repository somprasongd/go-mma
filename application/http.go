package application

import (
	"context"
	"fmt"
	"go-mma/config"
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
	// router.Use(gin.Logger())
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
