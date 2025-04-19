package cache

import (
	"context"

	"github.com/codepnw/gopher-social/internal/domains/users"
	"github.com/redis/go-redis/v9"
)

type Storage struct {
	Users interface {
		Get(context.Context, int64) (*users.User, error)
		Set(context.Context, *users.User) error
	}
}

func NewRedisStorage(rdb *redis.Client) Storage {
	return Storage{
		Users: &UserStore{rdb: rdb},
	}
}
