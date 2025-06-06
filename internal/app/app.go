package app

import (
	grpserver "github.com/legenda-hortici/hw-6-auth-service/internal/app/grpc"
	"github.com/legenda-hortici/hw-6-auth-service/internal/config"
	"github.com/legenda-hortici/hw-6-auth-service/internal/services"
	"github.com/legenda-hortici/hw-6-auth-service/internal/storage"
	"go.uber.org/zap"
)

type App struct {
	GRPCServer *grpserver.App
}

func NewApp(
	log *zap.SugaredLogger,
	cfg config.Config,
) *App {

	store, err := storage.NewStorage(cfg)
	if err != nil {
		panic(err)
	}

	authService := services.NewAuthService(cfg, log, store, store)

	log.Info("database initialized successfully")

	grpcApp := grpserver.NewGRPCApp(log, authService, cfg.GRPC.ListenPort)

	return &App{
		GRPCServer: grpcApp,
	}
}
