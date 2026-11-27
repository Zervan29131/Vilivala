package handler

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"sort"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/gin-gonic/gin"
	"vilivala-backend/config"
	"vilivala-backend/internal/model"
)

// SocialLink represents a blogger social link in dynamic site config.
type SocialLink struct {
	Title string `json:"title"`
	Link  string `json:"link"`
}

// SiteStats holds live aggregate counters returned alongside SiteConfig.
type SiteStats struct {
	PublishedPosts int64 `json:"published_posts"`
	TotalWords     int64 `json:"total_words"`
	TotalViews     int64 `json:"total_views"`
}

// SiteConfig holds customizable site-wide branding, widgets, and filing settings.
type SiteConfig struct {
	SiteTitle       string       `json:"site_title"`
	SiteDescription string       `json:"site_description"`
	HeroTitle       string       `json:"hero_title"`
	HeroTagline     string       `json:"hero_tagline"`
	Subtitles       []string     `json:"subtitles"`
	BloggerName     string       `json:"blogger_name"`
	BloggerSlogan   string       `json:"blogger_slogan"`
	BloggerAvatar   string       `json:"blogger_avatar"`
	NoticeEnabled   bool         `json:"notice_enabled"`
	NoticeMode      string       `json:"notice_mode"` // custom | hitokoto
	NoticeText      string       `json:"notice_text"`
	RewardEnabled   bool         `json:"reward_enabled"`
	RewardWording   string       `json:"reward_wording"`
	CopyrightEnable bool         `json:"copyright_enabled"`
	CopyrightNotice string       `json:"copyright_notice"`
	SinceYear       int          `json:"since_year"`
	ICPNumber       string       `json:"icp_number"`
	ICPURL          string       `json:"icp_url"`
	GonganNumber    string       `json:"gongan_number"`
	GonganURL       string       `json:"gongan_url"`
	Socials         []SocialLink `json:"socials"`
	Stats           *SiteStats   `json:"stats,omitempty"`
}

// DefaultSiteConfig returns the out-of-the-box configuration matching Zervan's Vilivala branding.
func DefaultSiteConfig() SiteConfig {
	return SiteConfig{
		SiteTitle:       "Zervan的小站",
		SiteDescription: "Zervanのblog",
		HeroTitle:       "这是存在于某处的幻想世界",
		HeroTagline:     "Be the change you want to see in the World",
		Subtitles: []string{
			"Be the change you want to see in the World",
			"记录代码、阅读与思考的角落",
			"Stay hungry, stay foolish",
		},
		BloggerName:     "Zervan",
		BloggerSlogan:   "天下最普通的人之一",
		BloggerAvatar:   "/img/console.jpg",
		NoticeEnabled:   true,
		NoticeMode:      "custom",
		NoticeText:      "Stay hungry, stay foolish · 欢迎来到 Zervan 的小站",
		RewardEnabled:   true,
		RewardWording:   "觉得文章有帮助？请我喝杯咖啡吧 ~",
		CopyrightEnable: true,
		CopyrightNotice: "本博客所有文章除特别声明外，著作权归作者所有。转载请注明出处！",
		SinceYear:       2021,
		ICPNumber:       "",
		ICPURL:          "https://beian.miit.gov.cn/",
		GonganNumber:    "",
		GonganURL:       "https://www.beian.gov.cn/",
		Socials: []SocialLink{
			{Title: "微博", Link: "https://weibo.com/zervan"},
			{Title: "bilibili", Link: "https://space.bilibili.com/1650366974"},
			{Title: "知乎", Link: "https://www.zhihu.com/people/zervan"},
			{Title: "抖音", Link: "https://v.douyin.com/L7jT1xp/"},
			{Title: "发邮件", Link: "mailto:nickmechlin01@outlook.com"},
		},
	}
}

func isSafeURL(raw string) bool {
	s := strings.TrimSpace(raw)
	if s == "" {
		return true
	}
	lower := strings.ToLower(s)
	return strings.HasPrefix(lower, "http://") ||
		strings.HasPrefix(lower, "https://") ||
		strings.HasPrefix(lower, "mailto:") ||
		(strings.HasPrefix(s, "/") && !strings.HasPrefix(s, "//"))
}

