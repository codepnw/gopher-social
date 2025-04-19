package users

import (
	"database/sql"

	"github.com/codepnw/gopher-social/cmd/config"
)

func InitUserDomain(db *sql.DB, cfg config.Config) UserHandler {
	repo := NewUserRepository(db)
	uc := NewUserUsecase(db, repo, cfg)
	hdl := NewUserHandler(uc)

	return hdl
}
