package app

import (
	grpcapp "go-gRPC-sso/internal/app/grpc"
	"go-gRPC-sso/internal/services/auth"
	"go-gRPC-sso/internal/storage/sqlite"
	"log/slog"
	"time"
)

type App struct {
	GRPCServer *grpcapp.App
}

func New(logger *slog.Logger, grpcPort int, storagePath string, tokenTTL time.Duration) *App {
	storage, err := sqlite.New(storagePath)
	if err != nil {
		panic(err)
	}

	authService := auth.New(logger, storage, storage, storage, tokenTTL)
	grpcApp := grpcapp.New(logger, authService, grpcPort)

	return &App{
		GRPCServer: grpcApp,
	}
}
