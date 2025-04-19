package response

import (
	"net/http"

	"github.com/codepnw/gopher-social/internal/utils/logger"
	"github.com/gin-gonic/gin"
)

func BadRequestResponse(c *gin.Context, err error) {
	logger.Warn(c, "bad request", err)
	c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
}

func NotFoundResponse(c *gin.Context, err error) {
	logger.Warn(c, "not found", err)
	c.JSON(http.StatusNotFound, gin.H{"status": "error", "message": err.Error()})
}

func InternalServerError(c *gin.Context, err error) {
	logger.Error(c, "internal server", err)
	c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
}

func ResponseData(c *gin.Context, code int, data any) {
	c.JSON(code, gin.H{"status": "success", "data": data})
}