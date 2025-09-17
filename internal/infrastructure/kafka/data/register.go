package data

import (
	"love-signal-users/internal/enum"
	"time"
)

// User is data for user from kafka.
type User struct {
	ID            int64        `json:"id"`
	FullName      string       `json:"full_name"`
	DateOfBirth   *time.Time   `json:"date_of_birth"`
	Gender        *enum.Gender `json:"gender"`
	AvatarFileKey *string      `json:"avatar_file_key"`
}
