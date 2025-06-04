package main

import (
	"github.com/legenda-hortici/hw-6-auth-service/internal/app"
	"github.com/legenda-hortici/hw-6-auth-service/internal/config"
	"github.com/legenda-hortici/hw-6-auth-service/pkg/logger"
	"os"
	"os/signal"
)

func main() {
	cfg := config.NewConfig()

	log, err := logger.NewLogger(cfg.LogLevel)
	if err != nil {
		panic(err)
	}

	application := app.NewApp(log, *cfg)

	log.Info("starting application")

	go func() {
		application.GRPCServer.MustRun()
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	<-stop

	application.GRPCServer.Stop()
	log.Info("gracefully stoped...")

	// TODO: написать тесты

}
