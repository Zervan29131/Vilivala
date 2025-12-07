package main

import (
	"back/config"
	"back/model"
	"back/router"
)

func main() {
	// 1. 加载配置
	if err := config.InitConfig(); err != nil {
		panic("配置加载失败：" + err.Error())
	}

	// 2. 初始化数据库
	if err := model.InitDB(); err != nil {
		panic("数据库连接失败：" + err.Error())
	}

	// 3. 自动迁移表结构（创建/更新表）
	model.DB.AutoMigrate(&model.User{}, &model.Article{}, &model.Category{})

	// 4. 初始化路由并启动服务
	r := router.InitRouter()
	if err := r.Run(":" + config.Config.Server.Port); err != nil {
		panic("服务启动失败：" + err.Error())
	}
}