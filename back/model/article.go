package model

import (
	"time"
)

// Article 文章模型
type Article struct {
	ID         uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Title      string    `gorm:"size:100;not null" json:"title"`       // 文章标题
	Content    string    `gorm:"type:text;not null" json:"content"`    // 文章内容
	ViewCount  int       `gorm:"default:0" json:"view_count"`          // 阅读量
	IsPublish  bool      `gorm:"default:true" json:"is_publish"`       // 是否发布
	CategoryID uint      `gorm:"not null" json:"category_id"`          // 分类ID
	UserID     uint      `gorm:"not null" json:"user_id"`              // 作者ID
	CreatedAt  time.Time `gorm:"autoCreateTime" json:"created_at"`     // 创建时间
	UpdatedAt  time.Time `gorm:"autoUpdateTime" json:"updated_at"`     // 更新时间

	// 关联查询
	User     *User     `gorm:"foreignKey:UserID" json:"user"`         // 作者信息
	Category *Category `gorm:"foreignKey:CategoryID" json:"category"` // 分类信息
}