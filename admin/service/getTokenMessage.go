package service

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"

	"gpt_mirror/admin/protocol"
	"gpt_mirror/pkg/proxy"

	"github.com/gin-gonic/gin"
)

func GetTokenMessageService(c *gin.Context, req protocol.GetTokenMessageReq) ([]protocol.GetTokenMessageResp, error) {
	// 使用带缓冲的 channel 来收集结果
	resultChan := make(chan protocol.GetTokenMessageResp, len(req.Token))
	errorChan := make(chan error, len(req.Token))

	// 使用 WaitGroup 等待所有 goroutine 完成
	var wg sync.WaitGroup

	// 为每个 token 启动一个 goroutine
	for _, token := range req.Token {
		wg.Add(1)
		go func(v string) {
			defer func() {
				if r := recover(); r != nil {
				}
				wg.Done()
			}()
			// eyJhbGciOi 这个开头
			if strings.HasPrefix(v, "eyJhbGciOi") {
				if isJWT(v) {
					expireTime, email := parseJWT(c, v)
					_, planType := proxy.CheckAccount(v)
					if len(planType) > 0 {
						temp := protocol.GetTokenMessageResp{
							AccessToken:  v,
							Email:        email,
							PlanType:     planType,
							ExpireTime:   expireTime,
							RefreshToken: "",
						}
						resultChan <- temp
					}
				} else {
					// sessionToken
					accessToken, email, planType, _ := sessionChange(c, v)
					expireTime, _ := parseJWT(c, accessToken)
					temp := protocol.GetTokenMessageResp{
						AccessToken:  accessToken,
						Email:        email,
						PlanType:     planType,
						ExpireTime:   expireTime,
						RefreshToken: "",
					}
					resultChan <- temp
				}
			} else {
				// refreshToken
				accessToken, refreshToken, expires, err := proxy.GetAccessToken(v)
				if err != nil {

					errorChan <- err
					return
				}
				_, email := parseJWT(c, accessToken)
				_, planType := proxy.CheckAccount(accessToken)
				temp := protocol.GetTokenMessageResp{
					AccessToken:  accessToken,
					Email:        email,
					PlanType:     planType,
					ExpireTime:   int64(expires),
					RefreshToken: refreshToken,
				}
				resultChan <- temp
			}
		}(token)
	}

	// 启动一个 goroutine 来关闭 channel
	go func() {
		wg.Wait()
		close(resultChan)
		close(errorChan)
	}()

	// 收集结果
	resp := make([]protocol.GetTokenMessageResp, 0, len(req.Token))
	for {
		select {
		case result, ok := <-resultChan:
			fmt.Printf("Received result: %+v,ok:%v\n", result, ok)
			if !ok {
				// channel 已关闭，检查是否有错误

				return resp, nil
			}
			resp = append(resp, result)
		case err := <-errorChan:
			if err != nil {
				return nil, err
			}
		}
	}
}

// 判断是否为 JWT
func isJWT(token string) bool {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return false
	}

	// 检查第一部分是否以 eyJ 开头（{"开头的base64）
	return strings.HasPrefix(parts[0], "eyJ")
}

func sessionChange(c *gin.Context, sessionToken string) (string, string, string, int64) {
	accessToken := ""
	email := ""
	planType := ""
	expireTime := int64(0)
	resp := proxy.SessionTokenToAccountMsg(sessionToken)

	if _, ok := resp["accessToken"]; ok {
		accessToken = resp["accessToken"].(string)
	}
	if _, ok := resp["user"].(map[string]interface{}); ok {
		email = resp["user"].(map[string]interface{})["email"].(string)
	}
	if _, ok := resp["account"].(map[string]interface{}); ok {
		planType = resp["account"].(map[string]interface{})["planType"].(string)
	}
	if _, ok := resp["expires"].(string); ok {
		// 格式化时间
		parseTime, err := time.Parse(time.RFC3339, resp["expires"].(string))
		if err != nil {
		} else {
			expireTime = parseTime.Unix()
		}
	}
	return accessToken, email, planType, expireTime
}

func parseJWT(c *gin.Context, token string) (int64, string) {
	expireTime := int64(0)
	email := ""
	parts := strings.Split(token, ".")
	if len(parts) < 2 {
		return expireTime, email
	}
	payloadPart := parts[1]
	payloadBytes, err := base64.RawURLEncoding.DecodeString(payloadPart)
	if err != nil {
		return expireTime, email
	}
	var claims map[string]interface{}
	if err := json.Unmarshal(payloadBytes, &claims); err != nil {
		return expireTime, email
	}
	if _, ok := claims["exp"].(float64); ok {
		expireTime = int64(claims["exp"].(float64))
	} else {
		expireTime = int64(claims["exp"].(int))
	}
	if _, ok := claims["https://api.openai.com/profile"].(map[string]interface{}); ok {
		email = claims["https://api.openai.com/profile"].(map[string]interface{})["email"].(string)
	}
	return expireTime, email
}
