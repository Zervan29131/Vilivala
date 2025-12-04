package service

import (
	"back/model"
	"back/repository"
)

type CategoryService struct {
	categoryRepo *repository.CategoryRepository
}

func NewCategoryService(repo *repository.CategoryRepository) *CategoryService {
	return &CategoryService{categoryRepo: repo}
}

// List 获取所有分类
func (s *CategoryService) List() ([]model.Category, error) {
	return s.categoryRepo.List()
}