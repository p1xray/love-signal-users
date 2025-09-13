package user

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"love-signal-users/internal/dto"
	"love-signal-users/internal/entity"
	"love-signal-users/internal/infrastructure"
	"love-signal-users/internal/usecase"
	"love-signal-users/pkg/logger/sl"
)

// Repository is a repository for user data use-case.
type Repository interface {
	User(ctx context.Context, id int64) (dto.User, error)
}

// UseCase is a use-case for getting user data.
type UseCase struct {
	log  *slog.Logger
	repo Repository
}

// New returns new user data use-case.
func New(log *slog.Logger, repo Repository) *UseCase {
	return &UseCase{
		log:  log,
		repo: repo,
	}
}

// Execute executes the use-case for getting user data.
func (uc *UseCase) Execute(ctx context.Context, id int64) (entity.User, error) {
	const op = "usecase.user.Execute"

	log := uc.log.With(
		slog.String("op", op),
		slog.Int64("user ID", id),
	)

	user, err := uc.repo.User(ctx, id)
	if err != nil {
		if errors.Is(err, infrastructure.ErrEntityNotFound) {
			log.Warn("user not found", sl.Err(err))

			return entity.User{}, fmt.Errorf("%s: %w", op, usecase.ErrUserNotFound)
		}

		log.Error("error getting user data by user ID", sl.Err(err))

		return entity.User{}, fmt.Errorf("%s: %w", op, err)
	}

	userEntity := entity.NewUser(user)

	return userEntity, nil
}
