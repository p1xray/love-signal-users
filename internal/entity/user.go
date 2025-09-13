package entity

import (
	"love-signal-users/internal/dto"
	"love-signal-users/internal/enum"
	"time"
)

// User is the user entity.
type User struct {
	ID            int64
	ExternalID    int64
	FullName      string
	Gender        *enum.Gender
	DateOfBirth   *time.Time
	AvatarFileKey *string

	dataStatus enum.DataStatus
}

// NewUser returns new user entity.
func NewUser(data dto.User) User {
	return User{
		ID:            data.ID,
		ExternalID:    data.ExternalID,
		FullName:      data.FullName,
		Gender:        data.Gender,
		DateOfBirth:   data.DateOfBirth,
		AvatarFileKey: data.AvatarFileKey,
	}
}

func (u *User) SetToCreate() {
	u.dataStatus = enum.ToCreate
}

func (u *User) SetToUpdate() {
	u.dataStatus = enum.ToUpdate
}

func (u *User) SetToRemove() {
	u.dataStatus = enum.ToRemove
}

func (u *User) IsToCreate() bool {
	return u.dataStatus == enum.ToCreate
}

func (u *User) IsToUpdate() bool {
	return u.dataStatus == enum.ToUpdate
}

func (u *User) IsToRemove() bool {
	return u.dataStatus == enum.ToRemove
}

func (u *User) ResetDataStatus() {
	u.dataStatus = enum.None
}
