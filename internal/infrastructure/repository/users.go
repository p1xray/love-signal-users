package repository

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"love-signal-users/internal/dto"
	"love-signal-users/internal/entity"
	"love-signal-users/internal/infrastructure"
	"love-signal-users/internal/infrastructure/converter"
	"love-signal-users/internal/infrastructure/storage/models"
	"love-signal-users/pkg/logger/sl"
)

const emptyID = 0

type Storage interface {
	Users(ctx context.Context, ids []int64) ([]models.User, error)
	User(ctx context.Context, id int64) (models.User, error)
	UserByExternalID(ctx context.Context, externalID int64) (models.User, error)
	FollowsByUserID(ctx context.Context, userID int64) ([]models.Follow, error)
	CreateFollow(ctx context.Context, follow models.Follow) (int64, error)
	UpdateFollow(ctx context.Context, follow models.Follow) error
	RemoveFollow(ctx context.Context, id int64) error
}

type Users struct {
	log     *slog.Logger
	storage Storage
}

func NewUsersRepository(log *slog.Logger, storage Storage) *Users {
	return &Users{
		log:     log,
		storage: storage,
	}
}

func (u *Users) User(ctx context.Context, id int64) (dto.User, error) {
	const op = "repository.users.User"

	log := u.log.With(
		slog.String("op", op),
		slog.Int64("user ID", id),
	)

	user, err := u.storage.User(ctx, id)
	if err != nil {
		if errors.Is(err, infrastructure.ErrEntityNotFound) {
			log.Warn("user not found", sl.Err(err))
		} else {
			log.Error("error getting user", sl.Err(err))
		}

		return dto.User{}, fmt.Errorf("%s: %w", op, err)
	}

	userDTO := converter.ToUserDTO(user)

	return userDTO, nil
}

func (u *Users) UserByExternalID(ctx context.Context, externalID int64) (dto.User, error) {
	const op = "repository.users.UserByExternalID"

	log := u.log.With(
		slog.String("op", op),
		slog.Int64("user external ID", externalID),
	)

	user, err := u.storage.UserByExternalID(ctx, externalID)
	if err != nil {
		if errors.Is(err, infrastructure.ErrEntityNotFound) {
			log.Warn("user not found", sl.Err(err))
		} else {
			log.Error("error getting user", sl.Err(err))
		}

		return dto.User{}, fmt.Errorf("%s: %w", op, err)
	}

	userDTO := converter.ToUserDTO(user)

	return userDTO, nil
}

func (u *Users) Follows(ctx context.Context, userID int64) ([]dto.Follow, error) {
	const op = "repository.users.Follows"

	log := u.log.With(
		slog.String("op", op),
		slog.Int64("user ID", userID),
	)

	follows, err := u.storage.FollowsByUserID(ctx, userID)
	if err != nil {
		log.Error("error getting follows", sl.Err(err))

		return []dto.Follow{}, fmt.Errorf("%s: %w", op, err)
	}

	userIDs := make([]int64, 0, len(follows)*2)
	for _, f := range follows {
		userIDs = append(userIDs, f.FollowingUserID)
		userIDs = append(userIDs, f.FollowedUserID)
	}

	users, err := u.storage.Users(ctx, userIDs)
	if err != nil {
		log.Error("error getting users", sl.Err(err))

		return []dto.Follow{}, fmt.Errorf("%s: %w", op, err)
	}

	followsDTO := make([]dto.Follow, len(follows))
	for i, f := range follows {
		followDTO, err := converter.ToFollowDTO(f, users)
		if err != nil {
			return []dto.Follow{}, fmt.Errorf("%s: %w", op, err)
		}

		followsDTO[i] = followDTO
	}

	return followsDTO, nil
}

func (u *Users) SaveFollow(ctx context.Context, follow *entity.Follow) error {
	const op = "repository.users.SaveFollow"

	log := u.log.With(
		slog.String("op", op),
	)

	if follow.IsToCreate() {
		if err := u.createFollow(ctx, follow); err != nil {
			log.Error("error creating follow", sl.Err(err))

			return fmt.Errorf("%s: %w", op, err)
		}
	}

	if follow.IsToUpdate() {
		if err := u.updateFollow(ctx, follow); err != nil {
			log.Error("error updating follow", sl.Err(err))

			return fmt.Errorf("%s: %w", op, err)
		}
	}

	if follow.IsToRemove() {
		if err := u.removeFollow(ctx, follow); err != nil {
			log.Error("error removing follow", sl.Err(err))

			return fmt.Errorf("%s: %w", op, err)
		}
	}

	return nil
}

func (u *Users) createFollow(ctx context.Context, follow *entity.Follow) error {
	followStorageModel := converter.ToFollowStorage(follow, models.FollowCreated())

	id, err := u.storage.CreateFollow(ctx, followStorageModel)
	if err != nil {
		return err
	}

	follow.ID = id
	follow.ResetDataStatus()

	return nil
}

func (u *Users) updateFollow(ctx context.Context, follow *entity.Follow) error {
	if follow.ID == emptyID {
		return infrastructure.ErrRequireIDToUpdate
	}

	userStorageModel := converter.ToFollowStorage(follow, models.FollowUpdated())

	err := u.storage.UpdateFollow(ctx, userStorageModel)
	if err != nil {
		return err
	}

	follow.ResetDataStatus()

	return nil
}

func (u *Users) removeFollow(ctx context.Context, follow *entity.Follow) error {
	if follow.ID == emptyID {
		return infrastructure.ErrRequireIDToRemove
	}

	err := u.storage.RemoveFollow(ctx, follow.ID)
	if err != nil {
		return err
	}

	follow.ResetDataStatus()

	return nil
}
