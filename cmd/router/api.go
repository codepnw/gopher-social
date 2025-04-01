package router

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type Application struct {
	Config Config
}

type Config struct {
	Addr    string
	Version string
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
