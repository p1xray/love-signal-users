package service

import (
	"context"
	"fmt"
	"log/slog"
	"love-signal-users/internal/dto"
	"love-signal-users/internal/lib/logger/sl"
)

type UsersService struct {
	log     *slog.Logger
	storage UserStorage
}

type UserStorage interface {
	// UserInfoByExternalId returns information about a user by their external identifier.
	UserInfoByExternalId(
		ctx context.Context,
		userExternalId int64,
	) (*dto.UserInfo, error)

	// UserProfileCard returns the user profile card.
	UserProfileCard(
		ctx context.Context,
		userId int64,
	) (*dto.UserProfileCard, error)
}

// New creates a new user service.
func New(log *slog.Logger) *UsersService {
	return &UsersService{
		log: log,
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
