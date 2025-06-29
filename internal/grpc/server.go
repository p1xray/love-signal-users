package server

import (
	"context"
	"love-signal-users/internal/dto"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// UsersService is a service for working with user data.
type UsersService interface {
	// UserData returns information about a user by their ID.
	UserData(
		ctx context.Context,
		userID int64,
	) (dto.UserData, error)

	// UserDataByExternalID returns information about a user by their external ID.
	UserDataByExternalID(
		ctx context.Context,
		userExternalID int64,
	) (dto.UserData, error)

	// FollowedUsers returns a list of users that the given user is followed to.
	FollowedUsers(
		ctx context.Context,
		userID int64,
	) ([]dto.FollowedUser, error)

	// FollowUser adds the user with userIdToFollow to the list of followed users with userId.
	FollowUser(
		ctx context.Context,
		userID int64,
		userIDToFollow int64,
	) error

	// UnfollowUser removes a user from the follow list.
	UnfollowUser(
		ctx context.Context,
		followLinkID int64,
	) error
}

func InvalidArgumentError(msg string) error {
	return status.Error(codes.InvalidArgument, msg)
}

func InternalError(msg string) error {
	return status.Error(codes.Internal, msg)
}

func NotFoundError(msg string) error {
	return status.Error(codes.NotFound, msg)
}
