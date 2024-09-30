package logger

import (
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

var (
	Logger                      *zap.Logger
	Error, Info, Warning, Debug func(string, ...zap.Field)
	Infof                       func(template string, args ...interface{})
)

func InitLogger(level string) {
	var err error
	switch level {
	case "production":
		Logger, err = zap.NewProductionConfig().Build()
	default:
		Logger, err = zap.NewDevelopmentConfig().Build()
	}
	if err != nil {
		panic(err)
	}
	Error = Logger.Error
	Warning = Logger.Warn
	Info = Logger.Info
	Debug = Logger.Debug
	Infof = Logger.Sugar().Infof
}

func LoggerForGin(c *gin.Context) {
	start := time.Now()
	c.Next()
	end := time.Now()
	Logger.Info("Incoming request",
		zap.String("method", c.Request.Method),
		zap.String("path", c.Request.URL.Path),
		zap.Int("status", c.Writer.Status()),
		zap.String("client_ip", c.ClientIP()),
		zap.Duration("latency", end.Sub(start)),
	)
}
