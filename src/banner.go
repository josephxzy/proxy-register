package main

import (
	"fmt"

	"github-proxy-registry/config"
)

func printBanner(cfg *config.AppConfig) {
	fmt.Println()
	fmt.Println("============================================")
	fmt.Printf(" 项目: %s\n", ProjectName)
	fmt.Printf(" 仓库: %s\n", ProjectURL)
	fmt.Printf(" 版本: %s\n", Version)
	fmt.Printf(" 构建: %s\n", BuildTime)
	fmt.Println("--------------------------------------------")
	fmt.Printf(" 监听: %s:%d\n", cfg.Server.Host, cfg.Server.Port)
	fmt.Println("============================================")
	fmt.Println()
}
