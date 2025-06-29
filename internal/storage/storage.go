package storage

import (
	"errors"
)

var (
	ErrEntityNotFound = errors.New("entity not found")
	ErrFollowExist    = errors.New("follow is already exist")
)
