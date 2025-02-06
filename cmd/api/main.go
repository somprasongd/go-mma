package main

import (
	"go-mma/application"
	"go-mma/config"
	"go-mma/data/sqldb"
	"go-mma/modules"
	"go-mma/modules/customer"
	"go-mma/modules/order"
	"go-mma/util/transactor"
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

	db, closeDB, err := sqldb.New(config.DSN)
	if err != nil {
		log.Panic(err)
	}
	defer func() {
		if err := closeDB(); err != nil {
			log.Println("Error closing database:", err)
		}
	}()

	transactor, dbCtx := transactor.New(db.DB(), transactor.NestedTransactionsSavepoints)

	app := application.New(*config, db)

	mCtx := modules.NewModuleContext(transactor, dbCtx)
	app.RegisterModules([]modules.Module{
		customer.NewModule(mCtx),
		order.NewModule(mCtx),
	})

	app.Run()
	// Handle graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	<-stop
	app.Shutdown()
}
