package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Config holds all application configuration
type Config struct {
	DBDsn              string
	ServerPort         string
	JWTSecret          string
	JWTExpires         string
	GinMode            string
	SiteURL            string
	SiteTitle          string
	SiteDescription    string
	CORSAllowedOrigins string
	FrontendDist       string
	UploadDir          string
}

// AppConfig is the global config instance
var AppConfig Config

// DB is the global database instance
var DB *gorm.DB

// Load reads environment variables from .env file and os.Environ
func Load() {
	// Load .env file (not required; falls back to OS env vars)
	if err := godotenv.Load(); err != nil {
		log.Println("[Config] No .env file found, using OS environment variables")
	}

	AppConfig = Config{
		DBDsn:              getEnv("DB_DSN", "bloguser:blogpass123@tcp(127.0.0.1:3306)/personal_blog?charset=utf8mb4&parseTime=True&loc=Local"),
		ServerPort:         getEnv("SERVER_PORT", "8080"),
		JWTSecret:          getEnv("JWT_SECRET", "vilivala-super-secret-key-change-in-production"),
		JWTExpires:         getEnv("JWT_EXPIRES_IN", "86400"),
		GinMode:            getEnv("GIN_MODE", "debug"),
		SiteURL:            getEnv("SITE_URL", "http://localhost:5173"),
		SiteTitle:          getEnv("SITE_TITLE", "Zervan的小站"),
		SiteDescription:    getEnv("SITE_DESCRIPTION", "这是存在于某处的幻想世界 · Zervanのblog"),
		CORSAllowedOrigins: getEnv("CORS_ALLOWED_ORIGINS", ""),
		FrontendDist:       getEnv("FRONTEND_DIST", ""),
		UploadDir:          getEnv("UPLOAD_DIR", "uploads"),
	}

	log.Printf("[Config] Server port: %s | GinMode: %s", AppConfig.ServerPort, AppConfig.GinMode)
}

// InitDB initializes the GORM database connection
func InitDB() {
	var err error
	gormCfg := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	}

	DB, err = gorm.Open(mysql.Open(AppConfig.DBDsn), gormCfg)
	if err != nil {
		log.Fatalf("[DB] Failed to connect to database: %v", err)
	}

	sqlDB, err := DB.DB()
	if err != nil {
		log.Fatalf("[DB] Failed to get underlying sql.DB: %v", err)
	}

	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(5)

	fmt.Println("[DB] Database connected successfully")
}

func getEnv(key, defaultVal string) string {
	if val, ok := os.LookupEnv(key); ok && val != "" {
		return val
	}
	return defaultVal
}
