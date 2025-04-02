package repository

import (
	"errors"
	"time"
)

var (
	ErrNotFound = errors.New("resource not found")

	QueryTimeout = time.Second * 5
)
