package store

import (
	"database/sql"

	"github.com/codepnw/gopher-social/internal/handler"
	"github.com/codepnw/gopher-social/internal/repository"
)

type Storage struct {
	Posts handler.PostsHandler
}

func NewStorage(db *sql.DB) Storage {
	postRepo := repository.NewPostRepository(db)
	commentRepo := repository.NewCommentRepository(db)

	return Storage{
		Posts: handler.NewPostsHandler(postRepo, commentRepo),
	}
}
