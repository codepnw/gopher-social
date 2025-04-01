package store

import (
	"database/sql"

	"github.com/codepnw/gopher-social/internal/repository"
)

type Storage struct {
	Users repository.UserRepository
	Posts repository.PostRepository
}

func NewStorage(db *sql.DB) Storage {
	return Storage{
		Users: repository.NewUserRepository(db),
		Posts: repository.NewPostRepository(db),
	}
}