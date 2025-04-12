package router

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/codepnw/gopher-social/cmd/config"
	"github.com/codepnw/gopher-social/internal/store"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Application struct {
	Config  config.Config
	Store  store.Storage
	Logger *zap.SugaredLogger
}

func (app *Application) Run(r *gin.Engine) error {
	addr := app.Config.App.Addr
	env := app.Config.App.Env

	server := &http.Server{
		Addr:         ":" + addr,
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

	app.Logger.Infow("server has started", "port", addr, "env", env)

	err := server.ListenAndServe()
	if errors.Is(err, http.ErrServerClosed) {
		return err
	}

	err = <-shutdown
	if err != nil {
		return err
	}

	app.Logger.Infow("server has stopped", "port", addr, "env", env)

	return nil
}
