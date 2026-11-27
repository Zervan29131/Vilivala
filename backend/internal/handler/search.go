package handler

import (
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"vilivala-backend/config"
	"vilivala-backend/internal/model"
)

type searchItem struct {
	ID           uint      `json:"id"`
	Title        string    `json:"title"`
	Slug         string    `json:"slug"`
	Snippet      string    `json:"snippet"`
	CategoryName string    `json:"category_name"`
	CreatedAt    time.Time `json:"created_at"`
}

var (
	reCodeBlock = regexp.MustCompile("(?s)```.*?```")
	reMdImage   = regexp.MustCompile(`!\[[^\]]*\]\([^)]*\)`)
	reMdLink    = regexp.MustCompile(`\[([^\]]+)\]\([^)]*\)`)
	reHTMLTag   = regexp.MustCompile(`<[^>]+>`)
	reMdMarks   = regexp.MustCompile("[#*_>`~|\\[\\]]+")
	reSpaces    = regexp.MustCompile(`\s+`)
)

// ExtractSnippet strips markdown noise and returns a plain-text excerpt centered on query.
func ExtractSnippet(content, summary, query string, maxRunes int) string {
	if maxRunes <= 0 {
		maxRunes = 120
	}
	plain := cleanMarkdown(content)
	if plain == "" {
		plain = cleanMarkdown(summary)
	}
	if plain == "" {
		return ""
	}

	runes := []rune(plain)
	lowerPlain := strings.ToLower(plain)
	lowerQ := strings.ToLower(strings.TrimSpace(query))

	idxByte := -1
	if lowerQ != "" {
		idxByte = strings.Index(lowerPlain, lowerQ)
	}
	if idxByte < 0 {
		if len(runes) <= maxRunes {
			return plain
		}
		return string(runes[:maxRunes]) + "…"
	}

	hitRune := utf8.RuneCountInString(lowerPlain[:idxByte])
	start := hitRune - maxRunes/3
	if start < 0 {
		start = 0
	}
	end := start + maxRunes
	if end > len(runes) {
		end = len(runes)
		if end-maxRunes > 0 {
			start = end - maxRunes
		}
	}

	var b strings.Builder
	if start > 0 {
		b.WriteString("…")
	}
	b.WriteString(string(runes[start:end]))
	if end < len(runes) {
		b.WriteString("…")
	}
	return b.String()
}

func cleanMarkdown(s string) string {
	s = reCodeBlock.ReplaceAllString(s, " ")
	s = reMdImage.ReplaceAllString(s, " ")
	s = reMdLink.ReplaceAllString(s, "$1")
	s = reHTMLTag.ReplaceAllString(s, " ")
	s = reMdMarks.ReplaceAllString(s, " ")
	s = reSpaces.ReplaceAllString(s, " ")
	return strings.TrimSpace(s)
}

// Search GET /api/search?q=keyword&limit=12
func Search(c *gin.Context) {
	q := strings.TrimSpace(c.Query("q"))
	if q == "" {
		OK(c, gin.H{"list": []searchItem{}, "total": 0})
		return
	}
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "12"))
	if limit < 1 {
		limit = 12
	}
	if limit > 30 {
		limit = 30
	}

	var posts []model.Post
	var total int64

	// Prefer MySQL FULLTEXT (ngram) when query has >= 2 runes; gracefully fall back to LIKE.
	useLike := utf8.RuneCountInString(q) < 2
	if !useLike {
		ft := config.DB.Model(&model.Post{}).
			Where("status = ? AND MATCH(title, summary, content) AGAINST (? IN NATURAL LANGUAGE MODE)", "published", q)
		if err := ft.Count(&total).Error; err != nil || total == 0 {
			useLike = true
		} else {
			err = ft.Select("id, title, slug, summary, content, category_id, created_at").
				Order(gorm.Expr("is_top DESC, MATCH(title, summary, content) AGAINST (? IN NATURAL LANGUAGE MODE) DESC, created_at DESC", q)).
				Limit(limit).Find(&posts).Error
			if err != nil {
				useLike = true
				posts = nil
			}
		}
	}

	if useLike {
		like := "%" + q + "%"
		base := config.DB.Model(&model.Post{}).
			Where("status = ? AND (title LIKE ? OR summary LIKE ? OR content LIKE ?)", "published", like, like, like)
		base.Count(&total)
		base.Select("id, title, slug, summary, content, category_id, created_at").
			Order(gorm.Expr("CASE WHEN title LIKE ? THEN 0 ELSE 1 END, is_top DESC, created_at DESC", like)).
			Limit(limit).Find(&posts)
	}

	cats := loadCats()
	items := make([]searchItem, 0, len(posts))
	for _, p := range posts {
		catName := ""
		if p.CategoryID != nil {
			if cat, ok := cats.byID[*p.CategoryID]; ok {
				catName = cat.Name
			}
		}
		items = append(items, searchItem{
			ID:           p.ID,
			Title:        p.Title,
			Slug:         p.Slug,
			Snippet:      ExtractSnippet(p.Content, p.Summary, q, 120),
			CategoryName: catName,
			CreatedAt:    p.CreatedAt,
		})
	}
	OK(c, gin.H{"list": items, "total": total})
}
