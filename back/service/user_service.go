package service

import (
	"back/middleware"
	"back/model"
	"back/repository"
	"back/util"
	"errors"

	"gorm.io/gorm"
)

// UserService 用户业务逻辑层
type UserService struct {
	userRepo *repository.UserRepository // 依赖数据访问层
}

// NewUserService 初始化UserService
func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{userRepo: repo}
}

// Register 注册业务逻辑
func (s *UserService) Register(user *model.User) error {
	// 1. 校验用户名是否存在
	existUser, _ := s.userRepo.GetByUsername(user.Username)
	if existUser != nil {
		return errors.New("用户名已存在")
	}
	// 2. 新增用户（密码加密在model钩子中处理）
	return s.userRepo.Create(user)
}

// Login 登录业务逻辑：验证密码 + 生成JWT
func (s *UserService) Login(username, password string) (token string, userInfo map[string]interface{}, err error) {
	// 1. 查询用户
	user, err := s.userRepo.GetByUsername(username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", nil, errors.New("用户名/密码错误")
		}
		return "", nil, err
	}

	// 2. 校验密码
	if !util.BcryptCompare(user.Password, password) {
		return "", nil, errors.New("用户名/密码错误")
	}

	// 3. 生成JWT令牌
	token, err = middleware.GenerateToken(user.ID, user.Role)
	if err != nil {
		return "", nil, errors.New("令牌生成失败")
	}

	// 4. 构造返回的用户信息（隐藏敏感字段）
	userInfo = map[string]interface{}{
		"id":       user.ID,
		"username": user.Username,
		"avatar":   user.Avatar,
		"role":     user.Role,
	}

	return token, userInfo, nil
}
// GetInfo 根据ID获取用户信息
func (s *UserService) GetInfo(id uint) (*model.User, error) {
	return s.userRepo.GetByID(id)
}

// ChangePassword 修改密码
func (s *UserService) ChangePassword(userID uint, oldPwd, newPwd string) error {
	// 查询用户
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return err
	}

	// 校验原密码
	if !util.BcryptCompare(user.Password, oldPwd) {
		return errors.New("原密码错误")
	}

	// 加密新密码
	hashPwd, err := util.BcryptEncrypt(newPwd)
	if err != nil {
		return err
	}

	// 更新密码
	user.Password = hashPwd
	return s.userRepo.UpdatePassword(user)
}