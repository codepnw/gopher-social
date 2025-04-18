package cache

import (
	"context"

	"github.com/codepnw/gopher-social/internal/entity"
	"github.com/redis/go-redis/v9"
)

type Storage struct {
	Users interface {
		Get(context.Context, int64) (*entity.User, error)
		Set(context.Context, *entity.User) error
	}
}

func NewRedisStorage(rdb *redis.Client) Storage {
	return Storage{
		Users: &UserStore{rdb: rdb},
	}
}