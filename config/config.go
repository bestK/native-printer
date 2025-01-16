package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	WebSocket struct {
		Port       int  `yaml:"port"`
		EnableCORS bool `yaml:"enable_cors"` // 添加跨域开关配置
	} `yaml:"websocket"`
	Log struct {
		Level  string `yaml:"level"`
		Folder string `yaml:"folder"`
	} `yaml:"log"`
}

// 创建默认配置
func NewDefaultConfig() *Config {
	cfg := &Config{}
	// 设置默认的 websocket 端口
	cfg.WebSocket.Port = 8080
	// 默认不启用跨域
	cfg.WebSocket.EnableCORS = false
	// 设置默认的日志级别和路径
	cfg.Log.Level = "info"
	cfg.Log.Folder = "logs"
	return cfg
}

// LoadConfig 从指定路径加载配置文件
func LoadConfig(path string) (*Config, error) {
	// 首先创建默认配置
	config := NewDefaultConfig()

	// 如果配置文件不存在，则创建默认配置
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return config, nil
	}

	// 读取配置文件
	data, err := os.ReadFile(path)
	if err != nil {
		return config, err
	}

	// 将配置文件内容解析到结构体
	err = yaml.Unmarshal(data, config)
	if err != nil {
		return config, err
	}

	return config, nil
}
