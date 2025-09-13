package models

import "time"

type FollowOption func(*Follow)

func FollowCreated() FollowOption {
	now := time.Now()
	return func(f *Follow) {
		f.CreatedAt = now
		f.UpdatedAt = now
	}
}

func FollowUpdated() FollowOption {
	return func(f *Follow) {
		f.UpdatedAt = time.Now()
	}
}
