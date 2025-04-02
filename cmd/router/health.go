package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (app *Application) HealthCheckHandler(c *gin.Context) {
	data := map[string]string{
		"status": "ok",
		"env":    app.Config.Env,
		"version": app.Config.AppVersion,
	}

	c.JSON(http.StatusOK, data)
}
