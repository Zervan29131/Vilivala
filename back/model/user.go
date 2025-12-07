package model

import (
	"back/util" // 后续引入加密工具

	"gorm.io/gorm"
)

// User 对应MySQL users表（GORM自动小写+复数）
type User struct {
	gorm.Model           // 内置字段：ID(uint)、CreatedAt、UpdatedAt、DeletedAt（软删除）
	Username string      `gorm:"size:20;not null;unique" json:"username"` // 用户名（唯一）
	Password string      `gorm:"size:100;not null" json:"-"`             // 密码（加密，返回前端时隐藏）
	Avatar   string      `gorm:"size:255" json:"avatar"`                 // 头像URL
	Role     string      `gorm:"size:10;default:'user'" json:"role"`     // 角色（admin/user）
}

// BeforeSave GORM钩子：保存用户前自动加密密码
func (u *User) BeforeSave(tx *gorm.DB) error {
	if u.Password != "" { // 密码有值时才加密（避免更新时覆盖）
		hashPwd, err := util.BcryptEncrypt(u.Password)
		if err != nil {
			return err
		}
		u.Password = hashPwd
	}
	return nil
}