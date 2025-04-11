package router

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
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
	MailExp    time.Duration
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

	shutdown := make(chan error)

	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		s := <-quit

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		app.Logger.Infow("signal caught", "signal", s.String())

		shutdown <- server.Shutdown(ctx)
	}()

	app.Logger.Infow("server has started", "port", app.Config.Addr, "env", app.Config.Env)

	err := server.ListenAndServe()
	if errors.Is(err, http.ErrServerClosed) {
		return err
	}

	err = <-shutdown
	if err != nil {
		return err
	}

	app.Logger.Infow("server has stopped", "port", app.Config.Addr, "env", app.Config.Env)

	return nil
}
