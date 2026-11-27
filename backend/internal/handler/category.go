package handler

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"vilivala-backend/config"
	"vilivala-backend/internal/model"
)

type catOut struct {
	model.Category
	Depth      int        `json:"depth"`
	PostCount  int64      `json:"post_count"`
	TotalCount int64      `json:"total_count"`
	Path       string     `json:"path"`
	Children   *[]*catOut `json:"children,omitempty"`
}

// ListCategories GET /api/categories[?tree=1]
func ListCategories(c *gin.Context) {
	ix := loadCats()
	var rows []struct {
		CategoryID uint
		N          int64
	}
	config.DB.Model(&model.Post{}).Select("category_id, COUNT(*) AS n").
		Where("status = ? AND category_id IS NOT NULL", "published").Group("category_id").Scan(&rows)
	direct := map[uint]int64{}
	for _, r := range rows {
		direct[r.CategoryID] = r.N
	}

	build := func(id uint) *catOut {
		var total int64
		for _, d := range ix.descendants(id) {
			total += direct[d]
		}
		return &catOut{Category: *ix.byID[id], Depth: ix.depth(id), PostCount: direct[id],
			TotalCount: total, Path: ix.pathString(id)}
	}

	if c.Query("tree") == "1" {
		var walk func(ids []uint) []*catOut
		walk = func(ids []uint) []*catOut {
			out := []*catOut{}
			for _, id := range ids {
				n := build(id)
				kids := walk(ix.children[id])
				n.Children = &kids
				out = append(out, n)
			}
			return out
		}
		OK(c, walk(ix.roots))
		return
	}

	// flat list ordered by sort, id
	out := []*catOut{}
	var all []model.Category
	config.DB.Order("sort ASC, id ASC").Find(&all)
	for _, cat := range all {
		out = append(out, build(cat.ID))
	}
	OK(c, out)
}

type categoryPayload struct {
	Name     string `json:"name" binding:"required"`
	Slug     string `json:"slug" binding:"required"`
	ParentID uint   `json:"parent_id"`
	Sort     int    `json:"sort"`
}

func categoryErr(c *gin.Context, err error) {
	if isDup(err) {
		msg := "同级分类名称已存在"
		if containsSlug(err) {
			msg = "slug 已存在"
		}
		Fail(c, http.StatusBadRequest, 400, msg)
		return
	}
	Fail(c, http.StatusInternalServerError, 500, "操作失败")
}

func containsSlug(err error) bool { return strings.Contains(err.Error(), "slug") }

// AdminCreateCategory POST /api/admin/categories
func AdminCreateCategory(c *gin.Context) {
	var req categoryPayload
	if err := c.ShouldBindJSON(&req); err != nil {
		Fail(c, http.StatusBadRequest, 400, "name 和 slug 为必填")
		return
	}
	ix := loadCats()
	if req.ParentID != 0 {
		if ix.byID[req.ParentID] == nil {
			Fail(c, http.StatusBadRequest, 400, "父分类不存在")
			return
		}
		if ix.depth(req.ParentID)+1 > maxCategoryDepth {
			Fail(c, http.StatusBadRequest, 400, "分类最多三级")
			return
		}
	}
	m := model.Category{Name: req.Name, Slug: req.Slug, ParentID: req.ParentID, Sort: req.Sort}
	if err := config.DB.Create(&m).Error; err != nil {
		categoryErr(c, err)
		return
	}
	OK(c, m)
}

// AdminUpdateCategory PUT /api/admin/categories/:id
func AdminUpdateCategory(c *gin.Context) {
	var req categoryPayload
	if err := c.ShouldBindJSON(&req); err != nil {
		Fail(c, http.StatusBadRequest, 400, "name 和 slug 为必填")
		return
	}
	id64, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	id := uint(id64)
	ix := loadCats()
	cur := ix.byID[id]
	if cur == nil {
		Fail(c, http.StatusNotFound, 404, "分类不存在")
		return
	}
	if req.ParentID != 0 {
		if ix.byID[req.ParentID] == nil {
			Fail(c, http.StatusBadRequest, 400, "父分类不存在")
			return
		}
		for _, d := range ix.descendants(id) {
			if d == req.ParentID {
				Fail(c, http.StatusBadRequest, 400, "不能移动到自身或子分类下")
				return
			}
		}
		if ix.depth(req.ParentID)+ix.height(id) > maxCategoryDepth {
			Fail(c, http.StatusBadRequest, 400, "分类最多三级")
			return
		}
	}
	err := config.DB.Model(cur).Updates(map[string]interface{}{
		"name": req.Name, "slug": req.Slug, "parent_id": req.ParentID, "sort": req.Sort,
	}).Error
	if err != nil {
		categoryErr(c, err)
		return
	}
	OK(c, nil)
}

