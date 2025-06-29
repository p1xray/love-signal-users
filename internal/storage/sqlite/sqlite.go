package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"love-signal-users/internal/storage"
	"love-signal-users/internal/storage/domain"
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

// UserData returns information about a user by their ID.
func (s *Storage) UserData(
	ctx context.Context,
	userID int64,
) (domain.User, error) {
	const op = "sqlite.UserData"

	stmt, err := s.db.PrepareContext(ctx,
		`select
    	u.id,
    	u.external_id,
    	u.full_name,
    	u.date_of_birth,
    	u.gender,
    	u.avatar_file_key,
    	u.deleted,
    	u.created_at,
    	u.updated_at
		from users u
		where u.deleted = false and u.id = ?;`)

	if err != nil {
		return domain.User{}, fmt.Errorf("%s: %w", op, err)
	}

	row := stmt.QueryRowContext(ctx, userID)

	var user domain.User
	err = row.Scan(
		&user.ID,
		&user.ExternalID,
		&user.FullName,
		&user.DateOfBirth,
		&user.Gender,
		&user.AvatarFileKey,
		&user.Deleted,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.User{}, fmt.Errorf("%s: %w", op, storage.ErrEntityNotFound)
		}

		return domain.User{}, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}

// UserDataByExternalID returns information about a user by their external ID.
func (s *Storage) UserDataByExternalID(
	ctx context.Context,
	userExternalID int64,
) (domain.User, error) {
	const op = "sqlite.UserDataByExternalID"

	stmt, err := s.db.PrepareContext(ctx,
		`select
    	u.id,
    	u.external_id,
    	u.full_name,
    	u.date_of_birth,
    	u.gender,
    	u.avatar_file_key,
    	u.deleted,
    	u.created_at,
    	u.updated_at
		from users u
		where u.deleted = false and u.id = ?;`)

	if err != nil {
		return domain.User{}, fmt.Errorf("%s: %w", op, err)
	}

	row := stmt.QueryRowContext(ctx, userExternalID)

	var user domain.User
	err = row.Scan(
		&user.ID,
		&user.ExternalID,
		&user.FullName,
		&user.DateOfBirth,
		&user.Gender,
		&user.AvatarFileKey,
		&user.Deleted,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.User{}, fmt.Errorf("%s: %w", op, storage.ErrEntityNotFound)
		}

		return domain.User{}, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}

// FollowedUsers returns a list of users that the given user is followed to.
func (s *Storage) FollowedUsers(
	ctx context.Context,
	userID int64,
) ([]domain.FollowedUser, error) {
	const op = "sqlite.FollowedUsers"

	stmt, err := s.db.PrepareContext(ctx,
		`select
    	f.id,
    	f.following_user_id,
    	f.followed_user_id,
    	f.number_of_likes,
    	f.created_at,
    	f.updated_at,
    	followed_user.id,
    	followed_user.external_id,
    	followed_user.full_name,
    	followed_user.date_of_birth,
    	followed_user.gender,
    	followed_user.avatar_file_key,
    	followed_user.deleted,
    	followed_user.created_at,
    	followed_user.updated_at
		from follows f
			join users user on f.following_user_id = user.id
			join users followed_user on f.followed_user_id = followed_user.id
		where user.deleted = false and followed_user.deleted = false and f.following_user_id = ?;`)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	rows, err := stmt.QueryContext(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	followedUsers := make([]domain.FollowedUser, 0)
	for rows.Next() {
		fu := domain.FollowedUser{}
		err = rows.Scan(
			&fu.FollowLink.ID,
			&fu.FollowLink.FollowingUserID,
			&fu.FollowLink.FollowedUserID,
			&fu.FollowLink.NumberOfLikes,
			&fu.FollowLink.CreatedAt,
			&fu.FollowLink.UpdatedAt,
			&fu.FollowedUser.ID,
			&fu.FollowedUser.ExternalID,
			&fu.FollowedUser.FullName,
			&fu.FollowedUser.DateOfBirth,
			&fu.FollowedUser.Gender,
			&fu.FollowedUser.AvatarFileKey,
			&fu.FollowedUser.Deleted,
			&fu.FollowedUser.CreatedAt,
			&fu.FollowedUser.UpdatedAt)
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
	userID int64,
	userIDToFollow int64,
) error {
	const op = "sqlite.AddFollowLink"

	stmt, err := s.db.PrepareContext(ctx,
		`insert into follows (following_user_id, followed_user_id, created_at, updated_at)
		values (?, ?, ?, ?);`)

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	now := time.Now()
	_, err = stmt.ExecContext(ctx, userID, userIDToFollow, now, now)
	if err != nil {
		var sqliteErr sqlite3.Error
		if errors.As(err, &sqliteErr) && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
			return fmt.Errorf("%s: %w", op, storage.ErrFollowExist)
		}

		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

// RemoveFollowLink removes a follow link by followLinkID.
func (s *Storage) RemoveFollowLink(
	ctx context.Context,
	followLinkID int64,
) error {
	const op = "sqlite.RemoveFollowLink"

	stmt, err := s.db.PrepareContext(ctx, "delete from follows where id = ?;")

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = stmt.ExecContext(ctx, followLinkID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
