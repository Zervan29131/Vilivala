package handler

import (
	"net/http"
	"regexp"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm/clause"
	"vilivala-backend/config"
	"vilivala-backend/internal/model"
)

var pageSlugRe = regexp.MustCompile(`^[a-z0-9-]{1,64}$`)

// ListPages GET /api/pages and GET /api/admin/pages (no content)
func ListPages(c *gin.Context) {
	var pages []model.Page
	if err := config.DB.Omit("content").Order("id ASC").Find(&pages).Error; err != nil {
		Fail(c, http.StatusInternalServerError, 500, "查询失败")
		return
	}
	out := make([]gin.H, 0, len(pages))
	for _, p := range pages {
		out = append(out, gin.H{"slug": p.Slug, "title": p.Title, "updated_at": p.UpdatedAt})
	}
	OK(c, out)
}

// GetPage GET /api/pages/:slug
func GetPage(c *gin.Context) {
	var p model.Page
	if err := config.DB.Where("slug = ?", c.Param("slug")).First(&p).Error; err != nil {
		Fail(c, http.StatusNotFound, 404, "页面不存在")
		return
	}
	OK(c, gin.H{"slug": p.Slug, "title": p.Title, "content": p.Content, "updated_at": p.UpdatedAt})
}

// AdminPutPage PUT /api/admin/pages/:slug (upsert)
func AdminPutPage(c *gin.Context) {
	slug := c.Param("slug")
	var req struct {
		Title   string `json:"title"`
		Content string `json:"content"`
	}
	if !pageSlugRe.MatchString(slug) {
		Fail(c, http.StatusBadRequest, 400, "slug 只能包含小写字母、数字和连字符（1-64 位）")
		return
	}
	if err := c.ShouldBindJSON(&req); err != nil || len([]rune(req.Title)) < 1 || len([]rune(req.Title)) > 64 {
		Fail(c, http.StatusBadRequest, 400, "title 为必填，1-64 字符")
		return
	}
	p := model.Page{Slug: slug, Title: req.Title, Content: req.Content}
	err := config.DB.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "slug"}},
		DoUpdates: clause.AssignmentColumns([]string{"title", "content", "updated_at"}),
	}).Create(&p).Error
	if err != nil {
		Fail(c, http.StatusInternalServerError, 500, "保存失败")
		return
	}
	OK(c, nil)
}

// AdminDeletePage DELETE /api/admin/pages/:slug
func AdminDeletePage(c *gin.Context) {
	res := config.DB.Where("slug = ?", c.Param("slug")).Delete(&model.Page{})
	if res.Error != nil {
		Fail(c, http.StatusInternalServerError, 500, "删除失败")
		return
	}
	if res.RowsAffected == 0 {
		Fail(c, http.StatusNotFound, 404, "页面不存在")
		return
	}
	OK(c, nil)
}
