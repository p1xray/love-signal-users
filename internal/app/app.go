package app

import (
	"log/slog"
	grpcapp "love-signal-users/internal/app/grpc"
	"love-signal-users/internal/config"
	"love-signal-users/internal/infrastructure/repository"
	"love-signal-users/internal/infrastructure/storage/sqlite"
	"love-signal-users/internal/usecase/externaluser"
	"love-signal-users/internal/usecase/follow"
	"love-signal-users/internal/usecase/followed"
	"love-signal-users/internal/usecase/unfollow"
	"love-signal-users/internal/usecase/user"
	"love-signal-users/pkg/logger/sl"
	"os"
	"os/signal"
	"syscall"
)

// App is an application.
type App struct {
	log     *slog.Logger
	grpcApp *grpcapp.App
}

// New creates a new application.
func New(
	log *slog.Logger,
	cfg *config.Config,
) *App {
	// Storages.
	storage, err := sqlite.New(cfg.StoragePath)
	if err != nil {
		panic(err)
	}

	// Repositories.
	usersRepository := repository.NewUsersRepository(log, storage)

	// Use-cases.
	userDataUseCase := user.New(log, usersRepository)
	userDataByExternalIDUseCase := externaluser.New(log, usersRepository)
	followedUsersUseCase := followed.New(log, usersRepository)
	followUserUseCase := follow.New(log, usersRepository)
	unfollowUserUseCase := unfollow.New(log, usersRepository)

	grpcApp := grpcapp.New(
		log,
		cfg.GRPC.Port,
		userDataUseCase,
		userDataByExternalIDUseCase,
		followedUsersUseCase,
		followUserUseCase,
		unfollowUserUseCase,
	)

	return &App{
		log:     log,
		grpcApp: grpcApp,
	}
}

// Start - starts the application.
func (a *App) Start() {
	const op = "app.Start"

	log := a.log.With(slog.String("op", op))
	log.Info("starting application")

	a.grpcApp.Start()
}

// GracefulStop - gracefully stops the application.
func (a *App) GracefulStop() {
	const op = "app.GracefulStop"

	log := a.log.With(slog.String("op", op))

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	select {
	case s := <-stop:
		log.Info("signal received from OS", slog.String("signal:", s.String()))
	case err := <-a.grpcApp.Notify():
		log.Error("received an error from the gRPC server:", sl.Err(err))
	}

	log.Info("stopping application")

	a.grpcApp.Stop()
}
