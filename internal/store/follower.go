package store

import (
	"context"
	"database/sql"
)

type FollowerStore struct {
	db *sql.DB
}

type Follower struct {
	UserID     int64  `json:"user_id"`
	FollowerID int64  `json:"follower_id"`
	CreatedAt  string `json:"created_at"`
}

func (s *FollowerStore) Follow(ctx context.Context, followerID int64, userID int64) error {
	query := `INSERT INTO followers
	(user_id, follower_id)
	VALUES ($1, $2)
	`
	ctx, cancel := context.WithTimeout(ctx, QueryContextTimeout)
	defer cancel()

	_, err := s.db.ExecContext(ctx, query, userID, followerID)
	return err
}

func (s *FollowerStore) Unfollow(ctx context.Context, followID int64, userID int64) error {
	query := `DELETE FROM followers
	WHERE user_id = $1 AND follower_id = $2
	`

	ctx, cancel := context.WithTimeout(ctx, QueryContextTimeout)
	defer cancel()

	_, err := s.db.ExecContext(ctx, query, userID, followID)
	return err
}
