package store

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

var (
	ErrNotFound         = errors.New("record not found")
	QueryContextTimeout = time.Second * 5
)

type Storage struct {
	Posts interface {
		Create(context.Context, *Post) error
		Get(context.Context, int64) (*Post, error)
		Update(context.Context, *Post) error
		Delete(context.Context, *Post) error
	}
	Users interface {
		Create(context.Context, *User) error
		GetUserByID(context.Context, int64) (*User, error)
	}
	Comments interface {
		Create(context.Context, *Comment) error
		GetByPostID(ctx context.Context, postID int64) ([]Comment, error)
	}
}

func NewStorage(db *sql.DB) Storage {
	return Storage{
		Posts:    &PostStore{db: db},
		Users:    &UserStore{db: db},
		Comments: &CommentStore{db: db},
	}
}
