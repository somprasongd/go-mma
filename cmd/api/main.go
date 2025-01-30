package main

import (
	"go-mma/application"
	"go-mma/config"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	config, err := config.Load()
	if err != nil {
		log.Panic(err)
	}

	app := application.New(*config)
	app.RegisterRoutes()
	app.Run()
	// Handle graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	<-stop
	app.Shutdown()
}
