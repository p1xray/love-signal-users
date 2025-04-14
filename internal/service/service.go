package service

import (
	"context"
	"fmt"
	"log/slog"
	"love-signal-users/internal/dto"
	"love-signal-users/internal/lib/logger/sl"
	"love-signal-users/internal/storage"
)

type UsersService struct {
	log     *slog.Logger
	storage storage.UserStorage
}

// New creates a new user service.
func New(log *slog.Logger, storage storage.UserStorage) *UsersService {
	return &UsersService{
		log:     log,
		storage: storage,
	}
}

// UserInfo returns information about a user by their external identifier.
func (us *UsersService) UserInfo(
	ctx context.Context,
	userExternalId int64,
) (*dto.UserInfo, error) {
	const op = "service.UserInfo"

	log := us.log.With(
		slog.String("op", op),
		slog.Int64("user external id", userExternalId),
	)

	userInfo, err := us.storage.UserInfoByExternalId(ctx, userExternalId)
	if err != nil {
		log.Error("failed to get user info by user external id", sl.Err(err))

		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return userInfo, nil
}

// UserProfileCard returns the user profile card.
func (us *UsersService) UserProfileCard(
	ctx context.Context,
	userId int64,
) (*dto.UserProfileCard, error) {
	const op = "service.UserProfileCard"

	log := us.log.With(
		slog.String("op", op),
		slog.Int64("user id", userId),
	)

	card, err := us.storage.UserProfileCard(ctx, userId)
	if err != nil {
		log.Error("failed to get user profile card data", sl.Err(err))

		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return card, nil
}

// FollowedUsers returns a list of users that the given user is followed to.
func (us *UsersService) FollowedUsers(
	ctx context.Context,
	userId int64,
) ([]*dto.FollowedUser, error) {
	const op = "service.FollowedUsers"

	log := us.log.With(
		slog.String("op", op),
		slog.Int64("user id", userId),
	)

	followedUsers, err := us.storage.FollowedUsersByUserId(ctx, userId)
	if err != nil {
		log.Error("failed to get followed users by user id", sl.Err(err))

		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return followedUsers, nil
}

// FollowUser adds the user with userIdToFollow to the list of followed users with userId.
func (us *UsersService) FollowUser(
	ctx context.Context,
	userId int64,
	userIdToFollow int64,
) error {
	const op = "service.FollowUser"

	log := us.log.With(
		slog.String("op", op),
		slog.Int64("user id", userId),
		slog.Int64("user id to follow", userIdToFollow),
	)

	user, err := us.storage.UserInfoById(ctx, userId)
	if err != nil {
		log.Error("failed to get user by user id", sl.Err(err))

		return fmt.Errorf("%s: %w", op, err)
	}

	userToFollow, err := us.storage.UserInfoById(ctx, userIdToFollow)
	if err != nil {
		log.Error("failed to get user to follow by user id", sl.Err(err))

		return fmt.Errorf("%s: %w", op, err)
	}

	if err := us.storage.AddFollowLink(ctx, user.Id, userToFollow.Id); err != nil {
		log.Error("failed to added followed user", sl.Err(err))

		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

// UnfollowUser removes a user from the follow list.
func (us *UsersService) UnfollowUser(
	ctx context.Context,
	followLinkId int64,
) error {
	const op = "service.UnfollowUser"

	log := us.log.With(
		slog.String("op", op),
		slog.Int64("follow link id", followLinkId),
	)

	if err := us.storage.RemoveFollowLink(ctx, followLinkId); err != nil {
		log.Error("failed to remove follow link", sl.Err(err))

		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
