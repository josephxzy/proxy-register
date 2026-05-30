// github-proxy-registry 节点调度中心
// 负责管理已注册的节点实例，分发节点列表，维护Token和挑战验证
package main

import (
	"embed"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github-proxy-registry/config"
	"github-proxy-registry/handlers"
	"github-proxy-registry/models"
	"github-proxy-registry/utils/nodereg"
)

func printHelp() {
	fmt.Printf("%s %s\n\n", ProjectName, Version)
	fmt.Println("Usage: github-proxy-registry [options]")
	fmt.Println("\nOptions:")
	fmt.Println("  --help, -h     Show this help message")
	fmt.Println("  --version, -v  Show version information")
	fmt.Printf("\nDefault address: 0.0.0.0:8080\n")
}

func printVersion() {
	fmt.Printf("%s %s\n", ProjectName, Version)
	if GitCommit != "unknown" {
		fmt.Printf("GitCommit: %s\n", GitCommit)
	}
	if BuildTime != "unknown" {
		fmt.Printf("BuildTime: %s\n", BuildTime)
	}
}

//go:embed public/*
var staticFiles embed.FS

func serveStaticFile(w http.ResponseWriter, r *http.Request, filepath string) {
	data, err := staticFiles.ReadFile(filepath)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	contentType := "text/html; charset=utf-8"
	if strings.HasSuffix(filepath, ".js") {
		contentType = "application/javascript"
	} else if strings.HasSuffix(filepath, ".css") {
		contentType = "text/css"
	} else if strings.HasSuffix(filepath, ".svg") {
		contentType = "image/svg+xml"
	} else if strings.HasSuffix(filepath, ".json") {
		contentType = "application/json"
	} else if strings.HasSuffix(filepath, ".png") {
		contentType = "image/png"
	} else if strings.HasSuffix(filepath, ".ico") {
		contentType = "image/x-icon"
	}
	w.Header().Set("Content-Type", contentType)
	w.Write(data)
}

func main() {
	// 解析命令行参数
	helpFlag := flag.Bool("help", false, "")
	helpShort := flag.Bool("h", false, "")
	versionFlag := flag.Bool("version", false, "")
	versionShort := flag.Bool("v", false, "")
	flag.Parse()

	if *helpFlag || *helpShort {
		printHelp()
		os.Exit(0)
	}

	if *versionFlag || *versionShort {
		printVersion()
		os.Exit(0)
	}

	// 加载配置
	cfg := config.LoadConfig()

	// 初始化全局注册表
	models.GlobalRegistry = nodereg.NewRegistry()

	// 定时清理超时实例
	go func() {
		ticker := time.NewTicker(5 * time.Minute)
		defer ticker.Stop()
		for range ticker.C {
			models.GlobalRegistry.Cleanup(24 * time.Hour)
		}
	}()

	// 注册路由
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1/register", handlers.Register)           // 节点注册
	mux.HandleFunc("/api/v1/nodes", handlers.GetNodes)              // 节点列表查询（Token认证）
	mux.HandleFunc("/api/v1/nodes/public", handlers.GetPublicNodes) // 公开节点列表（无需认证）

	// 前端静态文件服务
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			serveStaticFile(w, r, "public/index.html")
			return
		}
		if strings.HasPrefix(r.URL.Path, "/assets/") {
			filepath := "public" + r.URL.Path
			serveStaticFile(w, r, filepath)
			return
		}
		if r.URL.Path == "/favicon.svg" {
			serveStaticFile(w, r, "public/favicon.svg")
			return
		}
		http.NotFound(w, r)
	})

	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	printBanner(cfg)
	log.Fatal(http.ListenAndServe(addr, mux))
}
