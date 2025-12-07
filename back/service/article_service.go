package service

import (
	"back/model"
	"back/repository"

	"gorm.io/gorm"
)

// ArticleService 文章服务层
type ArticleService struct {
	articleRepo *repository.ArticleRepository
}

// NewArticleService 构造函数
func NewArticleService(articleRepo *repository.ArticleRepository) *ArticleService {
	return &ArticleService{
		articleRepo: articleRepo,
	}
}

// GetArticleList 获取文章列表
func (s *ArticleService) GetArticleList(page, size int, keyword string) ([]*model.Article, int64, string) {
	articles, total, err := s.articleRepo.GetList(page, size, keyword)
	if err != nil {
		return nil, 0, "获取文章列表失败：" + err.Error()
	}
	return articles, total, "success"
}

// GetArticleDetail 获取文章详情
func (s *ArticleService) GetArticleDetail(id uint) (*model.Article, string) {
	article, err := s.articleRepo.GetByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, "文章不存在"
		}
		return nil, "获取文章详情失败：" + err.Error()
	}
	// 阅读量+1
	s.articleRepo.IncrViewCount(id)
	return article, "success"
}

// PublishArticle 发布文章
func (s *ArticleService) PublishArticle(article *model.Article) (bool, string) {
	err := s.articleRepo.Create(article)
	if err != nil {
		return false, "发布文章失败：" + err.Error()
	}
	return true, "发布文章成功"
}

// UpdateArticle 编辑文章
func (s *ArticleService) UpdateArticle(id, userId uint, article *model.Article) (bool, string) {
	// 检查文章是否属于当前用户
	existArticle, err := s.articleRepo.GetByID(id)
	if err != nil {
		return false, "文章不存在"
	}
	if existArticle.UserID != userId {
		return false, "无权编辑他人文章"
	}

	// 更新文章
	err = s.articleRepo.Update(id, article)
	if err != nil {
		return false, "编辑文章失败：" + err.Error()
	}
	return true, "编辑文章成功"
}

// DeleteArticle 删除文章
func (s *ArticleService) DeleteArticle(id, userId uint) (bool, string) {
	// 检查文章归属
	existArticle, err := s.articleRepo.GetByID(id)
	if err != nil {
		return false, "文章不存在"
	}
	if existArticle.UserID != userId {
		return false, "无权删除他人文章"
	}

	// 删除文章
	err = s.articleRepo.Delete(id)
	if err != nil {
		return false, "删除文章失败：" + err.Error()
	}
	return true, "删除文章成功"
}