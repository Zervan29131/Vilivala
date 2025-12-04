package router

import (
	"back/controller"
	"back/middleware"
	"back/model"
	"back/repository"
	"back/service"

	"github.com/gin-gonic/gin"
)

// InitRouter 初始化路由（核心：保证所有接口路径为 /api/v1/xxx，与前端完全匹配）
func InitRouter() *gin.Engine {
	// 1. 创建Gin引擎（默认包含Logger和Recovery中间件）
	r := gin.Default()

	// 2. 全局中间件（必须放在最前面，解决跨域+404前置拦截）
	r.Use(middleware.CorsMiddleware()) // 跨域中间件（关键：避免前端跨域+404）

	// 3. 依赖注入（核心：按 Repository → Service → Controller 顺序初始化）
	// 3.1 获取全局数据库连接（从model包的DB）
	db := model.DB
	if db == nil {
		panic("数据库连接未初始化！") // 提前报错，避免运行时nil指针
	}

	// 3.2 初始化Repository层
	userRepo := repository.NewUserRepository(db)
	categoryRepo := repository.NewCategoryRepository(db)
	articleRepo := repository.NewArticleRepository(db)

	// 3.3 初始化Service层
	userService := service.NewUserService(userRepo)
	categoryService := service.NewCategoryService(categoryRepo)
	articleService := service.NewArticleService(articleRepo, categoryRepo)

	// 3.4 初始化Controller层
	userController := controller.NewUserController(userService)
	categoryController := controller.NewCategoryController(categoryService)
	articleController := controller.NewArticleController(articleService)

	// 4. 路由分组（核心：/api/v1，与前端请求路径完全一致）
	apiV1 := r.Group("/api/v1") // 必须是 /api/v1，前端请求路径以此开头
	{
		// 4.1 公开接口（无需登录，无JWT校验）
		public := apiV1.Group("")
		{
			// ====== 用户模块公开接口 ======
			public.POST("/user/register", userController.Register) // 注册：POST /api/v1/user/register
			public.POST("/user/login", userController.Login)       // 登录：POST /api/v1/user/login

			// ====== 分类模块公开接口 ======
			public.GET("/category/list", categoryController.List) // 分类列表：GET /api/v1/category/list

			// ====== 文章模块公开接口 ======
			public.GET("/article/list", articleController.List)   // 文章列表：GET /api/v1/article/list
			public.GET("/article/:id", articleController.GetDetail) // 文章详情：GET /api/v1/article/123
		}

		// 4.2 私有接口（需登录，JWT校验）
		private := apiV1.Group("")
		private.Use(middleware.JWTMiddleware()) // JWT中间件：校验token，解析userID到上下文
		{
			// ====== 用户模块私有接口 ======
			private.GET("/user/info", userController.GetInfo)         // 获取当前用户信息：GET /api/v1/user/info
			private.PUT("/user/password", userController.ChangePassword) // 修改密码：PUT /api/v1/user/password

			// ====== 文章模块私有接口 ======
			private.POST("/article", articleController.Create)       // 发布文章：POST /api/v1/article
			private.PUT("/article/:id", articleController.Update)    // 编辑文章：PUT /api/v1/article/123
			private.DELETE("/article/:id", articleController.Delete) // 删除文章：DELETE /api/v1/article/123
		}

		// 4.3 管理员接口（需JWT+管理员权限，可选扩展）
		admin := private.Group("")
		admin.Use(middleware.AdminMiddleware()) // 管理员权限校验
		{
			// 示例：管理员可删除任意文章（扩展用）
			// admin.DELETE("/admin/article/:id", articleController.AdminDelete)
		}
	}

	// 5. 404兜底处理（可选：返回自定义404响应）
	r.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{
			"code": 404,
			"msg":  "接口不存在：" + c.Request.Method + " " + c.Request.URL.Path,
		})
	})

	return r
}