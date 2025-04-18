package store

import (
	"database/sql"

	"github.com/codepnw/gopher-social/cmd/config"
	"github.com/codepnw/gopher-social/internal/auth"
	"github.com/codepnw/gopher-social/internal/handler"
	"github.com/codepnw/gopher-social/internal/repository"
	"github.com/codepnw/gopher-social/internal/store/cache"
	"github.com/codepnw/gopher-social/internal/utils/mailer"
)

var userRepo repository.UserRepository
var roleRepo repository.RoleRepository

type Storage struct {
	Posts handler.PostsHandler
	Users handler.UserHandler
}

func NewStorage(db *sql.DB, cfg config.Config, mailer mailer.MailtrapClient, cacheStorage cache.Storage) Storage {
	postRepo := repository.NewPostRepository(db)
	commentRepo := repository.NewCommentRepository(db)
	userRepo = repository.NewUserRepository(db)
	roleRepo = repository.NewRoleRepository(db)

	jwtAuth := auth.NewJWTAuthenticator(cfg.Auth.JWTSecret, "", "")

	return Storage{
		Posts: handler.NewPostsHandler(postRepo, commentRepo),
		Users: handler.NewUserHandler(cfg, userRepo, mailer, jwtAuth, cacheStorage),
	}
}

func GetUserRepo() repository.UserRepository {
	return userRepo
}

func GetRoleRepo() repository.RoleRepository {
	return roleRepo
}
