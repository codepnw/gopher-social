package users

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"time"

	"github.com/codepnw/gopher-social/internal/domains/commons"
)

type UserRepository interface {
	Create(ctx context.Context, tx *sql.Tx, user *User) error
	Activate(ctx context.Context, token string) error
	CreateAndInvite(ctx context.Context, user *User, token string, exp time.Duration) error
	GetByID(ctx context.Context, id int64) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	Delete(ctx context.Context, userID int64) error

	Follow(ctx context.Context, followerID, userID int64) error
	Unfollow(ctx context.Context, followerID, userID int64) error
}

type repository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &repository{db: db}
}

func (r *repository) Create(ctx context.Context, tx *sql.Tx, user *User) error {
	query := `
		INSERT INTO users (username, email, password, role_id)
		VALUES ($1, $2, $3, (SELECT id FROM roles WHERE name = $4)) 
		RETURNING id, created_at
	`
	err := r.db.QueryRowContext(
		ctx,
		query,
		user.Username,
		user.Email,
		user.Password,
		user.Role.Name,
	).Scan(&user.ID, &user.CreatedAt)

	if err != nil {
		return err
	}

	return nil
}

func (r *repository) CreateAndInvite(ctx context.Context, user *User, token string, invitationExp time.Duration) error {
	return commons.WithTransaction(ctx, r.db, func(tx *sql.Tx) error {
		// create user
		if err := r.Create(ctx, tx, user); err != nil {
			return err
		}

		// create user invite
		if err := r.createUserInvitation(ctx, tx, token, invitationExp, user.ID); err != nil {
			return err
		}

		return nil
	})
}

func (r *repository) createUserInvitation(ctx context.Context, tx *sql.Tx, token string, exp time.Duration, userID int64) error {
	query := `INSERT INTO user_invitations (token, user_id, expiry) VALUES ($1, $2, $3)`

	_, err := tx.ExecContext(ctx, query, token, userID, time.Now().Add(exp))
	if err != nil {
		return err
	}

	return nil
}

func (r *repository) Activate(ctx context.Context, token string) error {
	return commons.WithTransaction(ctx, r.db, func(tx *sql.Tx) error {
		// find user token
		user, err := r.getUserFromInvitation(ctx, tx, token)
		if err != nil {
			return err
		}

		// update user
		user.IsActive = true
		if err := r.update(ctx, tx, user); err != nil {
			return err
		}

		// clean invitations
		if err := r.deleteUserInvitations(ctx, tx, user.ID); err != nil {
			return err
		}

		return nil
	})
}

// Acticate Method
func (r *repository) getUserFromInvitation(ctx context.Context, tx *sql.Tx, token string) (*User, error) {
	query := `
		SELECT u.id, u.username, u.email, u.created_at, u.is_active
		FROM users u
		JOIN user_invitations ui ON u.id = ui.user_id
		WHERE ui.token = $1 AND ui.expiry > $2
	`
	hash := sha256.Sum256([]byte(token))
	hashToken := hex.EncodeToString(hash[:])

	var user User
	err := tx.QueryRowContext(ctx, query, hashToken, time.Now()).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.CreatedAt,
		&user.IsActive,
	)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, commons.ErrNotFound
		default:
			return nil, err
		}
	}

	return &user, nil
}

func (r *repository) update(ctx context.Context, tx *sql.Tx, user *User) error {
	query := `UPDATE users SET username = $1, email = $2, is_active = $3 WHERE id = $4`

	_, err := tx.ExecContext(ctx, query, user.Username, user.Email, user.IsActive, user.ID)
	if err != nil {
		return err
	}

	return nil
}

func (r *repository) deleteUserInvitations(ctx context.Context, tx *sql.Tx, userID int64) error {
	query := `DELETE FROM user_invitations WHERE user_id = $1`

	_, err := tx.ExecContext(ctx, query, userID)
	if err != nil {
		return err
	}

	return nil
}

// End Acticate Method

func (r *repository) GetByID(ctx context.Context, id int64) (*User, error) {
	query := `
		SELECT users.id, username, email, password, created_at, roles.* 
		FROM users
		JOIN roles ON (users.role_id = roles.id)
		WHERE users.id = $1 AND is_active = true
	`
	var user User
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
		&user.Role.ID,
		&user.Role.Name,
		&user.Role.Level,
		&user.Role.Description,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *repository) GetByEmail(ctx context.Context, email string) (*User, error) {
	query := `
		SELECT id, username, email, password, created_at FROM users
		WHERE email = $1 AND is_active = true
	`
	var user User
	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *repository) Delete(ctx context.Context, userID int64) error {
	return commons.WithTransaction(ctx, r.db, func(tx *sql.Tx) error {
		if err := r.delete(ctx, tx, userID); err != nil {
			return err
		}

		if err := r.deleteUserInvitations(ctx, tx, userID); err != nil {
			return err
		}

		return nil
	})
}

func (r *repository) delete(ctx context.Context, tx *sql.Tx, userID int64) error {
	_, err := tx.ExecContext(ctx, "DELETE FROM users WHERE id = $1", userID)
	if err != nil {
		return err
	}

	return nil
}

func (r *repository) Follow(ctx context.Context, followerID, userID int64) error {
	query := `INSERT INTO followers (user_id, follower_id) VALUES ($1, $2)`

	_, err := r.db.ExecContext(ctx, query, userID, followerID)
	if err != nil {
		return err
	}

	return nil
}

func (r *repository) Unfollow(ctx context.Context, followerID, userID int64) error {
	query := `DELETE FROM followers WHERE user_id = $1 AND follower_id = $2`

	_, err := r.db.ExecContext(ctx, query, userID, followerID)
	return err
}
