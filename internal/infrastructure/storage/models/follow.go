package models

import "time"

// Follow is data for follow link in storage.
type Follow struct {
	ID              int64
	FollowingUserID int64
	FollowedUserID  int64
	NumberOfLikes   uint32
	CreatedAt       time.Time
	UpdatedAt       time.Time
}
