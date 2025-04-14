package storage

import (
	"context"
	"errors"
	"love-signal-users/internal/dto"
)

var (
	ErrFollowExist = errors.New("follow is already exist")
)

type UserStorage interface {
	// UserInfoByExternalId returns information about a user by their external identifier.
	UserInfoByExternalId(
		ctx context.Context,
		userExternalId int64,
	) (*dto.UserInfo, error)

	// UserInfoById returns information about a user by their identifier.
	UserInfoById(
		ctx context.Context,
		userId int64,
	) (*dto.UserInfo, error)

	// UserProfileCard returns the user profile card.
	UserProfileCard(
		ctx context.Context,
		userId int64,
	) (*dto.UserProfileCard, error)

	// FollowedUsersByUserId returns a list of users that the given user is followed to.
	FollowedUsersByUserId(
		ctx context.Context,
		userId int64,
	) ([]*dto.FollowedUser, error)

	// AddFollowLink adds a follow link with the given user IDs.
	AddFollowLink(
		ctx context.Context,
		userId int64,
		userIdToFollow int64,
	) error

	// RemoveFollowLink removes a follow link by followLinkId.
	RemoveFollowLink(
		ctx context.Context,
		followLinkId int64,
	) error
}
