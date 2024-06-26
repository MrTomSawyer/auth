package main

import (
	"github.com/MrTomSawyer/sso/internal/app"
	"github.com/MrTomSawyer/sso/internal/config"
	"github.com/MrTomSawyer/sso/pkg/logger"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	cfg := config.MustLoad()
	log := logger.SetupLogger(cfg.Env)

	server := app.New(log, cfg.GRPC.Port, cfg.StoragePath, cfg.TokenTTL, cfg.SecretKey)

	go server.GRPCServer.MustRun()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	<-stop

	server.GRPCServer.Stop()
	log.Info("server stopped")
}

//go run main.go --config=../../config/local.yaml
