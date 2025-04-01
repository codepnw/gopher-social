package repository

import (
	"context"
	"database/sql"

	"github.com/codepnw/gopher-social/internal/entity"
	"github.com/lib/pq"
)

// type Storage struct {
// 	Posts interface {
// 		Create(context.Context) error
// 	}
// 	Users interface {
// 		Create(context.Context) error
// 	}
// }

type PostRepository interface {
	Create(ctx context.Context, post *entity.Post) error
}

type postRepository struct {
	db *sql.DB
}

func NewPostRepository(db *sql.DB) PostRepository {
	return &postRepository{db: db}
}

func (r *postRepository) Create(ctx context.Context, post *entity.Post) error {
	query := `
		INSERT INTO posts (title, content, user_id, tags)
		VALUES ($1, $2, $3, $4) RETURNING id, created_at, updated_at
	`
	err := r.db.QueryRowContext(
		ctx, 
		query, 
		post.Title, 
		post.Content, 
		post.UserID, 
		pq.Array(post.Tags),
	).Scan(&post.ID, &post.CreatedAt, &post.UpdatedAt)
	
	if err != nil {
		return err
	}

	return nil
}
