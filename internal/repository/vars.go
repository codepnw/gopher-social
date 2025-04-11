package repository

import (
	"errors"
	"time"
)

var (
	ErrNotFound = errors.New("resource not found")
	ErrConflict = errors.New("resource already exists")
	ErrDuplicateEmail = errors.New("a user with email already exists")
	ErrDuplicateUsername = errors.New("a user with username already exists")

	QueryTimeout = time.Second * 5
)
