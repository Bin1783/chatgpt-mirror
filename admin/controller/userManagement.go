package controller

import (
	"strconv"

	"gpt_mirror/admin/protocol"
	"gpt_mirror/admin/service"
	"gpt_mirror/response"

	"github.com/gin-gonic/gin"
)

// 创建用户
func CreateUser(c *gin.Context) {
	var req protocol.CreateUserReq
	if err := c.ShouldBindJSON(&req); err != nil {

		response.RespondBadRequest(c, "参数错误")
		return
	}

	err := service.CreateUser(req)
	if err != nil {

		response.RespondBadRequest(c, err.Error())
		return
	}

	response.RespondSuccess(c, gin.H{"message": "创建成功"})
}

// 更新用户
func UpdateUser(c *gin.Context) {
	var req protocol.UpdateUserReq
	if err := c.ShouldBindJSON(&req); err != nil {

		response.RespondBadRequest(c, "参数错误")
		return
	}

	err := service.UpdateUser(req)
	if err != nil {

		response.RespondBadRequest(c, err.Error())
		return
	}

	response.RespondSuccess(c, gin.H{"message": "更新成功"})
}

// 删除用户
func DeleteUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {

		response.RespondBadRequest(c, "ID参数错误")
		return
	}

	req := protocol.DeleteUserReq{ID: id}
	err = service.DeleteUser(req)
	if err != nil {

		response.RespondBadRequest(c, err.Error())
		return
	}

	response.RespondSuccess(c, gin.H{"message": "删除成功"})
}

// 获取用户列表
func GetUserList(c *gin.Context) {
	var req protocol.GetUserListReq
	if err := c.ShouldBindQuery(&req); err != nil {

		response.RespondBadRequest(c, "参数错误")
		return
	}

	resp, err := service.GetUserList(req)
	if err != nil {

		response.RespondBadRequest(c, err.Error())
		return
	}

	response.RespondSuccess(c, resp)
}

// 获取单个用户
func GetUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {

		response.RespondBadRequest(c, "ID参数错误")
		return
	}

	req := protocol.GetUserReq{ID: id}
	resp, err := service.GetUser(req)
	if err != nil {

		response.RespondBadRequest(c, err.Error())
		return
	}

	response.RespondSuccess(c, resp)
}

// 设置用户主账号关联
func SetUserMainAccount(c *gin.Context) {
	var req protocol.SetUserMainAccountReq
	if err := c.ShouldBindJSON(&req); err != nil {

		response.RespondBadRequest(c, "参数错误")
		return
	}

	err := service.SetUserMainAccount(req)
	if err != nil {

		response.RespondBadRequest(c, err.Error())
		return
	}

	response.RespondSuccess(c, gin.H{"message": "设置成功"})
}

// 获取可用主账号列表
func GetAvailableMainAccounts(c *gin.Context) {
	resp, err := service.GetAvailableMainAccounts()
	if err != nil {

		response.RespondBadRequest(c, err.Error())
		return
	}

	response.RespondSuccess(c, resp)
}
