package commons

import (
	"errors"
	"time"
)

const (
	ContextPostKey = "post"
	ContextUserKey = "user"

	ContextQueryTimeout = 5 * time.Second
)

var (
	ErrNotFound          = errors.New("resource not found")
	ErrConflict          = errors.New("resource already exists")
	ErrDuplicateEmail    = errors.New("a user with email already exists")
	ErrDuplicateUsername = errors.New("a user with username already exists")

	ErrInvalidEmailPassword = errors.New("invalid email or password")
)
