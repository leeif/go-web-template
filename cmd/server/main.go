package main

import (
	"context"
	"database/sql"
	"log"
	"os"
	"os/signal"
	"time"

	"go.uber.org/fx"

	_ "github.com/go-sql-driver/mysql"
	"github.com/leeif/go-web-template/config"
	"github.com/leeif/go-web-template/database"
	mylog "github.com/leeif/go-web-template/log"
	"github.com/leeif/go-web-template/manage"
	"github.com/leeif/go-web-template/route"
	"github.com/leeif/go-web-template/server"
)

// VERSION is the build version
var VERSION = ""

func register(route *route.Route, db *sql.DB, config *config.Config, logger *mylog.Log) error {
	// register routes
	route.Register()

	return nil
}

func main() {

	app := fx.New(
		fx.Provide(
			func() []string {
				return os.Args
			},
			func() string {
				return VERSION
			},
			config.NewConfig,
			database.NewDatabase,
			mylog.NewLogger,
			server.NewMux,
			route.NewRouter,
			manage.NewManager,
		),
		fx.Invoke(register),
		fx.NopLogger,
	)
	startCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	if err := app.Start(startCtx); err != nil {
		log.Fatal(err)
	}

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)

	<-quit

	stopCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	if err := app.Stop(stopCtx); err != nil {
		log.Fatal(err)
	}
}
