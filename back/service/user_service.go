package service

import (
	"back/middleware"
	"back/model"
	"back/repository"

	"golang.org/x/crypto/bcrypt"
)

// UserService结构体（依赖UserRepository，而非直接依赖db）
type UserService struct {
	userRepo *repository.UserRepository // 注入Repository
}

// 构造函数（参数改为*UserRepository）
func NewUserService(userRepo *repository.UserRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

// 注册逻辑（调用Repository层）
func (s *UserService) Register(username, password, avatar string) (bool, string) {
	// 调用Repository检查用户名
	if s.userRepo.IsUsernameExist(username) {
		return false, "用户名已存在"
	}

	// 密码加密
	hashPwd, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return false, "密码加密失败：" + err.Error()
	}

	// 调用Repository创建用户
	user := &model.User{
		Username: username,
		Password: string(hashPwd),
		Avatar:   avatar,
		Role:     "user",
	}
	if err := s.userRepo.CreateUser(user); err != nil {
		return false, "注册失败：" + err.Error()
	}
<<<<<<< HEAD
	return true, "注册成功"
}

// 登录逻辑
func (s *UserService) Login(username, password string) (*model.User, string, string) {
	// 调用Repository查用户
	user, err := s.userRepo.GetUserByUsername(username)
	if err != nil {
		return nil, "", "用户名不存在"
	}

	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, "", "密码错误"
	}

	// 生成Token
	token, err := middleware.GenerateToken(user.ID, user.Role)
	if err != nil {
		return nil, "", "生成Token失败：" + err.Error()
	}
	return user, token, "登录成功"
}

// 获取用户信息
func (s *UserService) GetUserInfo(userID uint) (*model.User, string) {
	user, err := s.userRepo.GetUserByID(userID)
	if err != nil {
		return nil, "用户信息不存在"
	}
	return user, "获取成功"
}

// （可选）补充修改密码方法（如果需要）
func (s *UserService) ChangePassword(userID uint, oldPwd, newPwd string) (bool, string) {
	user, err := s.userRepo.GetUserByID(userID)
	if err != nil {
		return false, "用户不存在"
	}

	// 验证旧密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(oldPwd)); err != nil {
		return false, "旧密码错误"
	}

	// 加密新密码
	hashNewPwd, err := bcrypt.GenerateFromPassword([]byte(newPwd), bcrypt.DefaultCost)
	if err != nil {
		return false, "密码加密失败"
	}

	// 更新密码
	if err := s.userRepo.UpdateUserPassword(userID, string(hashNewPwd)); err != nil {
		return false, "修改密码失败：" + err.Error()
	}
	return true, "修改密码成功"
=======

	// 4. 构造返回的用户信息（隐藏敏感字段）
	userInfo = map[string]interface{}{
		"id":       user.ID,
		"username": user.Username,
		"avatar":   user.Avatar,
		"role":     user.Role,
	}

	return token, userInfo, nil
>>>>>>> parent of 48535c3 (add)
}