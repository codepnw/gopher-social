package router

import (
	"log"
	"net/http"
	"time"

	"github.com/codepnw/gopher-social/internal/store"
	"github.com/gin-gonic/gin"
)

type Application struct {
	Config   Config
	Store    store.Storage
	DBConfig DBConfig
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
		Addr:         app.Config.Addr,
		Handler:      r,
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Minute,
	}

	log.Printf("server has started at %s", app.Config.Addr)

	return server.ListenAndServe()
}
