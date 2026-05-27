# EPay Go 部署指南

> 在 `D:\SRMTJ\LZClaw` 工作区内，推荐先通过 `D:\SRMTJ\LZClaw\LZDocker\docker-compose.yml` 启动统一 PostgreSQL/Redis 依赖，再启动本项目服务。本文档中的本仓库 `docker-compose` 方式保留为历史兼容路径。

## 环境要求

- Docker 20.10+
- Docker Compose 2.0+
- 2GB+ 内存
- 10GB+ 磁盘空间

## 快速开始

### 1. 克隆项目

```bash
git clone <repository-url>
cd epay-go
```

### 2. 配置环境变量

```bash
cp .env.example .env
# 编辑 .env 文件，配置必要的环境变量
```

### 3. 启动服务

```bash
# 开发环境
docker-compose up -d

# 生产环境
docker-compose -f docker-compose.yml -f docker-compose.prod.yml up -d
```

### 4. 访问服务

- 前端: http://localhost
- 后端 API: http://localhost:8080
- 管理后台: http://localhost/admin/login
- 商户中心: http://localhost/merchant/login

## 常用命令

```bash
# 查看日志
docker-compose logs -f

# 查看服务状态
docker-compose ps

# 停止服务
docker-compose down

# 重新构建
docker-compose build --no-cache

# 进入容器
docker-compose exec backend sh
docker-compose exec postgres psql -U epay
```

## 数据备份

```bash
# 备份数据库
docker-compose exec postgres pg_dump -U epay epay > backup.sql

# 恢复数据库
cat backup.sql | docker-compose exec -T postgres psql -U epay epay
```

## 更新部署

```bash
# 拉取最新代码
git pull

# 重新构建并启动
docker-compose build
docker-compose up -d
```

## 故障排查

### 后端无法连接数据库

检查 PostgreSQL 是否就绪：
```bash
docker-compose logs postgres
```

### 前端无法访问后端 API

检查 Nginx 代理配置和后端服务状态：
```bash
docker-compose logs frontend
docker-compose logs backend
```
