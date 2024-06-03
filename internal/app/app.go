package app

import (
	grpcapp "github.com/MrTomSawyer/sso/internal/app/grpc"
	"github.com/MrTomSawyer/sso/internal/services/auth"
	"github.com/MrTomSawyer/sso/internal/storage/sqlite"
	"log/slog"
	"time"
)

type App struct {
	GRPCServer *grpcapp.App
}

func New(log *slog.Logger, gRPCPort int, storagePath string, tokenTTL time.Duration, secret string) *App {
	storage, err := sqlite.New(storagePath)
	if err != nil {
		panic(err)
	}

	authService := auth.New(log, storage, storage, tokenTTL)

	grpcApp := grpcapp.New(log, authService, gRPCPort, secret)
	return &App{
		GRPCServer: grpcApp,
	}
}
