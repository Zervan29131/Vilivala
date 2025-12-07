package model

import (
	"time"
)

<<<<<<< HEAD
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
=======
// Article 对应articles表
type Article struct {
	gorm.Model
	Title      string    `gorm:"size:100;not null;index" json:"title"`   // 标题（索引）
	Content    string    `gorm:"type:text;not null" json:"content"`      // 内容（富文本）
	CoverImg   string    `gorm:"size:255" json:"cover_img"`              // 封面图URL
	CategoryID uint      `gorm:"not null;index" json:"category_id"`      // 分类ID（外键）
	UserID     uint      `gorm:"not null;index" json:"user_id"`          // 作者ID（外键）
	ViewCount  int       `gorm:"default:0" json:"view_count"`            // 阅读量
	IsPublish  bool      `gorm:"default:true" json:"is_publish"`         // 是否发布
	Category   Category  `gorm:"foreignKey:CategoryID" json:"category"`  // 关联分类
	User       User      `gorm:"foreignKey:UserID" json:"user"`          // 关联作者
}

// Category 对应categories表
type Category struct {
	gorm.Model
	Name     string     `gorm:"size:20;not null;unique" json:"name"` // 分类名（唯一）
	Articles []Article  `gorm:"foreignKey:CategoryID" json:"articles"` // 关联文章
>>>>>>> parent of 48535c3 (add)
}