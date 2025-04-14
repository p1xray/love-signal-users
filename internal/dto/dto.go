package dto

import "time"

type GenderEnum uint

// Gender enum.
const (
	MALE   = 1
	FEMALE = 2
)

// UserInfo is information about the user.
type UserInfo struct {
	Id   int64
	Name string
}

// UserProfileCard is user profile card.
type UserProfileCard struct {
	Id            int64
	Name          string
	DateOfBirth   *time.Time
	Gender        *GenderEnum
	AvatarFileKey *string
}

// FollowedUser is user that is followed to.
type FollowedUser struct {
	FollowLinkId     int64
	SendedLikesCount uint32
	UserId           int64
	Name             string
	AvatarFileKey    *string
}
