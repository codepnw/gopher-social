package repository

import (
	"context"
	"database/sql"

	"github.com/codepnw/gopher-social/internal/entity"
)

type CommentRepository interface {
	Create(ctx context.Context, comment *entity.Comment) error
	GetByPostID(ctx context.Context, postID int64) ([]entity.Comment, error)
}

type commentRepository struct {
	db *sql.DB
}

func NewCommentRepository(db *sql.DB) CommentRepository {
	return &commentRepository{db: db}
}

func (r *commentRepository) GetByPostID(ctx context.Context, postID int64) ([]entity.Comment, error) {
	query := `
		SELECT c.id, c.post_id, c.user_id, c.content, c.created_at, users.username, users.id FROM comments c
		JOIN users on users.id = c.user_id
		WHERE c.post_id = $1
		ORDER BY c.created_at DESC;
	`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeout)
	defer cancel()

	rows, err := r.db.QueryContext(ctx, query, postID)
	if err != nil {
		return nil, err
	}

	comments := []entity.Comment{}
	for rows.Next() {
		var c entity.Comment
		c.User = entity.User{}

		err := rows.Scan(
			&c.ID,
			&c.PostID,
			&c.UserID,
			&c.Content,
			&c.CreatedAt,
			&c.User.Username,
			&c.User.ID,
		)
		if err != nil {
			return nil, err
		}

		comments = append(comments, c)
	}
	return comments, nil
}

func (r *commentRepository) Create(ctx context.Context, comment *entity.Comment) error {
	query := `
		INSERT INTO comments (post_id, user_id, content)
		VALUES ($1, $2, $3) RETURNING id, created_at
	`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeout)
	defer cancel()

	err := r.db.QueryRowContext(
		ctx, 
		query,
		comment.PostID,
		comment.UserID,
		comment.Content,
	).Scan(&comment.ID, &comment.CreatedAt)

	if err != nil {
		return err
	}

	return nil
}