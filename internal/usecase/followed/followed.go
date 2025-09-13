package followed

import (
	"context"
	"fmt"
	"log/slog"
	"love-signal-users/internal/dto"
	"love-signal-users/internal/entity"
	"love-signal-users/internal/lib/logger/sl"
)

// Repository is a repository for followed users use-case.
type Repository interface {
	Follows(ctx context.Context, userID int64) ([]dto.Follow, error)
}

// UseCase is a use-case for getting followed users.
type UseCase struct {
	log  *slog.Logger
	repo Repository
}

// New returns new followed users use-case.
func New(log *slog.Logger, repo Repository) *UseCase {
	return &UseCase{
		log:  log,
		repo: repo,
	}
}

// Execute executes the use-case for getting followed users.
func (uc *UseCase) Execute(ctx context.Context, userID int64) ([]entity.Follow, error) {
	const op = "usecase.followed.Execute"

	log := uc.log.With(
		slog.String("op", op),
		slog.Int64("user ID", userID),
	)

	follows, err := uc.repo.Follows(ctx, userID)
	if err != nil {
		log.Error("error getting followed users data by user ID", sl.Err(err))

		return []entity.Follow{}, fmt.Errorf("%s: %w", op, err)
	}

	followEntities := make([]entity.Follow, len(follows))
	for i, f := range follows {
		followEntities[i] = entity.NewFollow(f)
	}

	return followEntities, nil
}
