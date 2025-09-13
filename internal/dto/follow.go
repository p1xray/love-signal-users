package dto

// Follow is a DTO with follow data.
type Follow struct {
	ID            int64
	FollowingUser User
	FollowedUser  User
	NumberOfLikes uint32
}
