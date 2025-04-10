package dto

import "time"

// UserInfo is information about the user.
type UserInfo struct {
	Id   int64
	Name string
}

// UserProfileCard is user profile card.
type UserProfileCard struct {
	Id          int64
	Name        string
	DateOfBirth time.Time
	Gender      int
}
