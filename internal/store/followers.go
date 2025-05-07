package store

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Follower struct {
	UserID     int64     `json:"user_id"`
	FollowerID int64     `json:"follower_id"`
	CreatedAt  time.Time `json:"created_at"`
}

type FollowerStore struct {
	db *pgxpool.Pool
}

func (s *FollowerStore) GetFollowers(ctx context.Context, userID int64) error {
	err := errors.New("error get followers")
	return err
}

func (s *FollowerStore) GetFollowing(ctx context.Context, userID int64) error {
	err := errors.New("error get following")
	return err
}

func (s *FollowerStore) Follow(ctx context.Context, followerID, userID int64) error {
	query := `
		INSERT INTO followers (user_id, follower_id) VALUES ($1, $2)
	`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()
	_, err := s.db.Exec(ctx, query, userID, followerID)
	return err
}

func (s *FollowerStore) Unfollow(ctx context.Context, followerID, userID int64) error {
	query := `
		DELETE followers WHERE user_id = $1 AND follower_id = $2
	`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()
	_, err := s.db.Exec(ctx, query, userID, followerID)
	return err
}
