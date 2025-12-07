package repository

import (
	"back/model"

	"gorm.io/gorm"
)

// ArticleRepository 文章仓库层
type ArticleRepository struct {
	db *gorm.DB
}

// NewArticleRepository 构造函数
func NewArticleRepository(db *gorm.DB) *ArticleRepository {
	return &ArticleRepository{
		db: db,
	}
}

// GetList 获取文章列表（分页+搜索）
func (r *ArticleRepository) GetList(page, size int, keyword string) ([]*model.Article, int64, error) {
	var articles []*model.Article
	var total int64

	// 构建查询条件
	query := r.db.Model(&model.Article{})
	if keyword != "" {
		query = query.Where("title LIKE ? OR content LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}

	// 统计总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (page - 1) * size
	if err := query.Preload("User").Preload("Category").Offset(offset).Limit(size).Order("created_at DESC").Find(&articles).Error; err != nil {
		return nil, 0, err
	}

	return articles, total, nil
}

// GetByID 通过ID获取文章
func (r *ArticleRepository) GetByID(id uint) (*model.Article, error) {
	var article model.Article
	err := r.db.Preload("User").Preload("Category").Where("id = ?", id).First(&article).Error
	return &article, err
}

// Create 创建文章
func (r *ArticleRepository) Create(article *model.Article) error {
	return r.db.Create(article).Error
}

// Update 更新文章
func (r *ArticleRepository) Update(id uint, article *model.Article) error {
	return r.db.Model(&model.Article{}).Where("id = ?", id).Updates(map[string]interface{}{
		"title":       article.Title,
		"content":     article.Content,
		"category_id": article.CategoryID,
		"is_publish":  article.IsPublish,
	}).Error
}

// Delete 删除文章
func (r *ArticleRepository) Delete(id uint) error {
	return r.db.Where("id = ?", id).Delete(&model.Article{}).Error
}

// IncrViewCount 阅读量+1
func (r *ArticleRepository) IncrViewCount(id uint) error {
	return r.db.Model(&model.Article{}).Where("id = ?", id).Update("view_count", gorm.Expr("view_count + ?", 1)).Error
}