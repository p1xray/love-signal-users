package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"love-signal-users/internal/dto"
	"love-signal-users/internal/storage"
	"time"

	"github.com/mattn/go-sqlite3"
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

// FollowedUsersByUserId returns a list of users that the given user is followed to.
func (s *Storage) FollowedUsersByUserId(
	ctx context.Context,
	userId int64,
) ([]*dto.FollowedUser, error) {
	const op = "sqlite.FollowedUsersByUserId"

	stmt, err := s.db.PrepareContext(ctx,
		`select f.id, f.sended_likes_count, followed_user.id, followed_user.name, followed_user.avatar_file__key
		from follows f
			join users user on f.following_user_id = user.id
			join users followed_user on f.followed_user_id = followed_user.id
		where f.deleted = false and user.deleted = false and followed_user.deleted = false and user.id = ?;`)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	rows, err := stmt.QueryContext(ctx, userId)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	followedUsers := make([]*dto.FollowedUser, 0)
	for rows.Next() {
		fu := &dto.FollowedUser{}
		err := rows.Scan(&fu.FollowLinkId, &fu.SendedLikesCount, &fu.UserId, &fu.Name, &fu.AvatarFileKey)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		followedUsers = append(followedUsers, fu)
	}

	return followedUsers, nil
}

// AddFollowLink adds a follow link with the given user IDs.
func (s *Storage) AddFollowLink(
	ctx context.Context,
	userId int64,
	userIdToFollow int64,
) error {
	const op = "sqlite.AddFollowedUser"

	stmt, err := s.db.PrepareContext(ctx,
		`insert into follows (following_user_id, followed_user_id, sended_likes_count, deleted, created_at, updated_at)
		values (?, ?, ?, ?, ?, ?);`)

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	now := time.Now()
	_, err = stmt.ExecContext(ctx, userId, userIdToFollow, 0, false, now, now)
	if err != nil {
		var sqliteErr sqlite3.Error
		if errors.As(err, &sqliteErr) && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
			return fmt.Errorf("%s: %w", op, storage.ErrFollowExist)
		}

		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

// RemoveFollowLink removes a follow link by followLinkId.
func (s *Storage) RemoveFollowLink(
	ctx context.Context,
	followLinkId int64,
) error {
	const op = "sqlite.RemoveFollowLink"

	stmt, err := s.db.PrepareContext(ctx,
		`update follows
		set deleted = true
		where id = ?;`)

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = stmt.ExecContext(ctx, followLinkId)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
