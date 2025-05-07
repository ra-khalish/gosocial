package store

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrNotFound          = errors.New("resource not found")
	QueryTimeoutDuration = time.Second * 5
)

type Storage struct {
	Posts interface {
		GetByID(context.Context, int64) (*Post, error)
		Create(context.Context, *Post) error
		Delete(context.Context, int64) (int64, error)
		Update(context.Context, *Post) error
	}
	Users interface {
		Create(context.Context, *User) error
		GetByID(context.Context, int64) (*User, error)
	}
	Comment interface {
		Create(context.Context, *Comment) error
		GetByPostID(context.Context, int64) ([]Comment, error)
	}
	Followers interface {
		Follow(ctx context.Context, followerId, userID int64) error
		Unfollow(ctx context.Context, followerId, userID int64) error
	}
}

func NewStorage(db *pgxpool.Pool) Storage {
	return Storage{
		Posts:     &PostStore{db},
		Users:     &UsersStore{db},
		Comment:   &CommentStore{db},
		Followers: &FollowerStore{db},
	}
}
