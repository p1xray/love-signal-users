package app

import (
	"log/slog"
	grpcapp "love-signal-users/internal/app/grpc"
	"love-signal-users/internal/service"
	"love-signal-users/internal/storage/sqlite"
)

// App is an application.
type App struct {
	GRPCServer *grpcapp.App
}

// New creates a new application.
func New(
	log *slog.Logger,
	grpcPort int,
	storagePath string,
) *App {
	storage, err := sqlite.New(storagePath)
	if err != nil {
		panic(err)
	}

	usersService := service.New(log, storage)

	grpcApp := grpcapp.New(log, grpcPort, usersService)

	return &App{
		GRPCServer: grpcApp,
	}
}
