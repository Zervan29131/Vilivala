package controller

import (
	"back/service"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// 1. 定义UserController结构体（依赖UserService）
type UserController struct {
	userService *service.UserService // 注入UserService
}

// 2. 构造函数（关键：修复语法错误！）
// 错误原因：之前缺少func关键字/参数/返回值格式错误
func NewUserController(userService *service.UserService) *UserController {
	return &UserController{
		userService: userService,
	}
}

// 3. 注册请求参数结构体
type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=2,max=20"` // 用户名2-20位
	Password string `json:"password" binding:"required,min=6"`       // 密码至少6位
	Avatar   string `json:"avatar"`                                  // 头像可选
}

// 4. 注册接口（修复语法错误，添加完整逻辑）
func (c *UserController) Register(ctx *gin.Context) {
	var req RegisterRequest
	// 绑定JSON参数并捕获错误
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "参数解析失败：" + err.Error(),
			"data": nil,
		})
		println("注册参数绑定错误：", err.Error()) // 终端打印错误
		return
	}

	// 密码加密（bcrypt）
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "密码加密失败：" + err.Error(),
			"data": nil,
		})
		return
	}

	// 调用Service层注册逻辑
	success, msg := c.userService.Register(req.Username, string(hashPassword), req.Avatar)
	if success {
		ctx.JSON(http.StatusOK, gin.H{
			"code": 200,
			"msg":  msg,
			"data": nil,
		})
	} else {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  msg,
			"data": nil,
		})
	}
}

// 5. 登录请求参数结构体
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// 6. 登录接口（补充完整逻辑，无语法错误）
func (c *UserController) Login(ctx *gin.Context) {
	var req LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "参数解析失败：" + err.Error(),
			"data": nil,
		})
		return
	}

	// 调用Service层登录逻辑
	user, token, msg := c.userService.Login(req.Username, req.Password)
	if user != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code": 200,
			"msg":  msg,
			"data": gin.H{
				"token":    token,
				"username": user.Username,
				"avatar":   user.Avatar,
				"id":       user.ID,
			},
		})
	} else {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"code": 401,
			"msg":  msg,
			"data": nil,
		})
	}
}

// 7. 获取用户信息接口
func (c *UserController) GetUserInfo(ctx *gin.Context) {
	// 从上下文获取userID（JWT解析后存入）
	userID, _ := ctx.Get("userID")
	user, msg := c.userService.GetUserInfo(userID.(uint))
	if user != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code": 200,
			"msg":  msg,
			"data": gin.H{
				"id":       user.ID,
				"username": user.Username,
				"avatar":   user.Avatar,
			},
		})
	} else {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  msg,
			"data": nil,
		})
	}
}