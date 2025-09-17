package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"love-signal-users/internal/infrastructure/converter"
	"love-signal-users/internal/infrastructure/kafka/data"
	"love-signal-users/internal/infrastructure/storage/models"
	"love-signal-users/pkg/logger/sl"
)

// UserHasRegisteredTopic is the name of the topic in kafka where data for the registered user handler is stored.
const UserHasRegisteredTopic = "test-register"

type UserHasRegisteredStorage interface {
	CreateUser(ctx context.Context, user models.User) (int64, error)
}

// UserHasRegistered is a register new user handler.
type UserHasRegistered struct {
	log     *slog.Logger
	storage UserHasRegisteredStorage
}

// NewUserHasRegistered returns new instance of UserHasRegistered handler.
func NewUserHasRegistered(log *slog.Logger, storage UserHasRegisteredStorage) *UserHasRegistered {
	return &UserHasRegistered{
		log:     log,
		storage: storage,
	}
}

// Execute  creates a new user in the storage based on the received data.
func (uhr *UserHasRegistered) Execute(ctx context.Context, dataAsBytes []byte) error {
	const op = "handlers.register.Execute"

	log := uhr.log.With(
		slog.String("op", op),
	)

	var user data.User
	if err := json.Unmarshal(dataAsBytes, &user); err != nil {
		log.Error("error unmarshalling user data", sl.Err(err))
	}

	userStorage := converter.ToUserStorage(user, models.UserCreated())
	if _, err := uhr.storage.CreateUser(ctx, userStorage); err != nil {
		log.Error("error creating user", sl.Err(err))

		return fmt.Errorf("%s:%w", op, err)
	}

	return nil
}
