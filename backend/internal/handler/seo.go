package handler

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"html"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"vilivala-backend/config"
	"vilivala-backend/internal/model"
)

func siteBaseURL(c *gin.Context) string {
	if u := strings.TrimRight(strings.TrimSpace(config.AppConfig.SiteURL), "/"); u != "" {
		return u
	}
	scheme := "http"
	if c.Request.TLS != nil || c.GetHeader("X-Forwarded-Proto") == "https" {
		scheme = "https"
	}
	return scheme + "://" + c.Request.Host
}

// Health GET /api/health
func Health(c *gin.Context) {
	if config.DB != nil {
		if sqlDB, err := config.DB.DB(); err == nil {
			if err := sqlDB.Ping(); err != nil {
				Fail(c, http.StatusServiceUnavailable, 503, "database unavailable")
				return
			}
		}
	}
	OK(c, gin.H{"status": "ok"})
}

// RobotsTxt GET /robots.txt
func RobotsTxt(c *gin.Context) {
	base := siteBaseURL(c)
	body := fmt.Sprintf("User-agent: *\nAllow: /\nDisallow: /admin\n\nSitemap: %s/sitemap.xml\n", base)
	c.Data(http.StatusOK, "text/plain; charset=utf-8", []byte(body))
}

type xmlURLSet struct {
	XMLName xml.Name `xml:"urlset"`
	Xmlns   string   `xml:"xmlns,attr"`
	URLs    []xmlURL `xml:"url"`
}

type xmlURL struct {
	Loc        string `xml:"loc"`
	LastMod    string `xml:"lastmod,omitempty"`
	ChangeFreq string `xml:"changefreq,omitempty"`
	Priority   string `xml:"priority,omitempty"`
}

// SitemapXML GET /sitemap.xml
func SitemapXML(c *gin.Context) {
	base := siteBaseURL(c)
	urls := []xmlURL{
		{Loc: base + "/", ChangeFreq: "daily", Priority: "1.0"},
		{Loc: base + "/archives", ChangeFreq: "weekly", Priority: "0.8"},
		{Loc: base + "/categories", ChangeFreq: "weekly", Priority: "0.8"},
		{Loc: base + "/terminal", ChangeFreq: "monthly", Priority: "0.5"},
	}

	var pages []model.Page
	config.DB.Select("slug, updated_at").Find(&pages)
	for _, p := range pages {
		urls = append(urls, xmlURL{
			Loc:        base + "/pages/" + p.Slug,
			LastMod:    p.UpdatedAt.Format("2006-01-02"),
			ChangeFreq: "monthly",
			Priority:   "0.6",
		})
	}

	var posts []model.Post
	config.DB.Select("slug, updated_at, created_at").
		Where("status = ?", "published").
		Order("created_at DESC").Find(&posts)
	for _, p := range posts {
		mod := p.UpdatedAt
		if mod.IsZero() {
			mod = p.CreatedAt
		}
		urls = append(urls, xmlURL{
			Loc:        base + "/posts/" + p.Slug,
			LastMod:    mod.Format("2006-01-02"),
			ChangeFreq: "monthly",
			Priority:   "0.7",
		})
	}

	out, err := xml.MarshalIndent(xmlURLSet{
		Xmlns: "http://www.sitemaps.org/schemas/sitemap/0.9",
		URLs:  urls,
	}, "", "  ")
	if err != nil {
		Fail(c, http.StatusInternalServerError, 500, "生成 sitemap 失败")
		return
	}
	c.Data(http.StatusOK, "application/xml; charset=utf-8", append([]byte(xml.Header), out...))
}

type rssDocument struct {
	XMLName xml.Name   `xml:"rss"`
	Version string     `xml:"version,attr"`
	Channel rssChannel `xml:"channel"`
}

type rssChannel struct {
	Title         string    `xml:"title"`
	Link          string    `xml:"link"`
	Description   string    `xml:"description"`
	Language      string    `xml:"language"`
	LastBuildDate string    `xml:"lastBuildDate"`
	Items         []rssItem `xml:"item"`
}

type rssItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	GUID        string `xml:"guid"`
	PubDate     string `xml:"pubDate"`
	Description string `xml:"description"`
}

// RSSFeed GET /rss.xml
func RSSFeed(c *gin.Context) {
	base := siteBaseURL(c)
	var posts []model.Post
	config.DB.Select("title, slug, summary, content, created_at").
		Where("status = ?", "published").
		Order("created_at DESC, id DESC").
		Limit(30).Find(&posts)

	items := make([]rssItem, 0, len(posts))
	for _, p := range posts {
		desc := strings.TrimSpace(p.Summary)
		if desc == "" {
			desc = ExtractSnippet(p.Content, "", "", 180)
		}
		link := base + "/posts/" + p.Slug
		items = append(items, rssItem{
			Title:       p.Title,
			Link:        link,
			GUID:        link,
			PubDate:     p.CreatedAt.Format(time.RFC1123Z),
			Description: desc,
		})
	}

	lastBuild := time.Now().Format(time.RFC1123Z)
	if len(posts) > 0 {
		lastBuild = posts[0].CreatedAt.Format(time.RFC1123Z)
	}

	doc := rssDocument{
		Version: "2.0",
		Channel: rssChannel{
			Title:         config.AppConfig.SiteTitle,
			Link:          base,
			Description:   config.AppConfig.SiteDescription,
			Language:      "zh-CN",
			LastBuildDate: lastBuild,
			Items:         items,
		},
	}
	out, err := xml.MarshalIndent(doc, "", "  ")
	if err != nil {
		Fail(c, http.StatusInternalServerError, 500, "生成 RSS 失败")
		return
	}
	c.Data(http.StatusOK, "application/rss+xml; charset=utf-8", append([]byte(xml.Header), out...))
}

