package main

import (
	"go-mma/application"
	"go-mma/config"
	"log"
	"os"
	"os/signal"
	"syscall"

	"go-mma/modules/customers"
	"go-mma/modules/notifications"
	"go-mma/modules/orders"
	"go-mma/shared/common/module"
	"go-mma/shared/common/storage/db"
)

func main() {
	config, err := config.Load()
	if err != nil {
		log.Panic(err)
	}

	transactor, dbCtx, closeDB, err := db.New(config.DSN)
	if err != nil {
		log.Panic(err)
	}
	defer func() {
		if err := closeDB(); err != nil {
			log.Println("Error closing database:", err)
		}
	}()

	app := application.New(*config)

	mCtx := module.NewModuleContext(transactor, dbCtx)
	app.RegisterModules([]module.Module{
		notifications.NewModule(mCtx),
		customers.NewModule(mCtx),
		orders.NewModule(mCtx),
	})

	app.Run()
	// Handle graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	<-stop
	app.Shutdown()
}
