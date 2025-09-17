package kafka

import "context"

type (
	// UserHasRegisteredHandler is a register new user handler.
	UserHasRegisteredHandler interface {
		// Execute  creates a new user in the storage based on the received data.
		Execute(ctx context.Context, dataAsBytes []byte) error
	}
)
