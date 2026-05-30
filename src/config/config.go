// Package config 管理调度中心配置
package config

import (
	"os"
	"sync"

	"github.com/pelletier/go-toml/v2"
)

// AppConfig 调度中心应用配置
type AppConfig struct {
	Server struct {
		Host string `toml:"host"` // 监听地址
		Port int    `toml:"port"` // 监听端口
	} `toml:"server"`
}

var (
	appConfig     *AppConfig
	appConfigLock sync.RWMutex
)

// DefaultConfig 返回默认配置
func DefaultConfig() *AppConfig {
	return &AppConfig{
		Server: struct {
			Host string `toml:"host"`
			Port int    `toml:"port"`
		}{
			Host: "0.0.0.0",
			Port: 8080,
		},
	}
}

// LoadConfig 加载配置，优先读取 config.toml，失败则使用默认配置
func LoadConfig() *AppConfig {
	cfg := DefaultConfig()
	if data, err := os.ReadFile("config.toml"); err == nil {
		_ = toml.Unmarshal(data, cfg)
	}

	appConfigLock.Lock()
	appConfig = cfg
	appConfigLock.Unlock()

	return cfg
}

// GetConfig 获取当前配置，线程安全
func GetConfig() *AppConfig {
	appConfigLock.RLock()
	if appConfig != nil {
		cfg := *appConfig
		appConfigLock.RUnlock()
		return &cfg
	}
	appConfigLock.RUnlock()
	return DefaultConfig()
}
