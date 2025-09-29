package middleware

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"gpt_mirror/pkg/createShareToken"
	productMainAccountToProduct "gpt_mirror/repo/producctMainAccountToProduct"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/spf13/viper"
)

const TraceIDKey = "trace_id"

func ShareTokenAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		url := viper.GetString(fmt.Sprintf("redirect.%s", viper.GetString("env")))
		fmt.Println("url", url)

		queryToken := c.Query("share_token")
		if queryToken != "" {

			tm, err := createShareToken.NewTokenManager()
			if err != nil {

				c.Redirect(http.StatusFound, url)
				c.Abort()
				return
			}

			userId, err := tm.DecryptToken(queryToken)
			if err != nil {

				c.Redirect(http.StatusFound, url)
				c.Abort()
				return
			}
			c.Set("userId", userId)

			expireTime := time.Date(2099, 1, 1, 0, 0, 0, 0, time.UTC)
			c.SetCookie("share_token", queryToken, int(expireTime.Sub(time.Now()).Seconds()), "/", "", false, true)

			cleanPath := c.Request.URL.Path

			c.Redirect(http.StatusTemporaryRedirect, cleanPath)
			c.Abort() // 必须终止，因为已经发出了重定向响应
			return
		}
		cookieToken, err := c.Cookie("share_token")
		if err != nil || cookieToken == "" {

			c.Redirect(http.StatusFound, url)
			c.Abort()
			return
		}

		tm, err := createShareToken.NewTokenManager()
		if err != nil {

			c.Redirect(http.StatusFound, url)
			c.Abort()
			return
		}

		userId, err := tm.DecryptToken(cookieToken)
		if err != nil {
			var errMsg string
			switch {
			case errors.Is(err, createShareToken.ErrTokenExpired):
				errMsg = "授权令牌已过期"
			case errors.Is(err, createShareToken.ErrInvalidToken):
				errMsg = "无效的授权令牌"
			default:
				errMsg = "内部服务器错误"

			}
			fmt.Println("errMsg", errMsg)
			c.SetCookie("share_token", "", -1, "/", "", false, true)

			c.Redirect(http.StatusFound, url)
			c.Abort()
			return
		}
		c.Set("userId", userId)
		if !strings.HasSuffix(c.Request.URL.Path, ".js") && !strings.HasSuffix(c.Request.URL.Path, ".js.map") && !strings.HasSuffix(c.Request.URL.Path, ".css") {
			id, err := getMianAccountId(userId)
			if err != nil {
				c.Redirect(http.StatusFound, url)
				c.Abort()
				return
			}

			c.Set("mainAccount", id)
		}
		traceID := c.GetHeader("X-Trace-ID")
		if traceID == "" {
			traceID = c.GetHeader("X-Request-ID")
		}
		if traceID == "" {
			traceID = uuid.New().String()
		}
		ctx := context.WithValue(c.Request.Context(), TraceIDKey, traceID)
		c.Request = c.Request.WithContext(ctx)
		c.Header("X-Trace-ID", traceID)
		c.Set(TraceIDKey, traceID)
		c.Next()
	}
}

func getMianAccountId(userId int) (int, error) {
	product, err := productMainAccountToProduct.GetProductMainAccountToUser(userId)
	if err != nil {
		return 0, err
	}
	return product.MainAccountID, nil
}
