package service

import (
	"errors"
	"math/rand"
	"time"

	"gpt_mirror/admin/protocol"
	gptMainAccountModel "gpt_mirror/models/gptMainAccount"
	"gpt_mirror/pkg/createShareToken"
	gptMainAccountRepo "gpt_mirror/repo/gptMainAccount"
	productMainAccountToProductRepo "gpt_mirror/repo/producctMainAccountToProduct"
	userRepo "gpt_mirror/repo/user"

	"github.com/spf13/viper"
	"gorm.io/gorm"
)

// 创建主账号
func CreateGptMainAccount(req protocol.CreateGptMainAccountReq) error {
	account := gptMainAccountModel.GptMainAccount{
		Token:        req.Token,
		RefreshToken: req.RefreshToken,
		Status:       req.Status,
		CreateTime:   time.Now(),
		UpdateTime:   time.Now(),
	}

	return gptMainAccountRepo.CreateAccount(account)
}

// 更新主账号
func UpdateGptMainAccount(req protocol.UpdateGptMainAccountReq) error {
	// 检查账号是否存在
	_, err := gptMainAccountRepo.GetAccountByID(req.ID)
	if err != nil {
		return errors.New("账号不存在")
	}

	updateData := make(map[string]interface{})
	if req.Token != "" {
		updateData["token"] = req.Token
	}
	if req.RefreshToken != "" {
		updateData["refresh_token"] = req.RefreshToken
	}
	if req.Status != 0 {
		updateData["status"] = req.Status
	}
	updateData["update_time"] = time.Now()

	return gptMainAccountRepo.UpdateAccountById(req.ID, updateData)
}

// 删除主账号
func DeleteGptMainAccount(req protocol.DeleteGptMainAccountReq) error {
	// 检查账号是否存在
	_, err := gptMainAccountRepo.GetAccountByID(req.ID)
	if err != nil {
		return errors.New("账号不存在")
	}

	return gptMainAccountRepo.DeleteAccountById(req.ID)
}

// 获取主账号列表
func GetGptMainAccountList(req protocol.GetGptMainAccountListReq) (*protocol.GptMainAccountListResp, error) {
	// 设置默认分页参数
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}

	where := make(map[string]interface{})
	if req.Status != 0 {
		where["status"] = req.Status
	}

	// 获取总数
	total, err := gptMainAccountRepo.GetAccountsCount(where)
	if err != nil {
		return nil, err
	}

	// 获取账号列表
	accounts, err := gptMainAccountRepo.GetAccountsByPage(where, req.Page, req.PageSize)
	if err != nil {
		return nil, err
	}

	// 转换为响应格式
	var list []protocol.GptMainAccountResp
	for _, account := range accounts {
		statusText := getStatusText(account.Status)
		list = append(list, protocol.GptMainAccountResp{
			ID:           account.ID,
			Token:        maskToken(account.Token),
			RefreshToken: maskToken(account.RefreshToken),
			Status:       account.Status,
			StatusText:   statusText,
			CreateTime:   account.CreateTime,
			UpdateTime:   account.UpdateTime,
		})
	}

	return &protocol.GptMainAccountListResp{
		List:  list,
		Total: total,
		Page:  req.Page,
		Size:  req.PageSize,
	}, nil
}

// 获取单个主账号
func GetGptMainAccount(req protocol.GetGptMainAccountReq) (*protocol.GptMainAccountResp, error) {
	account, err := gptMainAccountRepo.GetAccountByID(req.ID)
	if err != nil {
		return nil, errors.New("账号不存在")
	}

	statusText := getStatusText(account.Status)
	return &protocol.GptMainAccountResp{
		ID:           account.ID,
		Token:        account.Token,
		RefreshToken: account.RefreshToken,
		Status:       account.Status,
		StatusText:   statusText,
		CreateTime:   account.CreateTime,
		UpdateTime:   account.UpdateTime,
	}, nil
}

// 管理员登录
func AdminLogin(req protocol.AdminLoginReq) (*protocol.AdminLoginResp, error) {
	// 简单的用户名密码验证（实际项目中应该从数据库验证）
	if req.Username != "BinRoot" || req.Password != "BinRoot" {
		// 查询数据库
		user, err := userRepo.GetUserMsg(map[string]interface{}{"name": req.Username})
		if err != nil {
			return nil, errors.New("用户名或密码错误")
		}
		if user.Password != req.Password {
			return nil, errors.New("用户名或密码错误")
		}
		// 如果没有绑定主账号。随机绑定一个主账号
		mainAccount, err := productMainAccountToProductRepo.GetProductMainAccountToUser(user.ID)
		if err != nil && err != gorm.ErrRecordNotFound {
			return nil, errors.New("用户名或密码错误")
		}
		if mainAccount == nil {
			// 查询所有主账号
			mainAccounts, err := gptMainAccountRepo.GetAllAccounts()
			if err != nil {
				return nil, errors.New("用户名或密码错误")
			}
			mainAccountMsg := mainAccounts[rand.Intn(len(mainAccounts))]
			productMainAccountToProductRepo.UpsertUserMainAccountRelation(user.ID, mainAccountMsg.ID)
			// 生成token

		}
		manager, _ := createShareToken.NewTokenManager()
		token, _ := manager.GenerateToken(user.ID, time.Duration(viper.GetInt("expire.share_token"))*time.Second)
		// 重定向到主页
		return &protocol.AdminLoginResp{
			Token:       token,
			Username:    req.Username,
			IsRedirect:  true,
			RedirectUrl: "http://" + viper.GetString("domain."+viper.GetString("env")) + "?share_token=" + token,
		}, nil
	}

	// 生成简单的token（实际项目中应该使用JWT）
	token := "admin_token_" + time.Now().Format("20060102150405")

	return &protocol.AdminLoginResp{
		Token:      token,
		Username:   req.Username,
		IsRedirect: false,
	}, nil
}

// 获取状态文本
func getStatusText(status int) string {
	switch status {
	case gptMainAccountModel.AccountStatusNormal:
		return "正常"
	case gptMainAccountModel.AccountStatusOffline:
		return "手动下线"
	case gptMainAccountModel.AccountStatusError:
		return "异常下线"
	case gptMainAccountModel.AccountStatusDayLimit:
		return "当日已达上限"
	case gptMainAccountModel.AccountStatusCycleLimit:
		return "周期已达上限"
	case gptMainAccountModel.AccountStatusTokenExpired:
		return "token过期"
	default:
		return "未知状态"
	}
}

// 遮盖token显示
func maskToken(token string) string {
	if len(token) <= 10 {
		return token
	}
	return token[:5] + "..." + token[len(token)-5:]
}
