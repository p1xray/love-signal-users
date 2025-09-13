package grpcapp

import (
	"log/slog"
	"love-signal-users/internal/controller"
	"love-signal-users/internal/controller/grpc"
	"love-signal-users/pkg/grpcserver"
)

// App is an gRPC controller application.
type App struct {
	log        *slog.Logger
	port       string
	gRPCServer *grpcserver.Server
}

// New creates new gRPC controller application.
func New(
	log *slog.Logger,
	port string,
	userDataUseCase controller.UserData,
	userDataByExternalIDUseCase controller.UserDataByExternalID,
	followedUsersUseCase controller.Followed,
	followUserUseCase controller.Follow,
	unfollowUserUseCase controller.Unfollow,
) *App {
	gRPCServer := grpcserver.New(grpcserver.WithPort(port))

	grpc.NewRouter(
		gRPCServer.App,
		userDataUseCase,
		userDataByExternalIDUseCase,
		followedUsersUseCase,
		followUserUseCase,
		unfollowUserUseCase,
	)

	return &App{
		log:        log,
		port:       port,
		gRPCServer: gRPCServer,
	}
}

// Start - starts the gRPC controller application.
func (a *App) Start() {
	const op = "grpcapp.Start"

	log := a.log.With(
		slog.String("op", op),
		slog.String("port", a.port),
	)
	log.Info("running gRPC server")

	a.gRPCServer.Start()
}

// Stop - stops the gRPC controller application.
func (a *App) Stop() {
	const op = "grpcapp.Stop"

	log := a.log.With(
		slog.String("op", op),
		slog.String("port", a.port),
	)
	log.Info("stopping gRPC server")

	a.gRPCServer.Stop()
}

// Notify - notifies about gRPC controller application errors.
func (a *App) Notify() <-chan error {
	return a.gRPCServer.Notify()
}
