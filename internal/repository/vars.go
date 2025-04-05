package repository

import (
	"errors"
	"time"
)

var (
	ErrNotFound = errors.New("resource not found")
	ErrConflict = errors.New("resource already exists")

	QueryTimeout = time.Second * 5
)
