package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/codepnw/gopher-social/internal/entity"
	"github.com/lib/pq"
)

type PostRepository interface {
	Create(ctx context.Context, post *entity.Post) error
	GetByID(ctx context.Context, id int64) (*entity.Post, error)
	Delete(ctx context.Context, postID int64) error
	Update(ctx context.Context, post *entity.Post) error
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
	ctx, cancel := context.WithTimeout(ctx, QueryTimeout)
	defer cancel()

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

func (r *postRepository) GetByID(ctx context.Context, id int64) (*entity.Post, error) {
	query := `
		SELECT id, title, content, user_id, tags, created_at, updated_at, version
		FROM posts WHERE id = $1
	`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeout)
	defer cancel()

	var post entity.Post

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&post.ID,
		&post.Title,
		&post.Content,
		&post.UserID,
		pq.Array(&post.Tags),
		&post.CreatedAt,
		&post.UpdatedAt,
		&post.Version,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}

	return &post, nil
}

func (r *postRepository) Delete(ctx context.Context, postID int64) error {
	ctx, cancel := context.WithTimeout(ctx, QueryTimeout)
	defer cancel()

	res, err := r.db.ExecContext(ctx, "DELETE FROM posts WHERE id = $1", postID)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return ErrNotFound
	}

	return nil
}

func (r *postRepository) Update(ctx context.Context, post *entity.Post) error {
	query := `
		UPDATE posts SET title = $1, content = $2, version = version + 1
		WHERE id = $3 AND version = $4
		RETURNING version
	`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeout)
	defer cancel()

	err := r.db.QueryRowContext(
		ctx,
		query,
		post.Title,
		post.Content,
		post.ID,
		post.Version,
	).Scan(&post.Version)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrNotFound
		default:
			return err
		}
	}

	return nil
}
