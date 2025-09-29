package service

import (
	"context"
	"encoding/json"
	"fmt"
	"gpt_mirror/admin/protocol"
	"gpt_mirror/models/user"
	"gpt_mirror/pkg/redis"
	userRepo "gpt_mirror/repo/user"

	"github.com/gin-gonic/gin"
)

func AddUserService(c *gin.Context, req protocol.AddUserRateLimitReq) error {
	userInfo := make([]user.User, 0)
	pairsMap := make(map[string]interface{})
	for _, item := range req.UserRules {
		rule, _ := json.Marshal(item.Rule)
		userInfo = append(userInfo, user.User{
			ID:            item.UserID,
			ComboID:       item.UserID,
		})
		pairsMap[fmt.Sprintf("limit_config_%d", item.UserID)] = string(rule)
	}
	err := userRepo.SaveUserMsg(userInfo)
	if err != nil {
		return err
	}

	redis.MyRedis.MSet(context.Background(), pairsMap)
	return nil
}
