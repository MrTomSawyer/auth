package app

import (
	grpcapp "github.com/MrTomSawyer/sso/internal/app/grpc"
	"log/slog"
	"time"
)

type App struct {
	GRPCServer *grpcapp.App
}

func New(log *slog.Logger, gRPCPort int, storagePath string, tokenTTL time.Duration) *App {
	grpcApp := grpcapp.New(log, gRPCPort)
	return &App{
		GRPCServer: grpcApp,
	}
}
