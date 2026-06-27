package middleware

import (
	"blog/utils"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// LoggerMiddleware 日志记录中间件
func LoggerMiddleware() gin.HandlerFunc {
	// 创建logrus实例
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})

	return gin.LoggerWithFormatter(func(params gin.LogFormatterParams) string {
		// 记录请求日志
		logger.WithFields(logrus.Fields{
			"status_code": params.StatusCode,
			"latency":     params.Latency,
			"client_ip":   params.ClientIP,
			"method":      params.Method,
			"path":        params.Path,
			"user_agent":  params.Request.UserAgent(),
			"error":       params.ErrorMessage,
			"timestamp":   params.TimeStamp.Format(time.RFC3339),
		}).Info("HTTP Request")
		return ""
	})
}

// ErrorHandlerMiddleware 全局错误处理中间件
func ErrorHandlerMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				logrus.WithFields(logrus.Fields{
					"error":  err,
					"path":   ctx.Request.URL.Path,
					"method": ctx.Request.Method,
				}).Error("Panic recovered")
				utils.InternalServerError(ctx, "Internal server error")
				ctx.Abort()
			}
		}()
		ctx.Next()
	}
}
