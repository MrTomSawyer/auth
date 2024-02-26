package main

import (
	"github.com/MrTomSawyer/sso/internal/app"
	"github.com/MrTomSawyer/sso/internal/config"
	"github.com/MrTomSawyer/sso/internal/logger"
)

func main() {
	cfg := config.MustLoad()
	log := logger.SetupLogger(cfg.Env)

	server := app.New(log, cfg.GRPC.Port, cfg.StoragePath, cfg.TokenTTL)

	server.GRPCServer.MustRun()
}
