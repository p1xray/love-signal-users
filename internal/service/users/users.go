package users

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"love-signal-users/internal/dto"
	"love-signal-users/internal/lib/logger/sl"
	"love-signal-users/internal/service"
	"love-signal-users/internal/storage"
)

// Service is service for working with user data.
type Service struct {
	log     *slog.Logger
	storage service.UserStorage
}

// New creates a new instance of user service.
func New(log *slog.Logger, storage service.UserStorage) *Service {
	return &Service{
		log:     log,
		storage: storage,
	}
}

// UserData returns information about a user by their ID.
func (s *Service) UserData(
	ctx context.Context,
	userID int64,
) (dto.UserData, error) {
	const op = "users.UserData"

	log := s.log.With(
		slog.String("op", op),
		slog.Int64("user ID", userID),
	)

	userData, err := s.storage.UserData(ctx, userID)
	if err != nil {
		if errors.Is(err, storage.ErrEntityNotFound) {
			log.Warn("user not found", sl.Err(err))

			return dto.UserData{}, fmt.Errorf("%s: %w", op, service.ErrUserNotFound)
		}

		log.Error("error getting user data by user ID", sl.Err(err))

		return dto.UserData{}, fmt.Errorf("%s: %w", op, err)
	}

	return userData, nil
}

// UserDataByExternalID returns information about a user by their external ID.
func (s *Service) UserDataByExternalID(
	ctx context.Context,
	userExternalID int64,
) (dto.UserData, error) {
	const op = "users.UserDataByExternalID"

	log := s.log.With(
		slog.String("op", op),
		slog.Int64("user external ID", userExternalID),
	)

	userData, err := s.storage.UserDataByExternalId(ctx, userExternalID)
	if err != nil {
		if errors.Is(err, storage.ErrEntityNotFound) {
			log.Warn("user not found", sl.Err(err))

			return dto.UserData{}, fmt.Errorf("%s: %w", op, service.ErrUserNotFound)
		}

		log.Error("error getting user data by user external ID", sl.Err(err))

		return dto.UserData{}, fmt.Errorf("%s: %w", op, err)
	}

	return userData, nil
}

// FollowedUsers returns a list of users that the given user is followed to.
func (s *Service) FollowedUsers(
	ctx context.Context,
	userID int64,
) ([]dto.FollowedUser, error) {
	const op = "users.FollowedUsers"

	log := s.log.With(
		slog.String("op", op),
		slog.Int64("user ID", userID),
	)

	followedUsers, err := s.storage.FollowedUsers(ctx, userID)
	if err != nil {
		log.Error("error getting followed users by user ID", sl.Err(err))

		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return followedUsers, nil
}

// FollowUser adds the user with userIDToFollow to the list of followed users with userID.
func (s *Service) FollowUser(
	ctx context.Context,
	userID int64,
	userIDToFollow int64,
) error {
	const op = "users.FollowUser"

	log := s.log.With(
		slog.String("op", op),
		slog.Int64("user ID", userID),
		slog.Int64("user ID to follow", userIDToFollow),
	)

	user, err := s.storage.UserData(ctx, userID)
	if err != nil {
		if errors.Is(err, storage.ErrEntityNotFound) {
			log.Warn("user not found", sl.Err(err))

			return fmt.Errorf("%s: %w", op, service.ErrUserNotFound)
		}

		log.Error("error getting user by user ID", sl.Err(err))

		return fmt.Errorf("%s: %w", op, err)
	}

	userToFollow, err := s.storage.UserData(ctx, userIDToFollow)
	if err != nil {
		if errors.Is(err, storage.ErrEntityNotFound) {
			log.Warn("user not found", sl.Err(err))

			return fmt.Errorf("%s: %w", op, service.ErrUserNotFound)
		}

		log.Error("error getting user to follow by user ID", sl.Err(err))

		return fmt.Errorf("%s: %w", op, err)
	}

	if err := s.storage.AddFollowLink(ctx, user.ID, userToFollow.ID); err != nil {
		log.Error("error adding followed user", sl.Err(err))

		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

// UnfollowUser removes a user from the follow list.
func (s *Service) UnfollowUser(
	ctx context.Context,
	followLinkID int64,
) error {
	const op = "users.UnfollowUser"

	log := s.log.With(
		slog.String("op", op),
		slog.Int64("follow link ID", followLinkID),
	)

	if err := s.storage.RemoveFollowLink(ctx, followLinkID); err != nil {
		log.Error("error removing follow link", sl.Err(err))

		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
