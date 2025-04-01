package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (app *Application) HealthCheckHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "OK!"})
}