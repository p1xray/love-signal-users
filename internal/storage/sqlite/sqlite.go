package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"love-signal-users/internal/dto"

	_ "github.com/mattn/go-sqlite3"
)

type Storage struct {
	db *sql.DB
}

// New creates a new SQLite storage.
func New(storagePath string) (*Storage, error) {
	const op = "storage.sqlite.New"

	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db: db}, nil
}

// UserInfoByExternalId returns information about a user by their external identifier.
func (s *Storage) UserInfoByExternalId(
	ctx context.Context,
	userExternalId int64,
) (*dto.UserInfo, error) {
	const op = "sqlite.UserByExternalId"

	stmt, err := s.db.PrepareContext(ctx,
		"select u.id, u.name from users u where u.deleted = false and u.external_id = ?;")

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	row := stmt.QueryRowContext(ctx, userExternalId)

	var user dto.UserInfo
	err = row.Scan(&user.Id, &user.Name)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &user, nil
}

// UserProfileCard returns the user profile card.
func (s *Storage) UserProfileCard(
	ctx context.Context,
	userId int64,
) (*dto.UserProfileCard, error) {
	const op = "sqlite.UserProfileCard"

	stmt, err := s.db.PrepareContext(ctx,
		`select u.id, u.name, u.date_of_birth, u.gender, u.avatar_file__key
		from users u where u.deleted = false and u.id = ?;`)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	row := stmt.QueryRowContext(ctx, userId)

	var card dto.UserProfileCard
	err = row.Scan(&card.Id, &card.Name, &card.DateOfBirth, &card.Gender, &card.AvatarFileKey)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &card, nil
}
