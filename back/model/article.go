// model/article.go
package model

import "gorm.io/gorm"

// Article 文章模型（核心：匿名嵌入gorm.Model，ID由GORM自动管理）
type Article struct {
	gorm.Model // 关键：这行必须存在，且是匿名嵌入（无字段名）
	Title      string    `gorm:"size:100;not null;index" json:"title"`
	Content    string    `gorm:"type:text;not null" json:"content"`
	CoverImg   string    `gorm:"size:255" json:"cover_img"`
	CategoryID uint      `gorm:"not null;index" json:"category_id"`
	UserID     uint      `gorm:"not null;index" json:"user_id"`
	ViewCount  int       `gorm:"default:0" json:"view_count"`
	IsPublish  bool      `gorm:"default:true" json:"is_publish"`
	Category   Category  `gorm:"foreignKey:CategoryID" json:"category"`
	User       User      `gorm:"foreignKey:UserID" json:"user"`
}