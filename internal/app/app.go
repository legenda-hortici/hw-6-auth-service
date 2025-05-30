package app

import (
	"go.uber.org/zap"
	grpserver "skillsRockAuthService/internal/app/grpc"
	"skillsRockAuthService/internal/config"
	"skillsRockAuthService/internal/services"
	"skillsRockAuthService/internal/storage"
)

type App struct {
	GRPCServer *grpserver.App
}

func NewApp(
	log *zap.SugaredLogger,
	cfg config.Config,
) *App {

	storage, err := storage.NewStorage(log, cfg)
	if err != nil {
		panic(err)
	}

	authService := services.NewAuthService(log, storage)

	grpcApp := grpserver.NewGRPCApp(log, authService, cfg.GRPC.ListenPort)

	return &App{
		GRPCServer: grpcApp,
	}
}
