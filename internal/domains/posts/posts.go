package posts

import (
	"database/sql"
)

func InitPostDomain(db *sql.DB) PostHandler {
	postrepo := NewPostRepository(db)
	postusecase := NewPostUsecase(postrepo)
	posthandler := NewPostHandler(postusecase)

	return posthandler
}
