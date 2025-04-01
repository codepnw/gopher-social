package router

import "github.com/gin-gonic/gin"

func (app *Application) Routes() *gin.Engine {
	r := gin.Default()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	v := r.Group("/" + app.Config.Version)

	v.GET("/health", app.HealthCheckHandler)

	return r
}
