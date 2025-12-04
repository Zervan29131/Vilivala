package controller

import (
	"net/http"

	"back/service"

	"github.com/gin-gonic/gin"
)

type CategoryController struct {
	categoryService *service.CategoryService
}

func NewCategoryController(s *service.CategoryService) *CategoryController {
	return &CategoryController{categoryService: s}
}

// List 获取所有分类
// @Summary 获取分类列表
// @Produce json
// @Success 200 {object} gin.H{code:int, msg:string, data:[]model.Category}
// @Router /api/v1/category/list [get]
func (c *CategoryController) List(ctx *gin.Context) {
	categories, err := c.categoryService.List()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "获取分类失败：" + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "获取成功",
		"data": categories,
	})
}