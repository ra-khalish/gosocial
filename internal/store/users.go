package store

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type User struct {
	ID        int64     `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Password  string    `json:"-"`
	CreatedAt time.Time `json:"created_at"`
}

type UsersStore struct {
	db *pgxpool.Pool
}

func (s *UsersStore) Create(ctx context.Context, user *User) error {
	query := `
	INSERT INTO users (username, password, email) VALUES($1, $2, $3) 
	RETURNING id, created_at
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	err := s.db.QueryRow(
		ctx,
		query,
		user.Username,
		user.Password,
		user.Email,
	).Scan(
		&user.ID,
		&user.CreatedAt,
	)
	if err != nil {
		return err
	}
	return nil
}

func (s *UsersStore) GetByID(ctx context.Context, userID int64) (*User, error) {
	query := `
	SELECT id, username, email, password, created_at FROM users WHERE id = $1
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	var user User
	var password []byte

	err := s.db.QueryRow(
		ctx,
		query,
		userID,
	).Scan(
		&user.ID,
		&user.Username, &user.Email, &password, &user.CreatedAt)

	if err != nil {
		switch err {
		case pgx.ErrNoRows:
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}
	user.Password = string(password)

	return &user, nil
}
