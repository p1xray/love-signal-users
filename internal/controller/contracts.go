package controller

import (
	"context"
	"love-signal-users/internal/entity"
)

type (
	// UserData is a use-case for getting user data.
	UserData interface {
		Execute(ctx context.Context, id int64) (entity.User, error)
	}

	// UserDataByExternalID is a use-case for getting user data by external ID.
	UserDataByExternalID interface {
		// Execute executes the use-case for getting user data by external ID.
		Execute(ctx context.Context, externalID int64) (entity.User, error)
	}

	// Followed is a use-case for getting followed users.
	Followed interface {
		// Execute executes the use-case for getting followed users.
		Execute(ctx context.Context, userID int64) ([]entity.Follow, error)
	}

	// Follow is a use-case for following users.
	Follow interface {
		// Execute executes the use-case for following user.
		Execute(ctx context.Context, userID int64, userIDToFollow int64) error
	}

	// Unfollow is a use-case for unfollowing users.
	Unfollow interface {
		// Execute executes the use-case for unfollowing user.
		Execute(ctx context.Context, followLinkID int64) error
	}
)
