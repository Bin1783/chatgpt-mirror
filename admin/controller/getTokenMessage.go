package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gpt_mirror/admin/protocol"
	"gpt_mirror/admin/service"
	"gpt_mirror/pkg/errCode"
	"gpt_mirror/response"
)

func GetTokenMessage(c *gin.Context) {
	reqData := protocol.GetTokenMessageReq{}
	err := c.ShouldBindJSON(&reqData)
	if err != nil {
		response.RespondError(c, response.CodeBadRequest, errCode.ErrInvalidRequestCode, "请求数据无效")
		return
	}
	//参数校验
	validate := validator.New()
	err = validate.Struct(reqData)
	if err != nil {
		response.RespondError(c, response.CodeBadRequest, errCode.ErrInvalidRequestCode, err.Error())
		return
	}
	resp, err := service.GetTokenMessageService(c, reqData)
	if err != nil {
		response.RespondError(c, response.CodeBadRequest, errCode.ErrInvalidRequestCode, "获取token信息失败")
		return
	}
	response.RespondSuccess(c, resp)
}
