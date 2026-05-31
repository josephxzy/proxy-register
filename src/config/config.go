// Package config 管理调度中心配置
package config

import (
	"fmt"
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

// LoadConfig 加载配置，优先级：环境变量 > config.toml > 默认配置
func LoadConfig() *AppConfig {
	cfg := DefaultConfig()
	if data, err := os.ReadFile("config.toml"); err == nil {
		_ = toml.Unmarshal(data, cfg)
	}

	applyEnvOverrides(cfg)

	appConfigLock.Lock()
	appConfig = cfg
	appConfigLock.Unlock()

	return cfg
}

func applyEnvOverrides(cfg *AppConfig) {
	if host := os.Getenv("SERVER_HOST"); host != "" {
		cfg.Server.Host = host
	}
	if port := os.Getenv("SERVER_PORT"); port != "" {
		if p, err := parsePort(port); err == nil {
			cfg.Server.Port = p
		}
	}
}

func parsePort(portStr string) (int, error) {
	var port int
	_, err := fmt.Sscanf(portStr, "%d", &port)
	return port, err
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
