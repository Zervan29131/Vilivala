package controller

import (
	"back/model"
	"back/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// ArticleController 文章控制器
type ArticleController struct {
	articleService *service.ArticleService
}

// NewArticleController 构造函数（依赖注入ArticleService）
func NewArticleController(articleService *service.ArticleService) *ArticleController {
	return &ArticleController{
		articleService: articleService,
	}
}

// ========== 实现路由中调用的所有方法 ==========

// GetArticleList 获取文章列表（分页+搜索）
// 解决：router中调用的GetArticleList方法，此处必须首字母大写且名称完全一致
func (c *ArticleController) GetArticleList(ctx *gin.Context) {
	// 1. 获取分页/搜索参数
	pageStr := ctx.DefaultQuery("page", "1")
	sizeStr := ctx.DefaultQuery("size", "10")
	keyword := ctx.Query("keyword")

	// 2. 转换参数类型
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}
	size, err := strconv.Atoi(sizeStr)
	if err != nil || size < 1 {
		size = 10
	}

	// 3. 调用Service层获取列表
	articles, total, msg := c.articleService.GetArticleList(page, size, keyword)
	if articles == nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  msg,
			"data": nil,
		})
		return
	}

	// 4. 返回数据
	ctx.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "获取文章列表成功",
		"data": gin.H{
			"list":  articles,
			"total": total,
			"page":  page,
			"size":  size,
		},
	})
}

// GetArticleDetail 获取文章详情
func (c *ArticleController) GetArticleDetail(ctx *gin.Context) {
	// 1. 获取文章ID
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "文章ID格式错误",
			"data": nil,
		})
		return
	}

	// 2. 调用Service层获取详情
	article, msg := c.articleService.GetArticleDetail(uint(id))
	if article == nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"code": 404,
			"msg":  msg,
			"data": nil,
		})
		return
	}

	// 3. 返回数据
	ctx.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "获取文章详情成功",
		"data": article,
	})
}

// PublishArticle 发布文章
func (c *ArticleController) PublishArticle(ctx *gin.Context) {
	// 1. 定义请求参数结构体
	type PublishRequest struct {
		Title      string `json:"title" binding:"required,min=1,max=100"`
		Content    string `json:"content" binding:"required"`
		CategoryID uint   `json:"category_id" binding:"required"`
		IsPublish  bool   `json:"is_publish" default:"true"`
	}

	// 2. 绑定参数
	var req PublishRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "参数错误：" + err.Error(),
			"data": nil,
		})
		return
	}

	// 3. 获取当前登录用户ID（从JWT中间件存入的上下文）
	userId, _ := ctx.Get("userID")

	// 4. 构造文章模型
	article := &model.Article{
		Title:      req.Title,
		Content:    req.Content,
		CategoryID: req.CategoryID,
		UserID:     userId.(uint),
		IsPublish:  req.IsPublish,
	}

	// 5. 调用Service层发布
	success, msg := c.articleService.PublishArticle(article)
	if success {
		ctx.JSON(http.StatusOK, gin.H{
			"code": 200,
			"msg":  msg,
			"data": nil,
		})
	} else {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  msg,
			"data": nil,
		})
	}
}

// UpdateArticle 编辑文章
func (c *ArticleController) UpdateArticle(ctx *gin.Context) {
	// 1. 获取文章ID
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "文章ID格式错误",
			"data": nil,
		})
		return
	}

	// 2. 绑定更新参数
	type UpdateRequest struct {
		Title      string `json:"title" binding:"required,min=1,max=100"`
		Content    string `json:"content" binding:"required"`
		CategoryID uint   `json:"category_id" binding:"required"`
		IsPublish  bool   `json:"is_publish"`
	}
	var req UpdateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "参数错误：" + err.Error(),
			"data": nil,
		})
		return
	}

	// 3. 获取当前用户ID
	userId, _ := ctx.Get("userID")

	// 4. 调用Service层更新
	success, msg := c.articleService.UpdateArticle(uint(id), userId.(uint), &model.Article{
		Title:      req.Title,
		Content:    req.Content,
		CategoryID: req.CategoryID,
		IsPublish:  req.IsPublish,
	})
	if success {
		ctx.JSON(http.StatusOK, gin.H{
			"code": 200,
			"msg":  msg,
			"data": nil,
		})
	} else {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  msg,
			"data": nil,
		})
	}
}

// DeleteArticle 删除文章
func (c *ArticleController) DeleteArticle(ctx *gin.Context) {
	// 1. 获取文章ID
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "文章ID格式错误",
			"data": nil,
		})
		return
	}

	// 2. 获取当前用户ID
	userId, _ := ctx.Get("userID")

	// 3. 调用Service层删除
	success, msg := c.articleService.DeleteArticle(uint(id), userId.(uint))
	if success {
		ctx.JSON(http.StatusOK, gin.H{
			"code": 200,
			"msg":  msg,
			"data": nil,
		})
	} else {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  msg,
			"data": nil,
		})
	}
}