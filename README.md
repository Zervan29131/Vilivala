# Vilivala 个人博客

前后端分离的个人博客：**Vue 3 + TypeScript + Tailwind** / **Go (Gin + GORM)** / **MySQL 8.0**。

接口契约见 [docs/API.md](docs/API.md)（前后端唯一契约，改字段先改它）。

## 目录结构

```
backend/    Go 后端（:8080）
frontend/   Vue 3 前端（:5173，/api 代理到 :8080）
docs/       API.md 接口契约
docker-compose.yml   可选：Docker 启动 MySQL
```

## 环境要求

Go 1.22+、Node 18+、MySQL 8.0（本机安装或 Docker 二选一）

## 启动步骤

### 1. 数据库（二选一）

**A. 本机 MySQL（macOS / Homebrew）**

```bash
brew install mysql@8.0 && brew services start mysql@8.0
mysql -u root -e "CREATE DATABASE personal_blog CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
CREATE USER 'bloguser'@'127.0.0.1' IDENTIFIED BY 'blogpass123';
GRANT ALL ON personal_blog.* TO 'bloguser'@'127.0.0.1';"
```

**B. Docker**

```bash
docker compose up -d     # 账号 bloguser / blogpass123，库 personal_blog
```

### 2. 后端

```bash
cd backend
cp .env.example .env     # 按需修改
go run ./cmd/server      # 首次启动自动建表并写入示例数据
```

### 3. 前端

```bash
cd frontend
npm install
npm run dev              # http://localhost:5173
```

生产构建：`npm run build`（产物在 `frontend/dist`）。

## 默认管理员

| 项 | 值 |
| --- | --- |
| 后台入口 | http://localhost:5173/admin/login |
| 用户名 | `admin` |
| 密码 | `admin123` |

> 仅在 `users` 表为空时自动创建，**上线前请务必修改密码**。

## 环境变量（`backend/.env`）

| 变量 | 默认值 | 说明 |
| --- | --- | --- |
| `DB_DSN` | `bloguser:blogpass123@tcp(127.0.0.1:3306)/personal_blog?charset=utf8mb4&parseTime=True&loc=Local` | MySQL 连接串 |
| `SERVER_PORT` | `8080` | 后端端口 |
| `JWT_SECRET` | 示例值 | JWT 签名密钥，**生产必须替换** |
| `GIN_MODE` | `debug` | `debug` / `release` |

`.env` 已加入 `.gitignore`，不会提交密码。

## 跨平台迁移与部署（Linux / Windows / macOS / Docker）

项目已支持跨系统零改动迁移（所有路径使用 OS 无关处理，行尾由 `.gitattributes` 统一为 LF，Docker 构建开启多架构兼容）。

### 方式一：Docker Compose 一键部署（推荐，Linux / macOS / Windows 通用）

```bash
cp .env.example .env      # 修改 MYSQL_PASSWORD 与 JWT_SECRET
docker compose up -d --build
# 访问 http://localhost:8080 （包含 MySQL + Go 后端 + Nginx 前端与 SEO 代理）
```

### 方式二：单二进制免 Nginx 部署（支持 Linux amd64/arm64、macOS、Windows）

```bash
# 1. 一键打包前端静态产物 + 五大平台后端可执行文件（输出至 release/）
python3 tools/portable.py build-all

# 2. 在目标服务器直接用单个 Go 进程同时提供 API 与前端 SPA（含服务端轻量 SEO 注入）
FRONTEND_DIST=./frontend-dist UPLOAD_DIR=./uploads GIN_MODE=release ./vilivala-linux-amd64
```

### 数据打包与跨机迁移（数据库 + 上传图片）

```bash
# 在源机器导出完整迁移包（backups/vilivala-migration-*.tar.gz）
python3 tools/portable.py pack

# 在新机器一键恢复数据库与 uploads/ 目录
python3 tools/portable.py unpack backups/vilivala-migration-YYYYMMDD-HHMMSS.tar.gz
```

## 功能

- 前台：文章列表、搜索、分页、分类/标签筛选、文章详情（Markdown、代码高亮与复制、目录 TOC）
- 后台：登录鉴权、文章增删改查（草稿/发布）、Markdown 分屏编辑器、分类与标签管理
