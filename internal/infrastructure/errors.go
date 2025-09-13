package infrastructure

import "errors"

var (
	ErrEntityNotFound    = errors.New("entity not found")
	ErrRequireIDToUpdate = errors.New("a non-null identifier is required to update an entity in storage")
	ErrRequireIDToRemove = errors.New("a non-null identifier is required to remove an entity in storage")
	ErrFollowExist       = errors.New("follow already exists")
)
