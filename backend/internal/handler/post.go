package handler

import (
	"errors"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"vilivala-backend/config"
	"vilivala-backend/internal/model"
)

const postOrder = "is_top DESC, created_at DESC, id DESC"

func pageParams(c *gin.Context) (page, size int) {
	page, _ = strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ = strconv.Atoi(c.DefaultQuery("page_size", "10"))
	if page < 1 {
		page = 1
	}
	if size < 1 {
		size = 10
	}
	if size > 50 {
		size = 50
	}
	return
}

func listPosts(c *gin.Context, publishedOnly bool) {
	page, size := pageParams(c)
	q := config.DB.Model(&model.Post{})
	if publishedOnly {
		q = q.Where("status = ?", "published")
	} else if s := c.Query("status"); s == "draft" || s == "published" {
		q = q.Where("status = ?", s)
	}
	if v := c.Query("category_id"); v != "" {
		if id, err := strconv.ParseUint(v, 10, 64); err == nil {
			// a category includes all of its descendants
			q = q.Where("category_id IN ?", loadCats().descendants(uint(id)))
		}
	}
	if v := c.Query("tag_id"); v != "" {
		q = q.Where("id IN (?)", config.DB.Table("post_tags").Select("post_id").Where("tag_id = ?", v))
	}
	if kw := strings.TrimSpace(c.Query("keyword")); kw != "" {
		like := "%" + kw + "%"
		q = q.Where("title LIKE ? OR summary LIKE ? OR content LIKE ?", like, like, like)
	}
	var total int64
	if err := q.Count(&total).Error; err != nil {
		Fail(c, http.StatusInternalServerError, 500, "查询失败")
		return
	}
	posts := []model.Post{}
	err := q.Omit("content").Preload("Tags").
		Order(postOrder).Offset((page - 1) * size).Limit(size).Find(&posts).Error
	if err != nil {
		Fail(c, http.StatusInternalServerError, 500, "查询失败")
		return
	}
	fillCategories(posts)
	OK(c, gin.H{"list": posts, "total": total, "page": page, "page_size": size})
}

// ListPublishedPosts GET /api/posts
func ListPublishedPosts(c *gin.Context) { listPosts(c, true) }

// AdminListPosts GET /api/admin/posts
func AdminListPosts(c *gin.Context) { listPosts(c, false) }

type postLink struct {
	Title string `json:"title"`
	Slug  string `json:"slug"`
}

type relatedPost struct {
	ID          uint      `json:"id"`
	Title       string    `json:"title"`
	Slug        string    `json:"slug"`
	Summary     string    `json:"summary"`
	ReadingTime int       `json:"reading_time"`
	CreatedAt   time.Time `json:"created_at"`
}

type postDetail struct {
	model.Post
	Prev    *postLink     `json:"prev"`
	Next    *postLink     `json:"next"`
	Related []relatedPost `json:"related"`
}

// relatedPosts scores published posts by shared tags (+2 each) and same category (+1).
func relatedPosts(current model.Post, limit int) []relatedPost {
	var candidates []model.Post
	config.DB.Model(&model.Post{}).
		Select("id, title, slug, summary, word_count, category_id, created_at").
		Preload("Tags").
		Where("status = ? AND id <> ?", "published", current.ID).
		Find(&candidates)

	tagSet := make(map[uint]bool, len(current.Tags))
	for _, t := range current.Tags {
		tagSet[t.ID] = true
	}

	type scored struct {
		post  model.Post
		score int
	}
	list := make([]scored, 0, len(candidates))
	for _, cand := range candidates {
		score := 0
		if current.CategoryID != nil && cand.CategoryID != nil && *current.CategoryID == *cand.CategoryID {
			score += 1
		}
		for _, t := range cand.Tags {
			if tagSet[t.ID] {
				score += 2
			}
		}
		list = append(list, scored{post: cand, score: score})
	}

	sort.Slice(list, func(i, j int) bool {
		if list[i].score != list[j].score {
			return list[i].score > list[j].score
		}
		return list[i].post.CreatedAt.After(list[j].post.CreatedAt)
	})

	if len(list) > limit {
		list = list[:limit]
	}
	out := make([]relatedPost, 0, len(list))
	for _, item := range list {
		out = append(out, relatedPost{
			ID:          item.post.ID,
			Title:       item.post.Title,
			Slug:        item.post.Slug,
			Summary:     item.post.Summary,
			ReadingTime: model.ReadingMinutes(item.post.WordCount),
			CreatedAt:   item.post.CreatedAt,
		})
	}
	return out
}

