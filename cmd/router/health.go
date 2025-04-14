package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (app *Application) HealthCheckHandler(c *gin.Context) {
	user := c.MustGet(gin.AuthUserKey).(string)

	data := map[string]string{
		"status":  "ok",
		"env":     app.Config.App.Env,
		"version": app.Config.App.AppVersion,
		"user":    user,
	}

	c.JSON(http.StatusOK, data)
}

func (app *Application) BasicAuthLogout(c *gin.Context) {
	c.Header("WWW-Authenticate", `Basic realm="Please re-login"`)
	c.AbortWithStatus(http.StatusUnauthorized)
}
