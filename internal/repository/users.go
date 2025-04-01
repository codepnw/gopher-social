package repository

import (
	"context"
	"database/sql"

	"github.com/codepnw/gopher-social/internal/entity"
)

type UserRepository interface {
	Create(ctx context.Context, user *entity.User) error
}

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, user *entity.User) error {
	query := `
		INSERT INTO users (username, email, password)
		VALUES ($1, $2, $3) RETURNING id, created_at
	`
	err := r.db.QueryRowContext(
		ctx,
		query,
		user.Username,
		user.Email,
		user.Password,
	).Scan(&user.ID, &user.CreatedAt)
	
	if err != nil {
		return err
	}

	return nil
}
