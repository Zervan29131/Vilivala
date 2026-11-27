package router

import (
	"time"

	"github.com/gin-gonic/gin"
	"vilivala-backend/config"
	"vilivala-backend/internal/handler"
	"vilivala-backend/internal/middleware"
)

// Setup registers all routes and returns the configured Engine
func Setup() *gin.Engine {
	r := gin.Default()

	// Global middleware
	r.Use(middleware.CORS())

	r.Static("/uploads", handler.ResolveUploadDir())

	// SEO & feed endpoints
	r.GET("/robots.txt", handler.RobotsTxt)
	r.GET("/sitemap.xml", handler.SitemapXML)
	r.GET("/rss.xml", handler.RSSFeed)

	loginLimiter := middleware.NewRateLimiter(15, 5*time.Minute)

	api := r.Group("/api")
	{
		// --- Public routes ---
		api.GET("/health", handler.Health)
		api.POST("/auth/login", loginLimiter.Middleware("登录尝试过于频繁，请稍后再试"), handler.Login)

		api.GET("/posts", handler.ListPublishedPosts)
		api.GET("/posts/:slug", handler.GetPostBySlug)
		api.GET("/search", handler.Search)
		api.GET("/categories", handler.ListCategories)
		api.GET("/tags", handler.ListTags)
		api.GET("/archives", handler.Archives)
		api.GET("/pages", handler.ListPages)
		api.GET("/pages/:slug", handler.GetPage)
		api.GET("/site-config", handler.GetSiteConfig)

		// --- Admin routes (JWT protected) ---
		admin := api.Group("/admin")
		admin.Use(middleware.JWTAuth())
		{
			// Dashboard & Site Config
			admin.GET("/dashboard", handler.AdminDashboard)
			admin.PUT("/site-config", handler.AdminPutSiteConfig)

			// Posts
			admin.GET("/posts", handler.AdminListPosts)
			admin.GET("/posts/:id", handler.AdminGetPost)
			admin.POST("/posts", handler.AdminCreatePost)
			admin.PUT("/posts/:id", handler.AdminUpdatePost)
			admin.DELETE("/posts/:id", handler.AdminDeletePost)
			admin.PATCH("/posts/:id/status", handler.AdminUpdatePostStatus)

			// Categories
			admin.POST("/categories", handler.AdminCreateCategory)
			admin.PUT("/categories/:id", handler.AdminUpdateCategory)
			admin.DELETE("/categories/:id", handler.AdminDeleteCategory)

			// Pages & upload
			admin.GET("/pages", handler.ListPages)
			admin.PUT("/pages/:slug", handler.AdminPutPage)
			admin.DELETE("/pages/:slug", handler.AdminDeletePage)
			admin.POST("/upload", handler.AdminUpload)

			// Tags
			admin.POST("/tags", handler.AdminCreateTag)
			admin.PUT("/tags/:id", handler.AdminUpdateTag)
			admin.DELETE("/tags/:id", handler.AdminDeleteTag)
		}
	}

	// Optional lightweight SPA + server-side meta injection when FRONTEND_DIST is set
	if config.AppConfig.FrontendDist != "" {
		r.NoRoute(handler.ServeSPAWithMeta(config.AppConfig.FrontendDist))
	}

	return r
}
