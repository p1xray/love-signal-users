package app

import (
	"log/slog"
	grpcapp "love-signal-users/internal/app/grpc"
	"love-signal-users/internal/service"
)

// App is an application.
type App struct {
	GRPCServer *grpcapp.App
}

// New creates a new application.
func New(
	log *slog.Logger,
	grpcPort int,
) *App {
	usersService := service.New(log)

	grpcApp := grpcapp.New(log, grpcPort, usersService)

	return &App{
		GRPCServer: grpcApp,
	}
}
