package auth

import (
	"database/sql"

	"github.com/codepnw/gopher-social/cmd/config"
	"github.com/codepnw/gopher-social/internal/domains/users"
)

func InitAuthDomain(db *sql.DB, cfg config.Config) AuthHandler {
	userrepo := users.NewUserRepository(db)
	useruc := users.NewUserUsecase(userrepo, cfg)

	uc := NewAuthUsecase(useruc)
	hdl := NewAuthHandler(uc, cfg)

	return hdl
}
