package router

import (
	"back/controller" // 替换：my-blog-backend → back
	"back/middleware" // 替换：my-blog-backend → back
	"back/repository" // 替换：my-blog-backend → back
	"back/service"    // 替换：my-blog-backend → back

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// InitRouter 初始化路由（核心入口，接收数据库连接）
// 分层依赖：DB → Repository → Service → Controller
func InitRouter(db *gorm.DB) *gin.Engine {
    // 1. 初始化Gin引擎（开启默认日志和恢复中间件）
    r := gin.Default()

    // 2. 全局中间件配置
    r.Use(
        middleware.CorsMiddleware(),       // 跨域中间件
        gin.Recovery(),                    // 崩溃恢复
        gin.Logger(),                      // 请求日志
    )

    // 3. 初始化各层依赖（核心：解决依赖注入错误）
    // ========== 用户模块 ==========
    userRepo := repository.NewUserRepository(db)
    userService := service.NewUserService(userRepo)
    userController := controller.NewUserController(userService)

    // ========== 文章模块 ==========
    articleRepo := repository.NewArticleRepository(db)
    articleService := service.NewArticleService(articleRepo)
    articleController := controller.NewArticleController(articleService)

    // ========== 分类模块 ==========
    categoryRepo := repository.NewCategoryRepository(db)
    categoryService := service.NewCategoryService(categoryRepo)
    categoryController := controller.NewCategoryController(categoryService)

    // 4. API路由分组（版本控制：/api/v1）
    apiV1 := r.Group("/api/v1")
    {
        // ========== 用户相关路由 ==========
        userGroup := apiV1.Group("/user")
        {
            // 公开接口（无需登录）
            userGroup.POST("/register", userController.Register) // 注册
            userGroup.POST("/login", userController.Login)       // 登录

            // 需登录接口（JWT认证）
            userAuthGroup := userGroup.Group("/")
            userAuthGroup.Use(middleware.JWTMiddleware()) // JWT拦截器
            {
                userAuthGroup.GET("/info", userController.GetUserInfo)         // 获取用户信息
                userAuthGroup.PUT("/password", userController.ChangePassword) // 修改密码
            }
        }

        // ========== 文章相关路由 ==========
        articleGroup := apiV1.Group("/article")
        {
            // 公开接口
            articleGroup.GET("/list", articleController.GetArticleList)     // 文章列表（分页+搜索）
            articleGroup.GET("/:id", articleController.GetArticleDetail)    // 文章详情

            // 需登录接口
            articleAuthGroup := articleGroup.Group("/")
            articleAuthGroup.Use(middleware.JWTMiddleware())
            {
                articleAuthGroup.POST("/", articleController.PublishArticle)   // 发布文章
                articleAuthGroup.PUT("/:id", articleController.UpdateArticle)  // 编辑文章
                articleAuthGroup.DELETE("/:id", articleController.DeleteArticle) // 删除文章
            }
        }

        // ========== 分类相关路由 ==========
        categoryGroup := apiV1.Group("/category")
        {
            // 公开接口
            categoryGroup.GET("/list", categoryController.GetCategoryList) // 分类列表

            // 需登录接口
            categoryAuthGroup := categoryGroup.Group("/")
            categoryAuthGroup.Use(middleware.JWTMiddleware())
            {
                categoryAuthGroup.POST("/", categoryController.CreateCategory) // 创建分类
                categoryAuthGroup.PUT("/:id", categoryController.UpdateCategory) // 编辑分类
                categoryAuthGroup.DELETE("/:id", categoryController.DeleteCategory) // 删除分类
            }
        }
    }

    // 5. 404路由（兜底）
    r.NoRoute(func(c *gin.Context) {
        c.JSON(404, gin.H{
            "code": 404,
            "msg":  "接口不存在",
            "data": nil,
        })
    })

    return r
}