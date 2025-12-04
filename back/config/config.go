package config

import (
	"path/filepath"
	"runtime"

	"github.com/spf13/viper"
)

// 全局配置变量，全项目可直接调用
var Config AppConfig

// AppConfig 总配置结构体（与yaml字段对应）
type AppConfig struct {
	Server ServerConfig `mapstructure:"server"`
	MySQL  MySQLConfig  `mapstructure:"mysql"`
	JWT    JWTConfig    `mapstructure:"jwt"`
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Port string `mapstructure:"port"`
}

// MySQLConfig MySQL配置
type MySQLConfig struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	DBName   string `mapstructure:"dbname"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	Charset  string `mapstructure:"charset"`
}

// JWTConfig JWT配置
type JWTConfig struct {
	Secret string `mapstructure:"secret"`
	Expire int    `mapstructure:"expire"`
}

// InitConfig 加载配置文件到全局变量
func InitConfig() error {
	// 获取项目根目录（适配不同运行环境）
	_, filename, _, _ := runtime.Caller(0)
	rootDir := filepath.Dir(filepath.Dir(filename)) // 上一级目录（config/ → 根目录）

	// Viper配置
	viper.SetConfigName("config")       // 配置文件名（无后缀）
	viper.SetConfigType("yaml")         // 配置类型
	viper.AddConfigPath(rootDir)        // 配置文件路径（根目录）
	viper.AutomaticEnv()                // 支持环境变量覆盖（生产环境用）

	// 读取配置文件
	if err := viper.ReadInConfig(); err != nil {
		return err // 读取失败返回错误
	}

	// 解析到结构体
	if err := viper.Unmarshal(&Config); err != nil {
		return err
	}

	return nil
}