// neighbours finds the previous (older) and next (newer) published posts under postOrder.
func neighbours(id uint) (prev, next *postLink) {
	var rows []struct {
		ID    uint
		Title string
		Slug  string
	}
	config.DB.Model(&model.Post{}).Select("id, title, slug").
		Where("status = ?", "published").Order(postOrder).Find(&rows)
	for i, r := range rows {
		if r.ID != id {
			continue
		}
		if i+1 < len(rows) {
			prev = &postLink{rows[i+1].Title, rows[i+1].Slug}
		}
		if i > 0 {
			next = &postLink{rows[i-1].Title, rows[i-1].Slug}
		}
		return
	}
	return
}

// GetPostBySlug GET /api/posts/:slug (view_count +1)
func GetPostBySlug(c *gin.Context) {
	var p model.Post
	err := config.DB.Preload("Tags").
		Where("slug = ? AND status = ?", c.Param("slug"), "published").First(&p).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		Fail(c, http.StatusNotFound, 404, "文章不存在")
		return
	} else if err != nil {
		Fail(c, http.StatusInternalServerError, 500, "查询失败")
		return
	}
	config.DB.Model(&model.Post{}).Where("id = ?", p.ID).
		UpdateColumn("view_count", gorm.Expr("view_count + 1"))
	p.ViewCount++
	ps := []model.Post{p}
	fillCategories(ps)
	d := postDetail{Post: ps[0]}
	d.Prev, d.Next = neighbours(p.ID)
	d.Related = relatedPosts(p, 4)
	OK(c, d)
}

// AdminGetPost GET /api/admin/posts/:id
func AdminGetPost(c *gin.Context) {
	var p model.Post
	err := config.DB.Preload("Tags").First(&p, c.Param("id")).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		Fail(c, http.StatusNotFound, 404, "文章不存在")
		return
	} else if err != nil {
		Fail(c, http.StatusInternalServerError, 500, "查询失败")
		return
	}
	ps := []model.Post{p}
	fillCategories(ps)
	OK(c, ps[0])
}

type postPayload struct {
	Title      *string    `json:"title"`
	Slug       *string    `json:"slug"`
	Summary    *string    `json:"summary"`
	Content    *string    `json:"content"`
	CoverImage *string    `json:"cover_image"`
	Status     *string    `json:"status"`
	CategoryID *uint      `json:"category_id"`
	TagIDs     *[]uint    `json:"tag_ids"`
	IsTop      *bool      `json:"is_top"`
	CreatedAt  *time.Time `json:"created_at"`
}

func validStatus(s string) bool { return s == "draft" || s == "published" }

func isDup(err error) bool {
	return err != nil && strings.Contains(err.Error(), "Duplicate entry")
}

// AdminCreatePost POST /api/admin/posts
func AdminCreatePost(c *gin.Context) {
	var req postPayload
	if err := c.ShouldBindJSON(&req); err != nil ||
		req.Title == nil || *req.Title == "" || req.Slug == nil || *req.Slug == "" ||
		req.Content == nil || req.Status == nil || !validStatus(*req.Status) {
		Fail(c, http.StatusBadRequest, 400, "title、slug、content、status 为必填，status 须为 draft/published")
		return
	}
	p := model.Post{
		Title: *req.Title, Slug: *req.Slug, Content: *req.Content, Status: *req.Status,
		CategoryID: req.CategoryID, WordCount: model.CountWords(*req.Content),
	}
	if req.Summary != nil {
		p.Summary = *req.Summary
	}
	if req.CoverImage != nil {
		p.CoverImage = *req.CoverImage
	}
	if req.IsTop != nil {
		p.IsTop = *req.IsTop
	}
	if req.CreatedAt != nil {
		p.CreatedAt = *req.CreatedAt
	}
	if req.TagIDs != nil && len(*req.TagIDs) > 0 {
		config.DB.Find(&p.Tags, *req.TagIDs)
	}
	if err := config.DB.Create(&p).Error; err != nil {
		if isDup(err) {
			Fail(c, http.StatusBadRequest, 400, "slug 已存在")
		} else {
			Fail(c, http.StatusBadRequest, 400, "创建失败：请检查 category_id 是否有效")
		}
		return
	}
	OK(c, gin.H{"id": p.ID, "title": p.Title, "slug": p.Slug, "status": p.Status, "created_at": p.CreatedAt})
}

