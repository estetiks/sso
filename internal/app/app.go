package app

import (
	"log/slog"
	"time"

	grpcapp "github.com/estetiks/sso/internal/app/grpc"
)

type App struct {
	GRPCSrv *grpcapp.App
}

// New returns a new grpc application
func New(log *slog.Logger, grpcPort int, storagePath string, tokenTTL time.Duration) *App {
	// TODO: initialize storage

	//TODO: init auth service

	grpcApp := grpcapp.New(log, grpcPort)

	return &App{
		GRPCSrv: grpcApp,
	}
}
