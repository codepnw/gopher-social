package handler

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func internalServerError(c *gin.Context, err error) {
	log.Printf("internal server error: %s path: %s error: %s", c.Request.Method, c.Request.URL.Path, err.Error())
	c.JSON(http.StatusInternalServerError, gin.H{"error": "the server encountered a problem"})
}

func badRequestResponse(c *gin.Context, err error) {
	log.Printf("bad request error: %s path: %s error: %s", c.Request.Method, c.Request.URL.Path, err.Error())
	c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
}

func conflictResponse(c *gin.Context, err error) {
	log.Printf("conflict error: %s path: %s error: %s", c.Request.Method, c.Request.URL.Path, err.Error())
	c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
}

func notFoundResponse(c *gin.Context, err error) {
	log.Printf("not found error: %s path: %s error: %s", c.Request.Method, c.Request.URL.Path, err.Error())
	c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
}

func responseData(c *gin.Context, status int, data any) {
	c.JSON(status, gin.H{"data": data})
}