// AdminUpdatePost PUT /api/admin/posts/:id (partial update)
func AdminUpdatePost(c *gin.Context) {
	var p model.Post
	if err := config.DB.First(&p, c.Param("id")).Error; err != nil {
		Fail(c, http.StatusNotFound, 404, "文章不存在")
		return
	}
	var req postPayload
	if err := c.ShouldBindJSON(&req); err != nil {
		Fail(c, http.StatusBadRequest, 400, "请求参数错误")
		return
	}
	if req.Status != nil && !validStatus(*req.Status) {
		Fail(c, http.StatusBadRequest, 400, "status 须为 draft/published")
		return
	}
	upd := map[string]interface{}{}
	set := func(k string, v *string) {
		if v != nil {
			upd[k] = *v
		}
	}
	set("title", req.Title)
	set("slug", req.Slug)
	set("summary", req.Summary)
	set("content", req.Content)
	set("cover_image", req.CoverImage)
	set("status", req.Status)
	if req.Content != nil {
		upd["word_count"] = model.CountWords(*req.Content)
	}
	if req.CategoryID != nil {
		upd["category_id"] = *req.CategoryID
	}
	if req.IsTop != nil {
		upd["is_top"] = *req.IsTop
	}
	if req.CreatedAt != nil {
		upd["created_at"] = *req.CreatedAt
	}
	err := config.DB.Transaction(func(tx *gorm.DB) error {
		if len(upd) > 0 {
			if err := tx.Model(&p).Updates(upd).Error; err != nil {
				return err
			}
		}
		if req.TagIDs != nil {
			var tags []model.Tag
			if len(*req.TagIDs) > 0 {
				tx.Find(&tags, *req.TagIDs)
			}
			return tx.Model(&p).Association("Tags").Replace(tags)
		}
		return nil
	})
	if err != nil {
		if isDup(err) {
			Fail(c, http.StatusBadRequest, 400, "slug 已存在")
		} else {
			Fail(c, http.StatusBadRequest, 400, "更新失败")
		}
		return
	}
	OK(c, nil)
}

// AdminDeletePost DELETE /api/admin/posts/:id
func AdminDeletePost(c *gin.Context) {
	var p model.Post
	if err := config.DB.First(&p, c.Param("id")).Error; err != nil {
		Fail(c, http.StatusNotFound, 404, "文章不存在")
		return
	}
	err := config.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&p).Association("Tags").Clear(); err != nil {
			return err
		}
		return tx.Delete(&p).Error
	})
	if err != nil {
		Fail(c, http.StatusInternalServerError, 500, "删除失败")
		return
	}
	OK(c, nil)
}

// AdminUpdatePostStatus PATCH /api/admin/posts/:id/status
func AdminUpdatePostStatus(c *gin.Context) {
	var req struct {
		Status string `json:"status"`
	}
	if err := c.ShouldBindJSON(&req); err != nil || !validStatus(req.Status) {
		Fail(c, http.StatusBadRequest, 400, "status 须为 draft/published")
		return
	}
	var n int64
	config.DB.Model(&model.Post{}).Where("id = ?", c.Param("id")).Count(&n)
	if n == 0 {
		Fail(c, http.StatusNotFound, 404, "文章不存在")
		return
	}
	if err := config.DB.Model(&model.Post{}).Where("id = ?", c.Param("id")).Update("status", req.Status).Error; err != nil {
		Fail(c, http.StatusInternalServerError, 500, "更新失败")
		return
	}
	OK(c, nil)
}

// Archives GET /api/archives — published posts grouped by year and month (strictly chronological).
func Archives(c *gin.Context) {
	var rows []struct {
		ID        uint      `json:"id"`
		Title     string    `json:"title"`
		Slug      string    `json:"slug"`
		CreatedAt time.Time `json:"created_at"`
	}
	config.DB.Model(&model.Post{}).Select("id, title, slug, created_at").
		Where("status = ?", "published").Order("created_at DESC, id DESC").Find(&rows)

	type item struct {
		ID        uint      `json:"id"`
		Title     string    `json:"title"`
		Slug      string    `json:"slug"`
		CreatedAt time.Time `json:"created_at"`
	}
	byYM := map[[2]int][]item{}
	years, months := map[int]bool{}, map[int]map[int]bool{}
	for _, r := range rows {
		y, m := r.CreatedAt.Year(), int(r.CreatedAt.Month())
		byYM[[2]int{y, m}] = append(byYM[[2]int{y, m}], item(r))
		years[y] = true
		if months[y] == nil {
			months[y] = map[int]bool{}
		}
		months[y][m] = true
	}
	out := []gin.H{}
	for _, y := range sortedKeysDesc(years) {
		ms, count := []gin.H{}, 0
		for _, m := range sortedKeysDesc(months[y]) {
			posts := byYM[[2]int{y, m}]
			count += len(posts)
			ms = append(ms, gin.H{"month": m, "posts": posts})
		}
		out = append(out, gin.H{"year": y, "count": count, "months": ms})
	}
	OK(c, out)
}
