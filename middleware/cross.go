package middleware

import (
	"github.com/gin-gonic/gin"
	"strings"
)

// 允许跨域
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")

		// 允许的头部
		allowedHeaders := []string{
			"Content-Type",
			"Authorization",
			"X-Requested-With",
			"Accept",
			"Origin",
			"User-Agent",
			"DNT",
			"Cache-Control",
			"X-Mx-ReqToken",
			"Keep-Alive",
			"X-Requested-With",
			"If-Modified-Since",
			"X-CSRF-Token",
			"X-OpenAI-Assistant-App-Id",
			"OpenAI-Conversation-ID",
			"OpenAI-Ephemeral-User-ID",
			"OpenAI-Organization",
			"OAI-Device-Id",
			"OAI-Language",
		}

		// 设置CORS头部
		if origin != "" {
			c.Header("Access-Control-Allow-Origin", origin)
		} else {
			c.Header("Access-Control-Allow-Origin", "*")
		}

		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, PATCH, HEAD")
		c.Header("Access-Control-Allow-Headers", strings.Join(allowedHeaders, ", "))
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Content-Range, Content-Type, Authorization")
		c.Header("Access-Control-Max-Age", "86400")

		// 处理预检请求
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
