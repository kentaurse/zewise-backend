/*
Package configs - ZeWise 配置文件包
该文件用于定义 WeWise 后端服务配置文件
Copyright (c) [2024], Author(s):
- WhitePaper233<baizhiwp@gmail.com>
*/
package configs

import (
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/spf13/viper"
)

// Config 配置文件对象
type Config struct {
	// 服务器设置
	Server struct {
		// 服务器监听地址
		Host string `toml:"host"`
		// 服务器监听端口
		Port int `toml:"port"`
	} `toml:"server"`

	// Redis 设置
	Redis struct {
		// Redis 主机地址
		Host string `toml:"host"`
		// Redis 端口
		Port int `toml:"port"`
		// Redis 用户名
		Username string `toml:"username"`
		// Redis 密码
		Password string `toml:"password"`
		// Redis 数据库索引
		DB int `toml:"db"`
	} `toml:"redis"`

	// MinIO 设置
	MinIO struct {
		// MinIO 主机地址
		Host string `toml:"host"`
		// MinIO 端口
		Port int `toml:"port"`
		// MinIO access key
		Access string `toml:"access"`
		// MinIO secret key
		Secret string `toml:"secret"`
	} `toml:"minio"`

	// MongoDB 设置
	MongoDB struct {
		// MongoDB 主机地址
		Host string `toml:"host"`
		// MongoDB 端口
		Port int `toml:"port"`
		// MongoDB 数据库名称
		DBName string `toml:"dbname"`
	} `toml:"mongodb"`

	// 搜索服务设置
	SearchService struct {
		// 搜索服务主机地址
		Host string `toml:"host"`
		// 搜索服务端口
		Port int `toml:"port"`
	} `toml:"search_service"`

	// 压缩设置
	Compress struct {
		// 压缩等级
		Level compress.Level `toml:"level"`
	} `toml:"compress"`

	// 环境设置
	Env struct {
		// 环境类型 development, production
		Type string `toml:"type"`
	} `toml:"env"`
}

/*
NewConfig 创建配置文件对象

返回:
  - *Config: 配置文件对象
  - error: 错误信息
*/
func NewConfig() (*Config, error) {
	// 设置配置文件路径
	viper.SetConfigFile("./configuration.toml")

	// 读取配置文件
	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	// 解析配置文件
	config := &Config{}
	if err := viper.Unmarshal(config); err != nil {
		return nil, err
	}

	return config, nil
}
