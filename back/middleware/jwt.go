package middleware

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"back/config"
	"back/model"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// Claims 自定义JWT载荷（扩展默认字段）
type Claims struct {
	UserID uint   `json:"user_id"` // 用户ID
	Role   string `json:"role"`    // 用户角色
	jwt.RegisteredClaims           // JWT默认字段（过期时间/签发时间等）
}

// GenerateToken 生成JWT令牌
func GenerateToken(userID uint, role string) (string, error) {
	// 过期时间：当前时间 + 配置的过期秒数
	expireTime := time.Now().Add(time.Duration(config.Config.JWT.Expire) * time.Second)

	// 构造载荷
	claims := Claims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expireTime), // 过期时间
			IssuedAt:  jwt.NewNumericDate(time.Now()), // 签发时间
			Issuer:    "my-blog",                      // 签发者（自定义）
		},
	}

	// 生成令牌（HS256算法 + 配置密钥）
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.Config.JWT.Secret))
}

// JWTMiddleware Gin中间件：校验JWT令牌
func JWTMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. 获取请求头中的令牌（格式：Bearer <token>）
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "未携带令牌"})
			c.Abort() // 终止请求
			return
		}

		// 2. 拆分令牌格式
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "令牌格式错误"})
			c.Abort()
			return
		}

		// 3. 解析令牌
		// tokenStr := parts[1]
		// claims := &Claims{}
		// token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		// 	// 校验算法
		// 	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		// 		return nil, gin.H{"code": 401, "msg": "令牌算法非法"}
		// 	}
		// 	// 返回签名密钥
		// 	return []byte(config.Config.JWT.Secret), nil
		// })
tokenStr := parts[1]
		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
			// 校验算法
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				// 修正：返回合法的error类型，而非gin.H
				return nil, errors.New("令牌算法非法") 
				// 也可以用fmt.Errorf更灵活：return nil, fmt.Errorf("不支持的签名算法：%v", token.Header["alg"])
			}
			// 返回签名密钥
			return []byte(config.Config.JWT.Secret), nil
		})

		// 4. 校验令牌有效性
		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "令牌无效/已过期"})
			c.Abort()
			return
		}

		// 5. 验证用户是否存在（防止令牌合法但用户已删除）
		var user model.User
		if err := model.DB.First(&user, claims.UserID).Error; err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "用户不存在"})
			c.Abort()
			return
		}

		// 6. 将用户信息存入Gin上下文，供后续接口使用
		c.Set("userID", claims.UserID)
		c.Set("role", claims.Role)
		c.Next() // 继续执行后续逻辑
	}
}

// AdminMiddleware 管理员权限校验（需配合JWTMiddleware使用）
func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, _ := c.Get("role")
		if role != "admin" {
			c.JSON(http.StatusForbidden, gin.H{"code": 403, "msg": "无管理员权限"})
			c.Abort()
			return
		}
		c.Next()
	}
}

// CorsMiddleware 跨域中间件（解决前端跨域）
func CorsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*") // 生产环境替换为前端域名
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type,Authorization")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		// 处理OPTIONS预请求
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusOK)
			return
		}
		c.Next()
	}
}