// PageMeta holds lightweight server-side SEO tags injected into index.html.
type PageMeta struct {
	Title       string
	Description string
	URL         string
	Type        string // website | article
	PublishedAt string
	UpdatedAt   string
	Author      string
}

var (
	reTitleTag = regexp.MustCompile(`(?is)<title>.*?</title>`)
	reMetaDesc = regexp.MustCompile(`(?is)<meta\s+name="description"[^>]*>`)
)

// InjectMetaHTML replaces <title> and inserts OG/Twitter/canonical/JSON-LD before </head>.
func InjectMetaHTML(templateHTML string, meta PageMeta) string {
	titleEsc := html.EscapeString(meta.Title)
	descEsc := html.EscapeString(meta.Description)
	urlEsc := html.EscapeString(meta.URL)
	ogType := meta.Type
	if ogType == "" {
		ogType = "website"
	}

	res := reTitleTag.ReplaceAllString(templateHTML, "<title>"+titleEsc+"</title>")
	res = reMetaDesc.ReplaceAllString(res, "")

	var tags strings.Builder
	tags.WriteString("\n    <meta name=\"description\" content=\"" + descEsc + "\" />")
	tags.WriteString("\n    <link rel=\"canonical\" href=\"" + urlEsc + "\" />")
	tags.WriteString("\n    <meta property=\"og:type\" content=\"" + html.EscapeString(ogType) + "\" />")
	tags.WriteString("\n    <meta property=\"og:title\" content=\"" + titleEsc + "\" />")
	tags.WriteString("\n    <meta property=\"og:description\" content=\"" + descEsc + "\" />")
	tags.WriteString("\n    <meta property=\"og:url\" content=\"" + urlEsc + "\" />")
	tags.WriteString("\n    <meta name=\"twitter:card\" content=\"summary\" />")
	tags.WriteString("\n    <meta name=\"twitter:title\" content=\"" + titleEsc + "\" />")
	tags.WriteString("\n    <meta name=\"twitter:description\" content=\"" + descEsc + "\" />")

	if ogType == "article" {
		ld := map[string]any{
			"@context":      "https://schema.org",
			"@type":         "BlogPosting",
			"headline":      meta.Title,
			"description":   meta.Description,
			"url":           meta.URL,
			"datePublished": meta.PublishedAt,
			"dateModified":  meta.UpdatedAt,
			"author": map[string]string{
				"@type": "Person",
				"name":  meta.Author,
			},
		}
		if b, err := json.Marshal(ld); err == nil {
			safeJSON := strings.ReplaceAll(string(b), "</", `<\/`)
			tags.WriteString("\n    <script type=\"application/ld+json\">" + safeJSON + "</script>")
		}
	}

	if idx := strings.Index(strings.ToLower(res), "</head>"); idx >= 0 {
		return res[:idx] + tags.String() + "\n  " + res[idx:]
	}
	return res
}

// ServeSPAWithMeta serves static files from distDir and injects dynamic SEO meta into index.html.
func ServeSPAWithMeta(distDir string) gin.HandlerFunc {
	indexPath := filepath.Join(distDir, "index.html")
	return func(c *gin.Context) {
		reqPath := c.Request.URL.Path
		if strings.HasPrefix(reqPath, "/api/") {
			Fail(c, http.StatusNotFound, 404, "not found")
			return
		}

		// Try serving an existing static asset first.
		clean := filepath.Clean(strings.TrimPrefix(reqPath, "/"))
		if clean != "." && clean != "" && !strings.HasPrefix(clean, "..") {
			candidate := filepath.Join(distDir, clean)
			if info, err := os.Stat(candidate); err == nil && !info.IsDir() {
				c.File(candidate)
				return
			}
		}

		raw, err := os.ReadFile(indexPath)
		if err != nil {
			c.String(http.StatusNotFound, "index.html not found")
			return
		}

		base := siteBaseURL(c)
		meta := PageMeta{
			Title:       config.AppConfig.SiteTitle,
			Description: config.AppConfig.SiteDescription,
			URL:         base + reqPath,
			Type:        "website",
			Author:      "Zervan",
		}

		if strings.HasPrefix(reqPath, "/posts/") {
			slug := strings.TrimPrefix(reqPath, "/posts/")
			if slug != "" && !strings.Contains(slug, "/") {
				var p model.Post
				if err := config.DB.Select("title, slug, summary, content, created_at, updated_at").
					Where("slug = ? AND status = ?", slug, "published").First(&p).Error; err == nil {
					desc := strings.TrimSpace(p.Summary)
					if desc == "" {
						desc = ExtractSnippet(p.Content, "", "", 150)
					}
					meta.Title = p.Title + " - " + config.AppConfig.SiteTitle
					meta.Description = desc
					meta.Type = "article"
					meta.PublishedAt = p.CreatedAt.Format(time.RFC3339)
					meta.UpdatedAt = p.UpdatedAt.Format(time.RFC3339)
				}
			}
		} else if strings.HasPrefix(reqPath, "/pages/") {
			slug := strings.TrimPrefix(reqPath, "/pages/")
			if slug != "" && !strings.Contains(slug, "/") {
				var pg model.Page
				if err := config.DB.Select("title, slug, content").Where("slug = ?", slug).First(&pg).Error; err == nil {
					meta.Title = pg.Title + " - " + config.AppConfig.SiteTitle
					meta.Description = ExtractSnippet(pg.Content, "", "", 150)
				}
			}
		}

		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(InjectMetaHTML(string(raw), meta)))
	}
}
