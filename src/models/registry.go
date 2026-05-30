// Package models 存放调度中心的数据模型
package models

import "github-proxy-registry/utils/nodereg"

// GlobalRegistry 全局节点注册表，由 main.go 初始化，各 handler 共享
var GlobalRegistry *nodereg.Registry
