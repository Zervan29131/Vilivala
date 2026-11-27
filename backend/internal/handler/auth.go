package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"vilivala-backend/config"
	"vilivala-backend/internal/middleware"
	"vilivala-backend/internal/model"
)

// Login POST /api/auth/login
func Login(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		Fail(c, http.StatusBadRequest, 400, "用户名和密码不能为空")
		return
	}
	var u model.User
	if err := config.DB.Where("username = ?", req.Username).First(&u).Error; err != nil ||
		bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(req.Password)) != nil {
		Fail(c, http.StatusUnauthorized, 401, "用户名或密码错误")
		return
	}
	token, err := middleware.GenerateToken(u.ID, u.Username)
	if err != nil {
		Fail(c, http.StatusInternalServerError, 500, "生成 Token 失败")
		return
	}
	OK(c, gin.H{"token": token, "expires_in": 86400, "user": u})
}