// ValidateSiteConfig normalizes and validates incoming SiteConfig fields.
func ValidateSiteConfig(cfg *SiteConfig) error {
	cfg.SiteTitle = strings.TrimSpace(cfg.SiteTitle)
	if n := utf8.RuneCountInString(cfg.SiteTitle); n < 1 || n > 64 {
		return errors.New("站点名称需为 1–64 个字符")
	}
	if cfg.NoticeMode == "" {
		cfg.NoticeMode = "custom"
	}
	if cfg.NoticeMode != "custom" && cfg.NoticeMode != "hitokoto" {
		return errors.New("notice_mode 须为 custom 或 hitokoto")
	}
	if cfg.SinceYear < 1990 || cfg.SinceYear > 2100 {
		cfg.SinceYear = 2021
	}
	if !isSafeURL(cfg.BloggerAvatar) || !isSafeURL(cfg.ICPURL) || !isSafeURL(cfg.GonganURL) {
		return errors.New("链接地址仅支持 http(s)、mailto 或站内路径")
	}
	cleanedSubs := make([]string, 0, len(cfg.Subtitles))
	for _, s := range cfg.Subtitles {
		if t := strings.TrimSpace(s); t != "" {
			cleanedSubs = append(cleanedSubs, t)
		}
	}
	if len(cleanedSubs) == 0 && strings.TrimSpace(cfg.HeroTagline) != "" {
		cleanedSubs = []string{strings.TrimSpace(cfg.HeroTagline)}
	}
	cfg.Subtitles = cleanedSubs

	cleanedSocials := make([]SocialLink, 0, len(cfg.Socials))
	for _, item := range cfg.Socials {
		t := strings.TrimSpace(item.Title)
		l := strings.TrimSpace(item.Link)
		if t == "" && l == "" {
			continue
		}
		if t == "" || l == "" || !isSafeURL(l) {
			return errors.New("社交链接需填写名称且使用合法 URL 协议")
		}
		cleanedSocials = append(cleanedSocials, SocialLink{Title: t, Link: l})
	}
	cfg.Socials = cleanedSocials
	cfg.Stats = nil
	return nil
}

func loadStoredSiteConfig() SiteConfig {
	cfg := DefaultSiteConfig()
	var row model.SiteSetting
	if err := config.DB.First(&row, 1).Error; err == nil && row.Payload != "" {
		_ = json.Unmarshal([]byte(row.Payload), &cfg)
	}
	return cfg
}

func querySiteStats() *SiteStats {
	var s SiteStats
	config.DB.Raw(`
		SELECT
			COUNT(*) AS published_posts,
			COALESCE(SUM(word_count), 0) AS total_words,
			COALESCE(SUM(view_count), 0) AS total_views
		FROM posts
		WHERE status = 'published'
	`).Scan(&s)
	return &s
}

// GetSiteConfig GET /api/site-config
func GetSiteConfig(c *gin.Context) {
	cfg := loadStoredSiteConfig()
	cfg.Stats = querySiteStats()
	OK(c, cfg)
}

