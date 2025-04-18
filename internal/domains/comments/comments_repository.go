package comments

import (
	"context"
	"database/sql"
)

type CommentsRepository interface {
	create(ctx context.Context, comment *Comment) error
	getByPostID(ctx context.Context, postID int64) ([]Comment, error)
}

type repository struct {
	db *sql.DB
}

func NewCommentsRepository(db *sql.DB) CommentsRepository {
	return &repository{db: db}
}

func (r *repository) create(ctx context.Context, comment *Comment) error {
	query := `
		INSERT INTO comments (post_id, user_id, content)
		VALUES ($1, $2, $3) RETURNING id, created_at
	`
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

func (r *repository) getByPostID(ctx context.Context, postID int64) ([]Comment, error) {
	query := `
		SELECT c.id, c.post_id, c.user_id, c.content, c.created_at, users.username, users.id FROM comments c
		JOIN users on users.id = c.user_id
		WHERE c.post_id = $1
		ORDER BY c.created_at DESC;
	`
	_ = query
	

	return nil, nil
}

