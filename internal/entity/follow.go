package entity

import (
	"love-signal-users/internal/dto"
	"love-signal-users/internal/enum"
)

// Follow is the follow entity.
type Follow struct {
	ID            int64
	FollowingUser User
	FollowedUser  User
	NumberOfLikes uint32

	dataStatus enum.DataStatus
}

// NewFollow returns new follow entity.
func NewFollow(data dto.Follow) Follow {
	return Follow{
		ID:            data.ID,
		FollowingUser: NewUser(data.FollowingUser),
		FollowedUser:  NewUser(data.FollowedUser),
		NumberOfLikes: data.NumberOfLikes,
	}
}

func (f *Follow) SetToCreate() {
	f.dataStatus = enum.ToCreate
}

func (f *Follow) SetToUpdate() {
	f.dataStatus = enum.ToUpdate
}

func (f *Follow) SetToRemove() {
	f.dataStatus = enum.ToRemove
}

func (f *Follow) IsToCreate() bool {
	return f.dataStatus == enum.ToCreate
}

func (f *Follow) IsToUpdate() bool {
	return f.dataStatus == enum.ToUpdate
}

func (f *Follow) IsToRemove() bool {
	return f.dataStatus == enum.ToRemove
}

func (f *Follow) ResetDataStatus() {
	f.dataStatus = enum.None
}
