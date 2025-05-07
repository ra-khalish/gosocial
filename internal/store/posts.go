package store

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Post struct {
	ID        int64     `json:"id"`
	Content   string    `json:"content"`
	Title     string    `json:"title"`
	UserID    int64     `json:"user_id"`
	Tags      []string  `json:"tags"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Version   int       `json:"version"`
	Comments  []Comment `json:"comments"`
}

type PostStore struct {
	db *pgxpool.Pool
}

func (s *PostStore) Create(ctx context.Context, post *Post) error {
	query := `
	INSERT INTO posts (content, title, user_id, tags)
	VALUES ($1, $2, $3, $4) RETURNING id, created_at, updated_at
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	err := s.db.QueryRow(
		ctx,
		query,
		post.Content,
		post.Title,
		post.UserID,
		post.Tags,
	).Scan(
		&post.ID,
		&post.CreatedAt,
		&post.UpdatedAt,
	)
	if err != nil {
		return err
	}

	return nil
}

func (s *PostStore) GetByID(ctx context.Context, id int64) (*Post, error) {
	query := `
		SELECT id, user_id, title, content, created_at, updated_at, tags, version
		FROM posts
		WHERE id = $1
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	var post Post
	err := s.db.QueryRow(ctx, query, id).Scan(
		&post.ID,
		&post.UserID,
		&post.Title,
		&post.Content,
		&post.CreatedAt,
		&post.UpdatedAt,
		&post.Tags,
		&post.Version,
	)
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}
	return &post, nil
}

func (s *PostStore) Delete(ctx context.Context, postID int64) (int64, error) {
	query := `
	DELETE FROM posts WHERE id = $1
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	res, err := s.db.Exec(ctx, query, postID)

	if err != nil {
		return 0, err
	}

	rows := res.RowsAffected()

	if rows == 0 {
		return 0, ErrNotFound
	}
	return rows, nil
}

func (s *PostStore) Update(ctx context.Context, post *Post) error {
	query := `
	UPDATE posts
	SET title = $1, content = $2, version = version + 1
	WHERE id = $3 AND version = $4
	RETURNING version
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	err := s.db.QueryRow(ctx, query, post.Title, post.Content, post.ID, post.Version).Scan(&post.Version)
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return ErrNotFound
		default:
			return err
		}
	}
	// _, err := s.db.Exec(ctx, query, post.Title, post.Content, post.ID, post.Version)
	// if err != nil {
	// 	return err
	// }
	return nil
}
