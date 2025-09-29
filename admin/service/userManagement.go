package service

import (
	"errors"
	"time"

	"gpt_mirror/admin/protocol"
	gptMainAccountModel "gpt_mirror/models/gptMainAccount"
	"gpt_mirror/models/user"
	gptMainAccountRepo "gpt_mirror/repo/gptMainAccount"
	productMainAccountToUserRepo "gpt_mirror/repo/producctMainAccountToProduct"
	userRepo "gpt_mirror/repo/user"
)

// 创建用户
func CreateUser(req protocol.CreateUserReq) error {
	// 检查用户名是否已存在
	existingUser, err := userRepo.GetUserMsg(map[string]interface{}{"name": req.Name})
	if err == nil && existingUser != nil {
		return errors.New("用户名已存在")
	}

	// 如果指定了主账号ID，检查主账号是否存在
	if req.MainAccountID > 0 {
		_, err := gptMainAccountRepo.GetAccountByID(req.MainAccountID)
		if err != nil {
			return errors.New("指定的主账号不存在")
		}
	}

	// 创建用户
	userData := user.User{
		Name:        req.Name,
		Password:    req.Password,
		ComboID:     req.ComboID,
		CreatedTime: time.Now(),
		UpdatedTime: time.Now(),
	}

	err = userRepo.CreateUser(userData)
	if err != nil {
		return err
	}

	// 获取创建的用户ID
	createdUser, err := userRepo.GetUserMsg(map[string]interface{}{"name": req.Name})
	if err != nil {
		return err
	}

	// 如果指定了主账号，创建关联关系
	if req.MainAccountID > 0 {
		err = productMainAccountToUserRepo.UpsertUserMainAccountRelation(createdUser.ID, req.MainAccountID)
		if err != nil {
			return err
		}
	}

	return nil
}

// 更新用户
func UpdateUser(req protocol.UpdateUserReq) error {
	// 检查用户是否存在
	_, err := userRepo.GetUserById(req.ID)
	if err != nil {
		return errors.New("用户不存在")
	}

	// 如果要更新用户名，检查是否与其他用户重复
	if req.Name != "" {
		existingUser, err := userRepo.GetUserMsg(map[string]interface{}{"name": req.Name})
		if err == nil && existingUser != nil && existingUser.ID != req.ID {
			return errors.New("用户名已存在")
		}
	}

	// 如果指定了主账号ID，检查主账号是否存在
	if req.MainAccountID > 0 {
		_, err := gptMainAccountRepo.GetAccountByID(req.MainAccountID)
		if err != nil {
			return errors.New("指定的主账号不存在")
		}
	}

	// 更新用户基本信息
	updateData := make(map[string]interface{})
	if req.Name != "" {
		updateData["name"] = req.Name
	}
	if req.Password != "" {
		updateData["password"] = req.Password
	}
	if req.ComboID != 0 {
		updateData["combo_id"] = req.ComboID
	}
	updateData["update_time"] = time.Now()

	err = userRepo.UpdateUserMsg(map[string]interface{}{"id": req.ID}, updateData)
	if err != nil {
		return err
	}

	// 更新主账号关联关系
	err = productMainAccountToUserRepo.UpsertUserMainAccountRelation(req.ID, req.MainAccountID)
	if err != nil {
		return err
	}

	return nil
}

// 删除用户
func DeleteUser(req protocol.DeleteUserReq) error {
	// 检查用户是否存在
	_, err := userRepo.GetUserById(req.ID)
	if err != nil {
		return errors.New("用户不存在")
	}

	// 删除用户主账号关联关系
	err = productMainAccountToUserRepo.DeleteUserMainAccountRelation(req.ID)
	if err != nil {
		return err
	}

	// 删除用户
	err = userRepo.DeleteUserById(req.ID)
	if err != nil {
		return err
	}

	return nil
}

// 获取用户列表
func GetUserList(req protocol.GetUserListReq) (*protocol.UserListResp, error) {
	// 设置默认分页参数
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}

	// 获取总数和用户列表
	total, users, err := userRepo.GetUsersWithMainAccountAndCount(req.Name, req.ComboID, req.Page, req.PageSize)
	if err != nil {
		return nil, err
	}

	// 转换为响应格式
	var list []protocol.UserResp
	for _, user := range users {
		userResp := protocol.UserResp{
			ID:      user.ID,
			Name:    user.Name,
			ComboID: user.ComboID,
		}

		// 处理时间字段
		if !user.CreateTime.IsZero() {
			userResp.CreatedTime = user.CreateTime
		}
		if !user.UpdateTime.IsZero() {
			userResp.UpdatedTime = user.UpdateTime
		}

		// 处理主账号关联信息
		if user.MainAccountID != 0 {
			userResp.MainAccountID = user.MainAccountID
		}
		if user.MainAccountToken != "" {
			token := user.MainAccountToken
			userResp.MainAccountToken = maskToken(token)
		}

		list = append(list, userResp)
	}

	return &protocol.UserListResp{
		List:  list,
		Total: total,
		Page:  req.Page,
		Size:  req.PageSize,
	}, nil
}

// 获取单个用户
func GetUser(req protocol.GetUserReq) (*protocol.UserResp, error) {
	user, err := userRepo.GetUserById(req.ID)
	if err != nil {
		return nil, errors.New("用户不存在")
	}

	userResp := &protocol.UserResp{
		ID:          user.ID,
		Name:        user.Name,
		ComboID:     user.ComboID,
		CreatedTime: user.CreatedTime,
		UpdatedTime: user.UpdatedTime,
	}

	// 获取用户关联的主账号信息
	relation, err := productMainAccountToUserRepo.GetProductMainAccountToUser(user.ID)
	if err == nil && relation != nil {
		userResp.MainAccountID = relation.MainAccountID

		// 获取主账号token
		account, err := gptMainAccountRepo.GetAccountByID(relation.MainAccountID)
		if err == nil {
			userResp.MainAccountToken = maskToken(account.Token)
		}
	}

	return userResp, nil
}

// 设置用户主账号关联
func SetUserMainAccount(req protocol.SetUserMainAccountReq) error {
	// 检查用户是否存在
	_, err := userRepo.GetUserById(req.UserID)
	if err != nil {
		return errors.New("用户不存在")
	}

	// 检查主账号是否存在
	_, err = gptMainAccountRepo.GetAccountByID(req.MainAccountID)
	if err != nil {
		return errors.New("主账号不存在")
	}

	// 设置关联关系
	err = productMainAccountToUserRepo.UpsertUserMainAccountRelation(req.UserID, req.MainAccountID)
	if err != nil {
		return err
	}

	return nil
}

// 获取可用的主账号列表
func GetAvailableMainAccounts() ([]protocol.AvailableMainAccountResp, error) {
	// 获取所有正常状态的主账号
	accounts, err := gptMainAccountRepo.GetAccountsByStatus(map[string]interface{}{
		"status": gptMainAccountModel.AccountStatusNormal,
	})
	if err != nil {
		return nil, err
	}

	var list []protocol.AvailableMainAccountResp
	for _, account := range accounts {
		list = append(list, protocol.AvailableMainAccountResp{
			ID:         account.ID,
			Token:      maskToken(account.Token),
			Status:     account.Status,
			StatusText: getStatusText(account.Status),
		})
	}

	return list, nil
}
