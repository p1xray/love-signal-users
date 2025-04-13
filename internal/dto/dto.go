package dto

import "time"

type GenderEnum uint

// Gender enum.
const (
	MALE   = 0
	FEMALE = 1
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
