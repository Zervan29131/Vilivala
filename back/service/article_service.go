// service/article_service.go
package service

import (
	"back/model"
	"back/repository"
	"errors"

	"gorm.io/gorm"
)

type ArticleService struct {
	articleRepo  *repository.ArticleRepository
	categoryRepo *repository.CategoryRepository
}

func NewArticleService(articleRepo *repository.ArticleRepository, categoryRepo *repository.CategoryRepository) *ArticleService {
	return &ArticleService{
		articleRepo:  articleRepo,
		categoryRepo: categoryRepo,
	}
}

// List 分页查询文章
func (s *ArticleService) List(page, size int, keyword string) ([]model.Article, int64, error) {
	return s.articleRepo.List(page, size, keyword)
}

// GetDetail 获取文章详情
func (s *ArticleService) GetDetail(id uint) (*model.Article, error) {
	return s.articleRepo.GetByID(id)
}

// Create 发布文章
func (s *ArticleService) Create(article *model.Article) error {
	_, err := s.categoryRepo.GetByID(article.CategoryID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("分类不存在")
		}
		return err
	}
	return s.articleRepo.Create(article)
}

// Update 编辑文章（核心：接收ID和DTO，不操作结构体ID）
func (s *ArticleService) Update(articleID uint, userID uint, updateDTO interface{}) error {
	// 1. 查询文章是否存在
	existArticle, err := s.articleRepo.GetByID(articleID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("文章不存在")
		}
		return err
	}

	// 2. 权限校验
	if existArticle.UserID != userID {
		return errors.New("无权限修改该文章")
	}

	// 3. 解析更新DTO
	req := updateDTO.(struct {
		Title      string `json:"title"`
		Content    string `json:"content"`
		CoverImg   string `json:"cover_img"`
		CategoryID uint   `json:"category_id"`
		IsPublish  bool   `json:"is_publish"`
	})

	// 4. 校验分类
	_, err = s.categoryRepo.GetByID(req.CategoryID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("分类不存在")
		}
		return err
	}

	// 5. 构造更新数据（Map形式，无ID）
	updateData := map[string]interface{}{
		"title":        req.Title,
		"content":      req.Content,
		"cover_img":    req.CoverImg,
		"category_id":  req.CategoryID,
		"is_publish":   req.IsPublish,
	}

	// 6. 调用Repo更新（通过ID条件更新，无结构体ID）
	return s.articleRepo.UpdateByID(articleID, updateData)
}

// Delete 删除文章
func (s *ArticleService) Delete(id, userID uint) error {
	existArticle, err := s.articleRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("文章不存在")
		}
		return err
	}

	if existArticle.UserID != userID {
		return errors.New("无权限删除该文章")
	}

	return s.articleRepo.Delete(id, userID)
}