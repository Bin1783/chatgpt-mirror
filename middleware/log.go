package middleware

import (
	"github.com/gin-gonic/gin"
)

func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		//	startTime := time.Now()
		//	c.Next()
		//	endTime := time.Now()
		//
		//	latencyTime := endTime.Sub(startTime)
		//
		//	reqMethod := c.Request.Method
		//	reqURI := c.Request.RequestURI
		//
		//	statusCode := c.Writer.Status()
		//
		//	clientIP := c.ClientIP()
		//entry := logrus.WithFields(logrus.Fields{
		//	"status_code":  statusCode,
		//	"latency_time": latencyTime,
		//	"client_ip":    clientIP,
		//	"req_method":   reqMethod,
		//	"req_uri":      reqURI,
		//})
		//
		//if len(c.Errors) > 0 {
		//	entry.Error(c.Errors.ByType(gin.ErrorTypePrivate).String())
		//} else {
		//	if statusCode >= 500 {
		//		entry.Error("Server error")
		//	} else if statusCode >= 400 {
		//		entry.Warn("Client error")
		//	} else {
		//		entry.Info("Request handled successfully")
		//	}
		//}
	}
}
