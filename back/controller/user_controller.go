package controller

import (
	"net/http"
	"strconv"

	"back/model"
	"back/service"

	"github.com/gin-gonic/gin"
)

// UserController 用户控制器
type UserController struct {
	userService *service.UserService // 依赖业务层
}

// NewUserController 初始化控制器
func NewUserController(s *service.UserService) *UserController {
	return &UserController{userService: s}
}

// Register 注册接口
func (c *UserController) Register(ctx *gin.Context) {
	// 1. 解析前端参数
	var req struct {
		Username string `json:"username" binding:"required,min=3,max=20"`
		Password string `json:"password" binding:"required,min=6,max=20"`
		Avatar   string `json:"avatar"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "参数错误：" + err.Error()})
		return
	}

	// 2. 构造用户模型
	user := &model.User{
		Username: req.Username,
		Password: req.Password, // 自动加密
		Avatar:   req.Avatar,
	}

	// 3. 调用业务层
	if err := c.userService.Register(user); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": err.Error()})
		return
	}

	// 4. 返回响应
	ctx.JSON(http.StatusOK, gin.H{"code": 200, "msg": "注册成功"})
}

// Login 登录接口
func (c *UserController) Login(ctx *gin.Context) {
	// 1. 解析参数
	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "参数错误：" + err.Error()})
		return
	}

	// 2. 调用业务层
	token, userInfo, err := c.userService.Login(req.Username, req.Password)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": err.Error()})
		return
	}

	// 3. 返回响应（令牌+用户信息）
	ctx.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "登录成功",
		"data": gin.H{
			"token": token,
			"user":  userInfo,
		},
	})
}

// GetInfo 获取用户信息
// @Summary 获取当前用户信息
// @Produce json
// @Success 200 {object} gin.H{code:int, msg:string, data:model.User}
// @Router /api/v1/user/info [get]
// @Security Bearer
func (c *UserController) GetInfo(ctx *gin.Context) {
	// 从上下文获取用户ID（JWT中间件已存入）
	userID, _ := ctx.Get("userID")

	user, err := c.userService.GetInfo(userID.(uint))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "获取用户信息失败：" + err.Error(),
		})
		return
	}

	// 隐藏密码
	user.Password = ""

	ctx.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "获取成功",
		"data": user,
	})
}

// ChangePassword 修改密码
// @Summary 修改密码
// @Produce json
// @Param data body struct{OldPassword string; NewPassword string} true "密码信息"
// @Success 200 {object} gin.H{code:int, msg:string}
// @Router /api/v1/user/password [put]
// @Security Bearer
func (c *UserController) ChangePassword(ctx *gin.Context) {
	var req struct {
		OldPassword string `json:"old_password" binding:"required"`
		NewPassword string `json:"new_password" binding:"required,min=6,max=20"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "参数错误：" + err.Error(),
		})
		return
	}

	// 从上下文获取用户ID
	userID, _ := ctx.Get("userID")

	if err := c.userService.ChangePassword(userID.(uint), req.OldPassword, req.NewPassword); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "密码修改成功",
	})
}

// GetUserById 根据ID获取用户信息（供文章关联查询）
func (c *UserController) GetUserById(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "参数错误：" + err.Error(),
		})
		return
	}

	user, err := c.userService.GetInfo(uint(id))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "获取用户信息失败：" + err.Error(),
		})
		return
	}

	// 隐藏密码
	user.Password = ""

	ctx.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "获取成功",
		"data": user,
	})
}