package service

import (
	"context"
	"errors"
	"love-signal-users/internal/storage/domain"
)

var (
	ErrUserNotFound = errors.New("user not found")
)

type UserStorage interface {
	// UserData returns information about a user by their ID.
	UserData(
		ctx context.Context,
		userID int64,
	) (domain.User, error)

	// UserDataByExternalID returns information about a user by their external ID.
	UserDataByExternalID(
		ctx context.Context,
		userExternalID int64,
	) (domain.User, error)

	// FollowedUsers returns a list of users that the given user is followed to.
	FollowedUsers(
		ctx context.Context,
		userID int64,
	) ([]domain.FollowedUser, error)

	// AddFollowLink adds a follow link with the given user IDs.
	AddFollowLink(
		ctx context.Context,
		userID int64,
		userIDToFollow int64,
	) error

	// RemoveFollowLink removes a follow link by followLinkId.
	RemoveFollowLink(
		ctx context.Context,
		followLinkID int64,
	) error
}
