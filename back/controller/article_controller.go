package controller

import (
	"net/http"
	"strconv"

	"back/model"
	"back/service"

	"github.com/gin-gonic/gin"
)

// ArticleController 文章控制器
type ArticleController struct {
	articleService *service.ArticleService
}

// NewArticleController 创建控制器实例
func NewArticleController(s *service.ArticleService) *ArticleController {
	return &ArticleController{articleService: s}
}

// List 分页查询文章列表
func (c *ArticleController) List(ctx *gin.Context) {
	pageStr := ctx.DefaultQuery("page", "1")
	sizeStr := ctx.DefaultQuery("size", "10")
	keyword := ctx.Query("keyword")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	size, err := strconv.Atoi(sizeStr)
	if err != nil || size < 1 || size > 50 {
		size = 10
	}

	articles, total, err := c.articleService.List(page, size, keyword)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "获取文章列表失败：" + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "获取成功",
		"data": gin.H{
			"list":  articles,
			"total": total,
		},
	})
}

// GetDetail 获取文章详情
func (c *ArticleController) GetDetail(ctx *gin.Context) {
	idStr := ctx.Param("id")
	articleID, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "无效的文章ID：" + err.Error(),
		})
		return
	}

	article, err := c.articleService.GetDetail(uint(articleID))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "获取文章详情失败：" + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "获取成功",
		"data": article,
	})
}

// Create 发布文章
func (c *ArticleController) Create(ctx *gin.Context) {
	// 定义请求结构体（无ID）
	type CreateReq struct {
		Title      string `json:"title" binding:"required,max=100"`
		Content    string `json:"content" binding:"required"`
		CoverImg   string `json:"cover_img"`
		CategoryID uint   `json:"category_id" binding:"required"`
		IsPublish  bool   `json:"is_publish" default:"true"`
	}

	var req CreateReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "参数错误：" + err.Error(),
		})
		return
	}

	// 获取当前用户ID
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"code": 401,
			"msg":  "请先登录",
		})
		return
	}

	// 构造文章对象（绝对无ID字段）
	article := &model.Article{
		Title:      req.Title,
		Content:    req.Content,
		CoverImg:   req.CoverImg,
		CategoryID: req.CategoryID,
		UserID:     userID.(uint),
		IsPublish:  req.IsPublish,
	}

	if err := c.articleService.Create(article); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "发布文章失败：" + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "发布成功",
	})
}

// Update 编辑文章
func (c *ArticleController) Update(ctx *gin.Context) {
	// 1. 获取文章ID（仅作为参数传递，不赋值给结构体）
	idStr := ctx.Param("id")
	articleID, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "无效的文章ID：" + err.Error(),
		})
		return
	}

	// 2. 定义更新请求（无ID）
	type UpdateReq struct {
		Title      string `json:"title" binding:"required,max=100"`
		Content    string `json:"content" binding:"required"`
		CoverImg   string `json:"cover_img"`
		CategoryID uint   `json:"category_id" binding:"required"`
		IsPublish  bool   `json:"is_publish" default:"true"`
	}

	var req UpdateReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "参数错误：" + err.Error(),
		})
		return
	}

	// 3. 获取当前用户ID
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"code": 401,
			"msg":  "请先登录",
		})
		return
	}

	// 4. 调用服务层（传递ID和更新数据，无结构体ID赋值）
	err = c.articleService.Update(
		uint(articleID),
		userID.(uint),
		req,
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "编辑文章失败：" + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "编辑成功",
	})
}

// Delete 删除文章
func (c *ArticleController) Delete(ctx *gin.Context) {
	idStr := ctx.Param("id")
	articleID, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "无效的文章ID：" + err.Error(),
		})
		return
	}

	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"code": 401,
			"msg":  "请先登录",
		})
		return
	}

	err = c.articleService.Delete(uint(articleID), userID.(uint))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "删除文章失败：" + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "删除成功",
	})
}