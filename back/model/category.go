package model

import "gorm.io/gorm"

type Category struct {
	gorm.Model // 匿名嵌入gorm.Model
	Name     string    `gorm:"size:20;not null;unique" json:"name"`
	Articles []Article `gorm:"foreignKey:CategoryID" json:"-"`
}