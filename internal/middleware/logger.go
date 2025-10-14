package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func RequestLogger(logger *logrus.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		latency := time.Since(start).Milliseconds()
		status := c.Writer.Status()

		entry := logger.WithFields(map[string]interface{}{
			"status":     status,
			"method":     c.Request.Method,
			"path":       c.Request.URL.Path,
			"ip":         c.ClientIP(),
			"latency_ms": latency,
		})

		if len(c.Errors) > 0 {
			entry.WithField("errors", c.Errors.String()).Error("request error")
			return
		}

		switch {
		case status >= 500:
			entry.Error("internal error")
		case status >= 400:
			entry.Warn("client error")
		default:
			entry.Info("ok")
		}
	}
}
