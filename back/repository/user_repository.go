package repository

import (
	"back/model"

	"gorm.io/gorm"
)

// UserRepository结构体（依赖db）
type UserRepository struct {
	db *gorm.DB
}

// 构造函数
func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

// 检查用户名是否存在
func (r *UserRepository) IsUsernameExist(username string) bool {
	var count int64
	r.db.Model(&model.User{}).Where("username = ?", username).Count(&count)
	return count > 0
}

// 创建用户
func (r *UserRepository) CreateUser(user *model.User) error {
	return r.db.Create(user).Error
}

// 通过用户名查用户
func (r *UserRepository) GetUserByUsername(username string) (*model.User, error) {
	var user model.User
	err := r.db.Where("username = ?", username).First(&user).Error
	return &user, err
}

// 通过ID查用户
func (r *UserRepository) GetUserByID(userID uint) (*model.User, error) {
	var user model.User
<<<<<<< HEAD
	err := r.db.Where("id = ?", userID).First(&user).Error
	return &user, err
}

// 更新密码
func (r *UserRepository) UpdateUserPassword(userID uint, newPwd string) error {
	return r.db.Model(&model.User{}).Where("id = ?", userID).Update("password", newPwd).Error
=======
	err := r.db.Where("id = ?", id).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
>>>>>>> parent of 48535c3 (add)
}