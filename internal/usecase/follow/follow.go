package follow

import (
	"context"
	"fmt"
	"log/slog"
	"love-signal-users/internal/dto"
	"love-signal-users/internal/entity"
	"love-signal-users/internal/lib/logger/sl"
)

// Repository is a repository for follow user use-case.
type Repository interface {
	SaveFollow(ctx context.Context, follow *entity.Follow) error
}

// UseCase is a use-case for following users.
type UseCase struct {
	log  *slog.Logger
	repo Repository
}

// New returns new follow user use-case.
func New(log *slog.Logger, repo Repository) *UseCase {
	return &UseCase{
		log:  log,
		repo: repo,
	}
}

// Execute executes the use-case for following user.
func (uc *UseCase) Execute(
	ctx context.Context,
	userID int64,
	userIDToFollow int64,
) error {
	const op = "usecase.follow.Execute"

	log := uc.log.With(
		slog.String("op", op),
		slog.Int64("user ID", userID),
		slog.Int64("user ID to follow", userIDToFollow),
	)

	followDTO := dto.Follow{
		FollowingUser: dto.User{ID: userID},
		FollowedUser:  dto.User{ID: userIDToFollow},
	}

	followEntity := entity.NewFollow(followDTO)
	followEntity.SetToCreate()

	if err := uc.repo.SaveFollow(ctx, &followEntity); err != nil {
		// TODO: handle check FK constraint on user

		log.Error("error saving follow", sl.Err(err))

		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
