# 江湖小报

> 新闻聚合与事实核查平台

本仓库为 **江湖小报** 的 monorepo，包含前端 App 与后端 Server 两个子项目。

## 项目结构

```
jianghu/
├── app/          # 前端（UniApp + Vue 3 + TypeScript）
└── server/       # 后端（Go + Gin + SQLite）
```

## 快速开始

### 环境要求

- Node.js >= 18
- Go >= 1.22
- CGO 支持（SQLite 需要）

### 安装依赖

```bash
# 前端依赖
cd app && npm install

# 后端无需额外安装，直接 go build
```

### 开发模式

```bash
# 分别启动
make dev-server    # 后端：http://localhost:8080
make dev-app       # 前端：H5 开发模式

# 或者一键启动（后台运行 server）
make dev
```

### 构建

```bash
make build          # 构建前后端
make build-app      # 仅构建前端 H5
make build-server   # 仅构建后端
```

### 数据库迁移

```bash
make migrate
```

## 子项目文档

- [前端 README](app/README.md)
- [后端配置示例](server/config.example.yaml)
