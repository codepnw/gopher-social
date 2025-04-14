package handler

import (
	"net/http"

	"github.com/codepnw/gopher-social/internal/utils/logger"
	"github.com/gin-gonic/gin"
)

func internalServerError(c *gin.Context, err error) {
	logger.Error(c, "internal server error", err)
	c.JSON(http.StatusInternalServerError, gin.H{"error": "the server encountered a problem"})
}

func badRequestResponse(c *gin.Context, err error) {
	logger.Warn(c, "bad request response", err)
	c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
}

func conflictResponse(c *gin.Context, err error) {
	logger.Error(c, "conflict error", err)
	c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
}

func notFoundResponse(c *gin.Context, err error) {
	logger.Warn(c, "not found response", err)
	c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
}

func unauthorizedResponse(c *gin.Context, err error) {
	logger.Warn(c, "unauthorized", err)
	c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
}

func responseData(c *gin.Context, status int, data any) {
	c.JSON(status, gin.H{"data": data})
}
