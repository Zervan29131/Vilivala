package main

import (
	"back/config"
	"back/model"
	"back/router"
	"fmt"
)

func main() {
	// 1. 加载配置
	if err := config.InitConfig(); err != nil {
		panic(fmt.Sprintf("配置加载失败：%v", err))
	}

	// 2. 初始化数据库连接（必须成功）
	if err := model.InitDB(); err != nil {
		panic(fmt.Sprintf("数据库连接失败：%v", err))
	}
	if model.DB == nil {
		panic("数据库连接初始化后为nil！")
	}

	// 3. 核心：强制自动迁移表结构（加详细错误日志）
	fmt.Println("开始自动迁移表结构...")
	err := model.DB.AutoMigrate(
		&model.User{},      // 用户表
		&model.Category{},  // 分类表
		&model.Article{},   // 文章表（重点：确保这行存在）
	)
	if err != nil {
		panic(fmt.Sprintf("表结构迁移失败（关键！）：%v", err))
	}
	fmt.Println("表结构迁移成功！")

	// 4. 验证表是否存在（可选：打印确认）
	var tableExists bool
	// 检查articles表是否存在
	err = model.DB.Raw("SELECT 1 FROM information_schema.tables WHERE table_schema = ? AND table_name = ? LIMIT 1", "demo", "articles").Scan(&tableExists).Error
	if err != nil {
		panic(fmt.Sprintf("验证表存在性失败：%v", err))
	}
	if !tableExists {
		panic("AutoMigrate执行成功，但articles表仍未创建！请手动创建")
	}
	fmt.Println("验证：articles表已存在！")

	// 5. 启动服务
	fmt.Printf("启动服务，端口：%s\n", config.Config.Server.Port)
	if err := router.InitRouter().Run(":" + config.Config.Server.Port); err != nil {
		panic(fmt.Sprintf("服务启动失败：%v", err))
	}
}