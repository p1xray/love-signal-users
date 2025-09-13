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
)
