package store

import (
	"database/sql"

	"github.com/codepnw/gopher-social/internal/handler"
	"github.com/codepnw/gopher-social/internal/repository"
)

type Storage struct {
	Posts handler.PostsHandler
	Users handler.UserHandler
}

func NewStorage(db *sql.DB) Storage {
	postRepo := repository.NewPostRepository(db)
	commentRepo := repository.NewCommentRepository(db)
	userRepo := repository.NewUserRepository(db)

	return Storage{
		Posts: handler.NewPostsHandler(postRepo, commentRepo),
		Users: handler.NewUserHandler(userRepo),
	}
}
