package unfollow

import (
	"context"
	"fmt"
	"log/slog"
	"love-signal-users/internal/dto"
	"love-signal-users/internal/entity"
	"love-signal-users/pkg/logger/sl"
)

// Repository is a repository for unfollow user use-case.
type Repository interface {
	SaveFollow(ctx context.Context, follow *entity.Follow) error
}

// UseCase is a use-case for unfollowing users.
type UseCase struct {
	log  *slog.Logger
	repo Repository
}

// New returns new unfollow user use-case.
func New(log *slog.Logger, repo Repository) *UseCase {
	return &UseCase{
		log:  log,
		repo: repo,
	}
}

// Execute executes the use-case for unfollowing user.
func (uc *UseCase) Execute(ctx context.Context, followLinkID int64) error {
	const op = "usecase.unfollow.Execute"

	log := uc.log.With(
		slog.String("op", op),
		slog.Int64("follow link ID", followLinkID),
	)

	followDTO := dto.Follow{ID: followLinkID}
	followEntity := entity.NewFollow(followDTO)
	followEntity.SetToRemove()

	if err := uc.repo.SaveFollow(ctx, &followEntity); err != nil {
		log.Error("error saving follow", sl.Err(err))

		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