// AdminDeleteCategory DELETE /api/admin/categories/:id
func AdminDeleteCategory(c *gin.Context) {
	id64, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	id := uint(id64)
	ix := loadCats()
	if ix.byID[id] == nil {
		Fail(c, http.StatusNotFound, 404, "分类不存在")
		return
	}
	var n int64
	config.DB.Model(&model.Post{}).Where("category_id = ?", id).Count(&n)
	if n > 0 || len(ix.children[id]) > 0 {
		Fail(c, http.StatusBadRequest, 400, "请先移除子分类和文章")
		return
	}
	if err := config.DB.Delete(&model.Category{}, id).Error; err != nil {
		Fail(c, http.StatusInternalServerError, 500, "删除失败")
		return
	}
	OK(c, nil)
}

// ---- tags ----

type namePayload struct {
	Name string `json:"name" binding:"required"`
	Slug string `json:"slug" binding:"required"`
}

// ListTags GET /api/tags
func ListTags(c *gin.Context) {
	var list []struct {
		model.Tag
		PostCount int64 `json:"post_count"`
	}
	err := config.DB.Model(&model.Tag{}).
		Select("tags.*, (SELECT COUNT(*) FROM post_tags pt JOIN posts p ON p.id = pt.post_id WHERE pt.tag_id = tags.id AND p.status = 'published') AS post_count").
		Order("tags.id").Scan(&list).Error
	if err != nil {
		Fail(c, http.StatusInternalServerError, 500, "查询失败")
		return
	}
	if list == nil {
		OK(c, []int{})
		return
	}
	OK(c, list)
}

func dupOrFail(c *gin.Context, err error) {
	if isDup(err) {
		Fail(c, http.StatusBadRequest, 400, "名称或 slug 已存在")
	} else {
		Fail(c, http.StatusInternalServerError, 500, "操作失败")
	}
}

// AdminCreateTag POST /api/admin/tags
func AdminCreateTag(c *gin.Context) {
	var req namePayload
	if err := c.ShouldBindJSON(&req); err != nil {
		Fail(c, http.StatusBadRequest, 400, "name 和 slug 为必填")
		return
	}
	m := model.Tag{Name: req.Name, Slug: req.Slug}
	if err := config.DB.Create(&m).Error; err != nil {
		dupOrFail(c, err)
		return
	}
	OK(c, m)
}

// AdminUpdateTag PUT /api/admin/tags/:id
func AdminUpdateTag(c *gin.Context) {
	var req namePayload
	if err := c.ShouldBindJSON(&req); err != nil {
		Fail(c, http.StatusBadRequest, 400, "name 和 slug 为必填")
		return
	}
	var m model.Tag
	if err := config.DB.First(&m, c.Param("id")).Error; err != nil {
		Fail(c, http.StatusNotFound, 404, "标签不存在")
		return
	}
	if err := config.DB.Model(&m).Updates(map[string]interface{}{"name": req.Name, "slug": req.Slug}).Error; err != nil {
		dupOrFail(c, err)
		return
	}
	OK(c, nil)
}

// AdminDeleteTag DELETE /api/admin/tags/:id
func AdminDeleteTag(c *gin.Context) {
	var m model.Tag
	if err := config.DB.First(&m, c.Param("id")).Error; err != nil {
		Fail(c, http.StatusNotFound, 404, "标签不存在")
		return
	}
	err := config.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&m).Association("Posts").Clear(); err != nil {
			return err
		}
		return tx.Delete(&m).Error
	})
	if err != nil {
		Fail(c, http.StatusInternalServerError, 500, "删除失败")
		return
	}
	OK(c, nil)
}
