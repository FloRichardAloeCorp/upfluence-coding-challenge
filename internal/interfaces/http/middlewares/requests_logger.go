package middlewares

import (
	"time"

	"github.com/FloRichardAloeCorp/upfluence-coding-challenge/internal/logs"
	"github.com/gin-gonic/gin"
)

func RequestsLogger(log *logs.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		end := time.Now()

		latency := end.Sub(start)

		fields := []logs.Field{
			{Key: "path", Value: c.Request.URL.Path},
			{Key: "method", Value: c.Request.Method},
			{Key: "latency", Value: latency.String()},
		}

		log.Info("", fields...)
	}
}
