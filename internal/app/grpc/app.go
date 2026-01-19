package grpcapp

import (
	authgrpc "go-gRPC-sso/internal/grpc/auth"
	"log/slog"

	"google.golang.org/grpc"
)

type App struct {
	logger     *slog.Logger
	gRPCServer *grpc.Server
	port       int
}

func New(logger *slog.Logger, port int) *App {
	gRPCServer := grpc.NewServer()
	authgrpc.Register(gRPCServer)

	return &App{
		logger:     logger,
		gRPCServer: gRPCServer,
		port:       port,
	}
}
