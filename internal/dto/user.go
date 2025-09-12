package dto

import (
	"love-signal-users/internal/enum"
	"time"
)

// User is a DTO with user data.
type User struct {
	ID            int64
	ExternalID    int64
	FullName      string
	Gender        *enum.Gender
	DateOfBirth   *time.Time
	AvatarFileKey *string
}
