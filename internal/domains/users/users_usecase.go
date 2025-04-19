package users

import (
	"context"
	"database/sql"
	"time"

	"github.com/codepnw/gopher-social/cmd/config"
	"github.com/codepnw/gopher-social/internal/domains/commons"
	"github.com/lib/pq"
)

type UserUsecase interface {
	Create(ctx context.Context, user *UserReq) (*User, error)
	Activate(ctx context.Context, token string) error
	CreateAndInvite(ctx context.Context, user *UserReq, token string, exp time.Duration) error
	GetByID(ctx context.Context, id int64) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	Follow(ctx context.Context, followerID, userID int64) error
	Unfollow(ctx context.Context, followerID, userID int64) error
	Delete(ctx context.Context, userID int64) error
}

type usecase struct {
	db     *sql.DB
	repo   UserRepository
	config config.Config
}

func NewUserUsecase(db *sql.DB, repo UserRepository, config config.Config) UserUsecase {
	return &usecase{
		repo:   repo,
		config: config,
	}
}

func (uc *usecase) Create(ctx context.Context, user *UserReq) (*User, error) {
	var u User

	role := u.Role.Name
	if role == "" {
		role = "user"
	}

	u = User{
		Email:    user.Email,
		Username: user.Username,
		Password: user.Password,
		Role: Role{
			Name: role,
		},
	}

	ctx, cancel := context.WithTimeout(ctx, commons.ContextQueryTimeout)
	defer cancel()

	tx, _ := uc.db.BeginTx(ctx, nil)

	err := uc.repo.Create(ctx, tx, &u)
	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "users_email_key"`:
			return nil, commons.ErrDuplicateEmail
		case err.Error() == `pq: duplicate key value violates unique constraint "users_username_key"`:
			return nil, commons.ErrDuplicateUsername
		default:
			return nil, err
		}
	}

	return &u, nil
}

func (uc *usecase) Activate(ctx context.Context, token string) error {
	ctx, cancel := context.WithTimeout(ctx, commons.ContextQueryTimeout)
	defer cancel()

	return uc.repo.Activate(ctx, token)
}

func (uc *usecase) CreateAndInvite(ctx context.Context, user *UserReq, token string, exp time.Duration) error {
	ctx, cancel := context.WithTimeout(ctx, commons.ContextQueryTimeout)
	defer cancel()

	u := &User{
		Email:    user.Email,
		Username: user.Username,
		Password: user.Password,
		Role: Role{
			Name: "user",
		},
	}

	err := uc.repo.CreateAndInvite(ctx, u, token, uc.config.Auth.JWTExp)
	if err != nil {
		return err
	}

	return nil
}

func (uc *usecase) GetByID(ctx context.Context, id int64) (*User, error) {
	ctx, cancel := context.WithTimeout(ctx, commons.ContextQueryTimeout)
	defer cancel()

	user, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, commons.ErrNotFound
		default:
			return nil, err
		}
	}

	return user, nil
}

func (uc *usecase) GetByEmail(ctx context.Context, email string) (*User, error) {
	ctx, cancel := context.WithTimeout(ctx, commons.ContextQueryTimeout)
	defer cancel()

	user, err := uc.repo.GetByEmail(ctx, email)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, commons.ErrNotFound
		default:
			return nil, err
		}
	}

	return user, nil
}

func (uc *usecase) Delete(ctx context.Context, userID int64) error {
	ctx, cancel := context.WithTimeout(ctx, commons.ContextQueryTimeout)
	defer cancel()

	return uc.repo.Delete(ctx, userID)
}

func (uc *usecase) Follow(ctx context.Context, followerID, userID int64) error {
	ctx, cancel := context.WithTimeout(ctx, commons.ContextQueryTimeout)
	defer cancel()

	if err := uc.repo.Follow(ctx, followerID, userID); err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			return commons.ErrConflict
		}
	}

	return nil
}
func (uc *usecase) Unfollow(ctx context.Context, followerID, userID int64) error {
	ctx, cancel := context.WithTimeout(ctx, commons.ContextQueryTimeout)
	defer cancel()

	return uc.repo.Unfollow(ctx, followerID, userID)
}
