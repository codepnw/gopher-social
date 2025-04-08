package logger

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

var logger *zap.SugaredLogger

func InitLogger() (*zap.SugaredLogger, error) {
	config := zap.NewProductionConfig()
	config.DisableStacktrace = true
	config.EncoderConfig.TimeKey = "" // disable ts (timestamp) field

	build, err := config.Build()
	if err != nil {
		return nil, err
	}
	logger = build.Sugar()

	return logger, nil
}

func Error(c *gin.Context, msg string, err error) {
	logger.Errorw(msg, "method", c.Request.Method, "path", c.Request.URL.Path, "error", err.Error())
}

func Warn(c *gin.Context, msg string, err error) {
	logger.Warn(msg, "method", c.Request.Method, "path", c.Request.URL.Path, "error", err.Error())
}