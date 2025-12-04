// repository/article_repository.go
package repository

import (
	"back/model"

	"gorm.io/gorm"
)

type ArticleRepository struct {
	db *gorm.DB
}

func NewArticleRepository(db *gorm.DB) *ArticleRepository {
	return &ArticleRepository{db: db}
}

// List 分页查询
func (r *ArticleRepository) List(page, size int, keyword string) ([]model.Article, int64, error) {
	var articles []model.Article
	var total int64

	query := r.db.Model(&model.Article{}).
		Preload("User", func(db *gorm.DB) *gorm.DB {
			return db.Select("id, username")
		}).
		Preload("Category")

	if keyword != "" {
		query = query.Where("title LIKE ? OR content LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = query.
		Order("created_at DESC").
		Offset((page - 1) * size).
		Limit(size).
		Find(&articles).Error

	return articles, total, err
}

// GetByID 根据ID查询
func (r *ArticleRepository) GetByID(id uint) (*model.Article, error) {
	var article model.Article
	err := r.db.
		Preload("User", func(db *gorm.DB) *gorm.DB {
			return db.Select("id, username")
		}).
		Preload("Category").
		Where("id = ?", id).
		First(&article).Error
	if err != nil {
		return nil, err
	}

	r.db.Model(&article).Update("view_count", gorm.Expr("view_count + ?", 1))
	return &article, nil
}

// Create 新增
func (r *ArticleRepository) Create(article *model.Article) error {
	return r.db.Create(article).Error
}

// UpdateByID 根据ID更新（核心：无结构体ID赋值）
func (r *ArticleRepository) UpdateByID(articleID uint, updateData map[string]interface{}) error {
	return r.db.Model(&model.Article{}).
		Where("id = ?", articleID).
		Updates(updateData).Error
}

// Delete 删除
func (r *ArticleRepository) Delete(id uint, userID uint) error {
	return r.db.Where("id = ? AND user_id = ?", id, userID).Delete(&model.Article{}).Error
}