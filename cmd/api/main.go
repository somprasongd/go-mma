package main

import (
	"context"
	"fmt"
	"go-mma/config"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	config, err := config.Load()
	if err != nil {
		log.Panic(err)
	}

	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		time.Sleep(10 * time.Second)
		c.JSON(http.StatusOK, gin.H{
			"message": "Hello World",
		})
	})

	httpServer := &http.Server{
		Addr:    fmt.Sprintf(":%d", config.HTTPPort),
		Handler: r,
	}
	go func() {
		log.Printf("Starting server on port %d", config.HTTPPort)
		if err := httpServer.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("Error starting server: %v", err)
		}
	}()

	// Handle graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	<-stop
	log.Println("Shutting down server")
	if err := httpServer.Shutdown(context.Background()); err != nil {
		log.Printf("Server shutdown error: %v", err)
	}
	log.Println("Server stopped")
}
