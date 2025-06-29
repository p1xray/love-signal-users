package dto

import "time"

type GenderEnum uint

// Gender enum.
const (
	MALE   = 1
	FEMALE = 2
)

// UserData is information about the user.
type UserData struct {
	ID            int64
	FullName      string
	DateOfBirth   *time.Time
	Gender        *GenderEnum
	AvatarFileKey *string
}

// FollowedUser is user that is followed to.
type FollowedUser struct {
	FollowLinkID  int64
	NumberOfLikes uint32
	UserID        int64
	FullName      string
	AvatarFileKey *string
}
