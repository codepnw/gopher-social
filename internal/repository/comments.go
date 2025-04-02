package repository

import (
	"context"
	"database/sql"

	"github.com/codepnw/gopher-social/internal/entity"
)

type CommentRepository interface {
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
