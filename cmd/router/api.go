package router

import (
	"net/http"
	"time"

	"github.com/codepnw/gopher-social/internal/store"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Application struct {
	Config   Config
	Store    store.Storage
	DBConfig DBConfig
	Logger   *zap.SugaredLogger
}

type Config struct {
	Addr       string
	AppVersion string
	ApiVersion string
	Env        string
}

type DBConfig struct {
	Addr         string
	MaxOpenConns int
	MaxIdleConns int
	MaxIdleTime  string
}

func (app *Application) Run(r *gin.Engine) error {
	server := &http.Server{
		Addr:         ":" + app.Config.Addr,
		Handler:      r,
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Minute,
	}

	app.Logger.Infow("server has started", "port", app.Config.Addr)

	return server.ListenAndServe()
}
