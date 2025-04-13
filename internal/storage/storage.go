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

	// UserProfileCard returns the user profile card.
	UserProfileCard(
		ctx context.Context,
		userId int64,
	) (*dto.UserProfileCard, error)
}