// AdminPutSiteConfig PUT /api/admin/site-config
func AdminPutSiteConfig(c *gin.Context) {
	var cfg SiteConfig
	if err := c.ShouldBindJSON(&cfg); err != nil {
		Fail(c, http.StatusBadRequest, 400, "请求参数格式错误")
		return
	}
	if err := ValidateSiteConfig(&cfg); err != nil {
		Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	raw, err := json.Marshal(cfg)
	if err != nil {
		Fail(c, http.StatusInternalServerError, 500, "序列化失败")
		return
	}
	row := model.SiteSetting{ID: 1, Payload: string(raw), UpdatedAt: time.Now()}
	if err := config.DB.Save(&row).Error; err != nil {
		Fail(c, http.StatusInternalServerError, 500, "保存配置失败")
		return
	}
	log.Printf("[Admin] site-config updated by user=%v", c.GetString("username"))
	cfg.Stats = querySiteStats()
	OK(c, cfg)
}

// AdminDashboard GET /api/admin/dashboard
func AdminDashboard(c *gin.Context) {
	var counts struct {
		TotalPosts     int64 `json:"total_posts"`
		PublishedPosts int64 `json:"published_posts"`
		DraftPosts     int64 `json:"draft_posts"`
		Categories     int64 `json:"categories"`
		Tags           int64 `json:"tags"`
		Pages          int64 `json:"pages"`
		TotalWords     int64 `json:"total_words"`
		TotalViews     int64 `json:"total_views"`
	}
	config.DB.Raw(`
		SELECT
			(SELECT COUNT(*) FROM posts) AS total_posts,
			(SELECT COUNT(*) FROM posts WHERE status = 'published') AS published_posts,
			(SELECT COUNT(*) FROM posts WHERE status = 'draft') AS draft_posts,
			(SELECT COUNT(*) FROM categories) AS categories,
			(SELECT COUNT(*) FROM tags) AS tags,
			(SELECT COUNT(*) FROM pages) AS pages,
			(SELECT COALESCE(SUM(word_count), 0) FROM posts) AS total_words,
			(SELECT COALESCE(SUM(view_count), 0) FROM posts) AS total_views
	`).Scan(&counts)

	type topPost struct {
		ID        uint      `json:"id"`
		Title     string    `json:"title"`
		Slug      string    `json:"slug"`
		Status    string    `json:"status"`
		ViewCount uint      `json:"view_count"`
		WordCount int       `json:"word_count"`
		CreatedAt time.Time `json:"created_at"`
	}
	topPosts := []topPost{}
	config.DB.Model(&model.Post{}).
		Select("id, title, slug, status, view_count, word_count, created_at").
		Order("view_count DESC, id DESC").
		Limit(10).Find(&topPosts)

	type recentPost struct {
		ID        uint      `json:"id"`
		Title     string    `json:"title"`
		Slug      string    `json:"slug"`
		Status    string    `json:"status"`
		UpdatedAt time.Time `json:"updated_at"`
	}
	recentPosts := []recentPost{}
	config.DB.Model(&model.Post{}).
		Select("id, title, slug, status, updated_at").
		Order("updated_at DESC, id DESC").
		Limit(6).Find(&recentPosts)

	type catDist struct {
		ID         uint   `json:"id"`
		Name       string `json:"name"`
		Slug       string `json:"slug"`
		TotalCount int64  `json:"total_count"`
	}
	idx := loadCats()
	var catCounts []struct {
		CategoryID uint  `gorm:"column:category_id"`
		Cnt        int64 `gorm:"column:cnt"`
	}
	config.DB.Model(&model.Post{}).
		Select("category_id, COUNT(*) AS cnt").
		Where("status = ? AND category_id IS NOT NULL", "published").
		Group("category_id").Scan(&catCounts)
	directCnt := map[uint]int64{}
	for _, r := range catCounts {
		directCnt[r.CategoryID] = r.Cnt
	}

	dist := []catDist{}
	for _, rootID := range idx.roots {
		cat := idx.byID[rootID]
		if cat == nil {
			continue
		}
		var sum int64
		for _, did := range idx.descendants(rootID) {
			sum += directCnt[did]
		}
		dist = append(dist, catDist{
			ID:         cat.ID,
			Name:       cat.Name,
			Slug:       cat.Slug,
			TotalCount: sum,
		})
	}
	sort.Slice(dist, func(i, j int) bool { return dist[i].TotalCount > dist[j].TotalCount })

	OK(c, gin.H{
		"total_posts":           counts.TotalPosts,
		"published_posts":       counts.PublishedPosts,
		"draft_posts":           counts.DraftPosts,
		"categories":            counts.Categories,
		"tags":                  counts.Tags,
		"pages":                 counts.Pages,
		"total_words":           counts.TotalWords,
		"total_views":           counts.TotalViews,
		"top_posts":             topPosts,
		"recent_posts":          recentPosts,
		"category_distribution": dist,
	})
}
