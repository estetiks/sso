package app

import (
	"log/slog"
	"time"

	grpcapp "github.com/estetiks/sso/internal/app/grpc"
	"github.com/estetiks/sso/internal/services/auth"
	"github.com/estetiks/sso/internal/storage/sqlite"
)

type App struct {
	GRPCSrv *grpcapp.App
}

// New returns a new grpc application
func New(log *slog.Logger, grpcPort int, storagePath string, tokenTTL time.Duration) *App {

	storage, err := sqlite.New(storagePath)

	if err != nil {
		panic(err)
	}

	authApp := auth.New(log, storage, storage, storage, tokenTTL)

	grpcApp := grpcapp.New(log, authApp, grpcPort)

	return &App{
		GRPCSrv: grpcApp,
	}
}
