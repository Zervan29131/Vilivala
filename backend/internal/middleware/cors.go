package middleware

import (
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"vilivala-backend/config"
)

// CORS returns a configured CORS middleware that allows the Vue dev server or configured origins
func CORS() gin.HandlerFunc {
	origins := []string{
		"http://localhost:5173",
		"http://localhost:3000",
		"http://127.0.0.1:5173",
		"http://localhost:8080",
	}
	if raw := strings.TrimSpace(config.AppConfig.CORSAllowedOrigins); raw != "" {
		custom := []string{}
		for _, part := range strings.Split(raw, ",") {
			if o := strings.TrimSpace(part); o != "" {
				custom = append(custom, o)
			}
		}
		if len(custom) > 0 {
			origins = custom
		}
	}
	cfg := cors.Config{
		AllowOrigins:     origins,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "Accept"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}
	return cors.New(cfg)
}
