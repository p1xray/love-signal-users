package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/mattn/go-sqlite3"
	_ "github.com/mattn/go-sqlite3"
	"love-signal-users/internal/infrastructure"
	"love-signal-users/internal/infrastructure/storage/models"
	"strings"
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

// Users returns slice of users by ids from storage.
func (s *Storage) Users(ctx context.Context, ids []int64) ([]models.User, error) {
	const op = "sqlite.Users"

	// Generate the placeholders for the IN clause.
	placeholders := make([]string, len(ids))
	for i := range ids {
		placeholders[i] = "?"
	}
	inClause := strings.Join(placeholders, ",")

	query := fmt.Sprintf(
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
		where u.deleted = false and u.id in (%s);`, inClause)

	stmt, err := s.db.PrepareContext(ctx, query)
	if err != nil {
		return []models.User{}, fmt.Errorf("%s: %w", op, err)
	}

	args := make([]interface{}, len(ids))
	for i, rc := range ids {
		args[i] = rc
	}

	rows, err := stmt.QueryContext(ctx, args...)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	users := make([]models.User, 0)
	for rows.Next() {
		user := models.User{}
		err = rows.Scan(
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
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		users = append(users, user)
	}

	return users, nil
}

// User returns information about a user by their ID from storage.
func (s *Storage) User(ctx context.Context, userID int64) (models.User, error) {
	const op = "sqlite.User"

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
		return models.User{}, fmt.Errorf("%s: %w", op, err)
	}

	row := stmt.QueryRowContext(ctx, userID)

	var user models.User
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
			return models.User{}, fmt.Errorf("%s: %w", op, infrastructure.ErrEntityNotFound)
		}

		return models.User{}, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}

// UserByExternalID returns information about a user by their external ID from storage.
func (s *Storage) UserByExternalID(ctx context.Context, externalID int64) (models.User, error) {
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
		return models.User{}, fmt.Errorf("%s: %w", op, err)
	}

	row := stmt.QueryRowContext(ctx, externalID)

	var user models.User
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
			return models.User{}, fmt.Errorf("%s: %w", op, infrastructure.ErrEntityNotFound)
		}

		return models.User{}, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}

// FollowsByUserID returns a list of follow links that the given user is followed to from storage.
func (s *Storage) FollowsByUserID(
	ctx context.Context,
	userID int64,
) ([]models.Follow, error) {
	const op = "sqlite.FollowedUsers"

	stmt, err := s.db.PrepareContext(ctx,
		`select
    	f.id,
    	f.following_user_id,
    	f.followed_user_id,
    	f.number_of_likes,
    	f.created_at,
    	f.updated_at
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

	follows := make([]models.Follow, 0)
	for rows.Next() {
		follow := models.Follow{}
		err = rows.Scan(
			&follow.ID,
			&follow.FollowingUserID,
			&follow.FollowedUserID,
			&follow.NumberOfLikes,
			&follow.CreatedAt,
			&follow.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		follows = append(follows, follow)
	}

	return follows, nil
}

// CreateFollow creates the follow link in storage.
func (s *Storage) CreateFollow(ctx context.Context, follow models.Follow) (int64, error) {
	const op = "sqlite.CreateFollow"

	stmt, err := s.db.PrepareContext(ctx,
		`insert into follows (following_user_id, followed_user_id, number_of_likes, created_at, updated_at)
		values (?, ?, ?, ?);`)

	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	res, err := stmt.ExecContext(
		ctx,
		follow.FollowingUserID,
		follow.FollowedUserID,
		follow.NumberOfLikes,
		follow.CreatedAt,
		follow.UpdatedAt,
	)
	if err != nil {
		var sqliteErr sqlite3.Error
		if errors.As(err, &sqliteErr) && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
			return 0, fmt.Errorf("%s: %w", op, infrastructure.ErrFollowExist)
		}

		return 0, fmt.Errorf("%s: %w", op, err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

// UpdateFollow updates the follow link in storage.
func (s *Storage) UpdateFollow(ctx context.Context, follow models.Follow) error {
	const op = "sqlite.UpdateFollow"

	stmt, err := s.db.PrepareContext(ctx,
		`update follows
		 set following_user_id = ?,
			 followed_user_id = ?,
			 number_of_likes = ?,
			 updated_at = ?
		 where id = ?;`)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = stmt.ExecContext(
		ctx,
		follow.FollowingUserID,
		follow.FollowedUserID,
		follow.NumberOfLikes,
		follow.UpdatedAt,
		follow.ID,
	)

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

// RemoveFollow removes the follow link from storage by followLinkID.
func (s *Storage) RemoveFollow(ctx context.Context, followLinkID int64) error {
	const op = "sqlite.RemoveFollow"

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
