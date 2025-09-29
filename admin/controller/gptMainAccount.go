package controller

import (
	"net/http"
	"strconv"

	"gpt_mirror/admin/protocol"
	"gpt_mirror/admin/service"
	"gpt_mirror/response"

	"github.com/gin-gonic/gin"
)

// 管理员登录
func AdminLogin(c *gin.Context) {
	var req protocol.AdminLoginReq
	if err := c.ShouldBindJSON(&req); err != nil {

		response.RespondBadRequest(c, "参数错误")
		return
	}

	resp, err := service.AdminLogin(req)
	if err != nil {

		response.RespondBadRequest(c, err.Error())
		return
	}

	response.RespondSuccess(c, resp)
}

// 创建主账号
func CreateGptMainAccount(c *gin.Context) {
	var req protocol.CreateGptMainAccountReq
	if err := c.ShouldBindJSON(&req); err != nil {

		response.RespondBadRequest(c, "参数错误")
		return
	}

	err := service.CreateGptMainAccount(req)
	if err != nil {

		response.RespondBadRequest(c, err.Error())
		return
	}

	response.RespondSuccess(c, gin.H{"message": "创建成功"})
}

// 更新主账号
func UpdateGptMainAccount(c *gin.Context) {
	var req protocol.UpdateGptMainAccountReq
	if err := c.ShouldBindJSON(&req); err != nil {

		response.RespondBadRequest(c, "参数错误")
		return
	}

	err := service.UpdateGptMainAccount(req)
	if err != nil {

		response.RespondBadRequest(c, err.Error())
		return
	}

	response.RespondSuccess(c, gin.H{"message": "更新成功"})
}

// 删除主账号
func DeleteGptMainAccount(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {

		response.RespondBadRequest(c, "ID参数错误")
		return
	}

	req := protocol.DeleteGptMainAccountReq{ID: id}
	err = service.DeleteGptMainAccount(req)
	if err != nil {

		response.RespondBadRequest(c, err.Error())
		return
	}

	response.RespondSuccess(c, gin.H{"message": "删除成功"})
}

// 获取主账号列表
func GetGptMainAccountList(c *gin.Context) {
	var req protocol.GetGptMainAccountListReq
	if err := c.ShouldBindQuery(&req); err != nil {

		response.RespondBadRequest(c, "参数错误")
		return
	}

	resp, err := service.GetGptMainAccountList(req)
	if err != nil {

		response.RespondBadRequest(c, err.Error())
		return
	}

	response.RespondSuccess(c, resp)
}

// 获取单个主账号
func GetGptMainAccount(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {

		response.RespondBadRequest(c, "ID参数错误")
		return
	}

	req := protocol.GetGptMainAccountReq{ID: id}
	resp, err := service.GetGptMainAccount(req)
	if err != nil {

		response.RespondBadRequest(c, err.Error())
		return
	}

	response.RespondSuccess(c, resp)
}

// 管理后台首页
func AdminIndex(c *gin.Context) {
	c.HTML(http.StatusOK, "admin.html", gin.H{
		"title": "GPT主账号管理后台",
	})
}
