package protocol

import "time"

// 创建主账号请求
type CreateGptMainAccountReq struct {
	Token        string `json:"token" validate:"required"`
	RefreshToken string `json:"refresh_token"`
	Status       int    `json:"status" validate:"required"`
}

// 更新主账号请求
type UpdateGptMainAccountReq struct {
	ID           int    `json:"id" validate:"required"`
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
	Status       int    `json:"status"`
}

// 删除主账号请求
type DeleteGptMainAccountReq struct {
	ID int `json:"id" validate:"required"`
}

// 获取主账号列表请求
type GetGptMainAccountListReq struct {
	Page     int `json:"page" form:"page"`
	PageSize int `json:"page_size" form:"page_size"`
	Status   int `json:"status" form:"status"`
}

// 获取单个主账号请求
type GetGptMainAccountReq struct {
	ID int `json:"id" validate:"required" form:"id"`
}

// 主账号响应
type GptMainAccountResp struct {
	ID           int       `json:"id"`
	Token        string    `json:"token"`
	RefreshToken string    `json:"refresh_token"`
	Status       int       `json:"status"`
	StatusText   string    `json:"status_text"`
	CreateTime   time.Time `json:"create_time"`
	UpdateTime   time.Time `json:"update_time"`
}

// 主账号列表响应
type GptMainAccountListResp struct {
	List  []GptMainAccountResp `json:"list"`
	Total int64                `json:"total"`
	Page  int                  `json:"page"`
	Size  int                  `json:"size"`
}

// 管理员登录请求
type AdminLoginReq struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

// 管理员登录响应
type AdminLoginResp struct {
	Token    string `json:"token"`
	Username string `json:"username"`
	IsRedirect bool `json:"is_redirect"`
	RedirectUrl string `json:"redirect_url"`
}
