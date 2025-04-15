package repository

import (
	"context"
	"database/sql"

	"github.com/codepnw/gopher-social/internal/entity"
)

type RoleRepository interface {
	GetByName(ctx context.Context, name string) (*entity.Role, error)
}

type roleRepository struct {
	db *sql.DB
}

func NewRoleRepository(db *sql.DB) RoleRepository {
	return &roleRepository{db: db}
}

func (r *roleRepository) GetByName(ctx context.Context, name string) (*entity.Role, error) {
	query := `SELECT id, name, description, level FROM roles WHERE name = $1`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeout)
	defer cancel()

	var role entity.Role
	err := r.db.QueryRowContext(ctx, query, name).Scan(
		&role.ID,
		&role.Name,
		&role.Description,
		&role.Level,
	)
	if err != nil {
		return nil, err
	}

	return &role, nil
}
