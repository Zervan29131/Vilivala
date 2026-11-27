package router

import (
	"github.com/gin-gonic/gin"
	"vilivala-backend/internal/handler"
	"vilivala-backend/internal/middleware"
)

// Setup registers all routes and returns the configured Engine
func Setup() *gin.Engine {
	r := gin.Default()

	// Global middleware
	r.Use(middleware.CORS())

	r.Static("/uploads", "./"+handler.UploadDir)

	api := r.Group("/api")
	{
		// --- Public routes ---
		api.POST("/auth/login", handler.Login)

		api.GET("/posts", handler.ListPublishedPosts)
		api.GET("/posts/:slug", handler.GetPostBySlug)
		api.GET("/categories", handler.ListCategories)
		api.GET("/tags", handler.ListTags)
		api.GET("/archives", handler.Archives)
		api.GET("/pages", handler.ListPages)
		api.GET("/pages/:slug", handler.GetPage)

		// --- Admin routes (JWT protected) ---
		admin := api.Group("/admin")
		admin.Use(middleware.JWTAuth())
		{
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

	return r
}
