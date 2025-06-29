package domain

import (
	"github.com/guregu/null/v6"
	"time"
)

// User is data for user in storage.
type User struct {
	ID            int64
	ExternalID    int64
	FullName      string
	DateOfBirth   null.Time
	Gender        null.Int16
	AvatarFileKey null.String
	Deleted       bool
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

// FollowLink is data for follow link in storage.
type FollowLink struct {
	ID              int64
	FollowingUserID int64
	FollowedUserID  int64
	NumberOfLikes   uint32
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

// FollowedUser is data for follow link and followed user in storage.
type FollowedUser struct {
	FollowLink   FollowLink
	FollowedUser User
}
