# Github-Proxy-Registry

**Github 加速代理节点注册中心** —— 管理 [Github-Proxy](https://github.com/xzyuse/Github-Proxy) 实例的自动注册、凭证分发与节点调度的轻量级服务。

## 核心功能

- **自动注册**：节点实例启动时自动注册，生成 Token 和 Challenge
- **节点分发**：向已注册实例返回其他可用节点列表
- **公开查询**：无需认证的公开节点列表接口
- **Web 仪表盘**：内置前端界面，实时查看节点状态
- **定时清理**：自动清理超时离线实例

## 架构设计

```
  Github-Proxy 节点 A,B,C ──────┐
                                ▼
                    ┌─────────────────────┐
                    │ Github-Proxy-Registry│
                    │   (注册调度中心)      │
                    │                     │
                    │ • 实例注册/认证      │
                    │ • Token/Challenge   │
                    │ • 节点列表分发       │
                    │ • 定时清理超时实例   │
                    └─────────────────────┘
```

## 快速部署

### 一键构建

```bash
# Linux / macOS
./build.sh v1.0.0

# Windows (PowerShell)
.\build.ps1 -Version v1.0.0
```

### 直接运行

```bash
cd src && go build -o registry . && ./registry
```

默认监听 `0.0.0.0:8080`。

### Docker

```bash
# 构建镜像
docker build -t github-proxy-registry:latest .
docker run -d -p 8080:8080 --name registry github-proxy-registry:latest

# 或使用 Docker Compose
docker compose up -d --build
```

## 配置

`config.toml` 文件配置：

```toml
[server]
host = "0.0.0.0"
port = 8080
```

支持的环境变量：

| 环境变量 | 默认值 | 说明 |
|---------|--------|------|
| `SERVER_HOST` | `0.0.0.0` | 监听地址 |
| `SERVER_PORT` | `8080` | 监听端口 |

## API 端点

基础路径：`/api/v1`

| 方法 | 路径 | 说明 | 认证 |
|------|------|------|------|
| POST | `/register` | 实例注册 | 无 |
| GET | `/nodes?token=xxx` | 查询节点列表 | Token |
| GET | `/nodes/public` | 公开节点列表 | 无 |

### 注册示例

```bash
curl -X POST http://localhost:8080/api/v1/register \
  -H "Content-Type: application/json" \
  -d '{"domain": "github.example.com"}'
```

## 技术栈

- Go 1.23+
- [go-toml/v2](https://github.com/pelletier/go-toml)
- Vue 3 + Vite + TailwindCSS

## 项目结构

```
src/
├── main.go                    # 入口（路由注册、HTTP 服务器启动）
├── config/config.go           # 配置管理
├── frontend/                  # Vue 3 前端源码
├── handlers/                  # HTTP 处理器
│   ├── nodes.go               # 节点列表 Handler
│   └── register.go            # 注册 Handler
├── models/registry.go         # 全局注册表引用
├── utils/nodereg/registry.go  # 注册表核心逻辑
└── public/                    # 前端构建产物（嵌入二进制）
```

## License

MIT
