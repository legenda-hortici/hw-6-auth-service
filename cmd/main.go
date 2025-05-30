package main

import (
	"os"
	"os/signal"
	"skillsRockAuthService/internal/app"
	"skillsRockAuthService/internal/config"
	"skillsRockAuthService/internal/logger"
)

func main() {
	cfg := config.NewConfig()

	log, err := logger.NewLogger(cfg.LogLevel)
	if err != nil {
		panic(err)
	}

	log.Info(cfg)

	application := app.NewApp(log, *cfg)

	go func() {
		application.GRPCServer.MustRun()
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	<-stop

	application.GRPCServer.Stop()
	log.Info("gracefully stoped...")
}
