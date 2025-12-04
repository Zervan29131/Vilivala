package router

import (
	"back/controller"
	"back/middleware"
	"back/model"
	"back/repository"
	"back/service"

	"github.com/gin-gonic/gin"
)

// InitRouter 初始化路由
func InitRouter() *gin.Engine {
	// 1. 创建Gin引擎
	r := gin.Default()

	// 2. 全局中间件
	r.Use(middleware.CorsMiddleware()) // 跨域

	// 3. 初始化依赖（依赖注入）
	userRepo := repository.NewUserRepository(model.DB)
	userService := service.NewUserService(userRepo)
	userController := controller.NewUserController(userService)

	// 4. 路由分组
	// 公开接口（无需登录）
	publicGroup := r.Group("/api/v1")
	{
		// 用户接口
		publicGroup.POST("/user/register", userController.Register)
		publicGroup.POST("/user/login", userController.Login)
	}

	// 私有接口（需JWT认证）
	privateGroup := r.Group("/api/v1")
	privateGroup.Use(middleware.JWTMiddleware())
	{
		// 示例：获取用户信息（后续扩展）
		// privateGroup.GET("/user/info", userController.Info)
	}

	// 管理员接口（需管理员权限）
	adminGroup := privateGroup.Group("")
	adminGroup.Use(middleware.AdminMiddleware())
	{
		// 示例：管理用户（后续扩展）
		// adminGroup.GET("/user/list", userController.List)
	}

	return r
}