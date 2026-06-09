package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
)

func Logger() gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		return param.TimeStamp.Format(time.RFC3339) + " " +
			param.Method + " " + param.Path + " " +
			param.StatusCodeString() + " " + param.Latency.String() + "
"
	})
}
