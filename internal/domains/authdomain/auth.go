package authdomain

import (
	"database/sql"

	"github.com/codepnw/gopher-social/cmd/config"
	"github.com/codepnw/gopher-social/internal/auth"
	"github.com/codepnw/gopher-social/internal/domains/users"
)

func InitAuthDomain(db *sql.DB, cfg config.Config, jwt *auth.JWTAuthenticator) AuthHandler {
	userrepo := users.NewUserRepository(db)
	useruc := users.NewUserUsecase(db, userrepo, cfg)

	uc := NewAuthUsecase(useruc)
	hdl := NewAuthHandler(uc, cfg, jwt)

	return hdl
}
