package model

import (
	"back/config" // 注意：这里是go mod init的模块名+包路径
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// 全局DB对象，数据访问层直接调用
var DB *gorm.DB

// InitDB 初始化数据库连接
func InitDB() error {
	// 构建MySQL连接DSN（Data Source Name）
	dsn := config.Config.MySQL.Username + ":" + config.Config.MySQL.Password +
		"@tcp(" + config.Config.MySQL.Host + ":" + config.Config.MySQL.Port + ")/" +
		config.Config.MySQL.DBName + "?charset=" + config.Config.MySQL.Charset +
		"&parseTime=True&loc=Local"

	// 连接MySQL并配置
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info), // 开发环境打印SQL日志（生产可关）
	})
	if err != nil {
		return err
	}

	// 配置连接池（优化性能）
	sqlDB, _ := db.DB()
	sqlDB.SetMaxOpenConns(100)                // 最大连接数
	sqlDB.SetMaxIdleConns(20)                 // 最大空闲连接数
	sqlDB.SetConnMaxLifetime(30*time.Minute)  // 连接存活时间

	// 赋值给全局DB
	DB = db
	return nil
}