package app

import (
	"log/slog"
	"time"

	"google.golang.org/grpc"
)

type App struct {
	GRPCServer *grpc.Server
}

func New(logger *slog.Logger, grpcPort int, storagePath string, tokenTTL time.Duration) *App {
	// TODO: init storage
}
