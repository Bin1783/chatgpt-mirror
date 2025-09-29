package protocol

import "time"

// 创建用户请求
type CreateUserReq struct {
	Name          string `json:"name" validate:"required"`
	Password      string `json:"password" validate:"required"`
	ComboID       int    `json:"combo_id" validate:"required"`
	MainAccountID int    `json:"main_account_id"` // 关联的主账号ID
}

// 更新用户请求
type UpdateUserReq struct {
	ID            int    `json:"id" validate:"required"`
	Name          string `json:"name"`
	Password      string `json:"password"`
	ComboID       int    `json:"combo_id"`
	MainAccountID int    `json:"main_account_id"` // 关联的主账号ID
}

// 删除用户请求
type DeleteUserReq struct {
	ID int `json:"id" validate:"required"`
}

// 获取用户列表请求
type GetUserListReq struct {
	Page     int    `json:"page" form:"page"`
	PageSize int    `json:"page_size" form:"page_size"`
	Name     string `json:"name" form:"name"`         // 用户名搜索
	ComboID  int    `json:"combo_id" form:"combo_id"` // 套餐ID筛选
}

// 获取单个用户请求
type GetUserReq struct {
	ID int `json:"id" validate:"required" form:"id"`
}

// 用户响应
type UserResp struct {
	ID               int       `json:"id"`
	Name             string    `json:"name"`
	ComboID          int       `json:"combo_id"`
	MainAccountID    int       `json:"main_account_id"`
	MainAccountToken string    `json:"main_account_token"` // 关联的主账号token（脱敏）
	CreatedTime      time.Time `json:"created_time"`
	UpdatedTime      time.Time `json:"updated_time"`
}

// 用户列表响应
type UserListResp struct {
	List  []UserResp `json:"list"`
	Total int64      `json:"total"`
	Page  int        `json:"page"`
	Size  int        `json:"size"`
}

// 设置用户主账号关联请求
type SetUserMainAccountReq struct {
	UserID        int `json:"user_id" validate:"required"`
	MainAccountID int `json:"main_account_id" validate:"required"`
}

// 获取可用主账号列表响应
type AvailableMainAccountResp struct {
	ID         int    `json:"id"`
	Token      string `json:"token"` // 脱敏后的token
	Status     int    `json:"status"`
	StatusText string `json:"status_text"`
}
