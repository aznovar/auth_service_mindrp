package app

import (
	grpcapp "auth_service/internal/app/grpc"
	"auth_service/internal/config"
	"auth_service/internal/services/auth"
	"auth_service/internal/storage/postgres"
	"log/slog"
)

type App struct {
	GRPCServer *grpcapp.App
}

func New(
	log *slog.Logger,
	cfg *config.Config,
) *App {
	storage, err := postgres.New(cfg.DB)
	if err != nil {
		panic(err)
	}

	authService := auth.New(log, storage, storage, storage, cfg.TokenTTL)

	grpcApp := grpcapp.New(log, authService, cfg.GRPC.Port)

	return &App{
		GRPCServer: grpcApp,
	}
}
