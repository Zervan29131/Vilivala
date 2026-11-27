# API 接口契约文档

> **版本**：v1.1（本轮实现：分类树、文章增强、独立页面、图片上传、归档；账号体系/评论等见文末「v1.2 规划」，**暂缓，尚未实现**）  
> **更新日期**：2026-11-25  
> **说明**：本文件是前后端唯一接口契约，任何字段变更须先修改此文档，再同步修改代码。

---

## 统一响应格式

```json
{
  "code": 0,
  "message": "ok",
  "data": {}
}
```

| 字段 | 类型 | 说明 |
|------|------|------|
| `code` | int | `0` 表示成功，非 `0` 表示错误（见错误码表） |
| `message` | string | 状态描述，成功为 `"ok"`，失败为错误原因 |
| `data` | any | 业务数据，失败时为 `null` |

### 常用错误码

| code | 含义 |
|------|------|
| 0 | 成功 |
| 400 | 请求参数错误 |
| 401 | 未认证 / Token 无效或过期 |
| 403 | 无权限 |
| 404 | 资源不存在 |
| 500 | 服务器内部错误 |

---

## Base URL

- **开发环境**：`http://localhost:8080`
- **前端代理**：Vite 将 `/api/*` 代理到 `http://localhost:8080`

所有接口路径均以 `/api` 开头。

---

## 分页公共参数

需要分页的接口均支持以下 Query 参数：

| 参数 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| `page` | int | 1 | 当前页码（从 1 开始） |
| `page_size` | int | 10 | 每页条数（最大 50） |

分页响应统一包含：

```json
{
  "code": 0,
  "message": "ok",
  "data": {
    "list": [],
    "total": 100,
    "page": 1,
    "page_size": 10
  }
}
```

---

## 一、公开接口（无需鉴权）

### 1.1 获取已发布文章列表

**GET** `/api/posts`

**Query 参数**

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| `page` | int | 否 | 页码，默认 1 |
| `page_size` | int | 否 | 每页条数，默认 10 |
| `category_id` | int | 否 | 按分类 ID 筛选 |
| `tag_id` | int | 否 | 按标签 ID 筛选 |
| `keyword` | string | 否 | 关键词搜索（匹配 title 和 summary） |

**响应示例**

```json
{
  "code": 0,
  "message": "ok",
  "data": {
    "list": [
      {
        "id": 1,
        "title": "Go 语言入门指南",
        "slug": "go-getting-started",
        "summary": "从零开始学习 Go 语言的基础知识...",
        "cover_image": "https://example.com/cover.jpg",
        "status": "published",
        "view_count": 128,
        "created_at": "2026-11-01T10:00:00+08:00",
        "updated_at": "2026-11-10T15:30:00+08:00",
        "category": {
          "id": 1,
          "name": "技术笔记",
          "slug": "tech-notes"
        },
        "tags": [
          { "id": 1, "name": "Go", "slug": "go" },
          { "id": 2, "name": "后端", "slug": "backend" }
        ]
      }
    ],
    "total": 42,
    "page": 1,
    "page_size": 10
  }
}
```

---

### 1.2 获取文章详情

**GET** `/api/posts/:slug`

> 访问时自动将 `view_count` +1。

**路径参数**

| 参数 | 类型 | 说明 |
|------|------|------|
| `slug` | string | 文章 slug（URL 友好标识符） |

**响应示例**

```json
{
  "code": 0,
  "message": "ok",
  "data": {
    "id": 1,
    "title": "Go 语言入门指南",
    "slug": "go-getting-started",
    "summary": "从零开始学习 Go 语言的基础知识...",
    "content": "# Go 语言入门\n\n## 安装\n\n```bash\nbrew install go\n```\n\n...",
    "cover_image": "https://example.com/cover.jpg",
    "status": "published",
    "view_count": 129,
    "created_at": "2026-11-01T10:00:00+08:00",
    "updated_at": "2026-11-10T15:30:00+08:00",
    "category": {
      "id": 1,
      "name": "技术笔记",
      "slug": "tech-notes"
    },
    "tags": [
      { "id": 1, "name": "Go", "slug": "go" },
      { "id": 2, "name": "后端", "slug": "backend" }
    ]
  }
}
```

---

### 1.3 获取所有分类及文章数

**GET** `/api/categories`

**响应示例**

```json
{
  "code": 0,
  "message": "ok",
  "data": [
    {
      "id": 1,
      "name": "技术笔记",
      "slug": "tech-notes",
      "post_count": 15,
      "created_at": "2026-10-01T00:00:00+08:00"
    },
    {
      "id": 2,
      "name": "生活随笔",
      "slug": "life-essay",
      "post_count": 8,
      "created_at": "2026-10-01T00:00:00+08:00"
    }
  ]
}
```

---

### 1.4 获取所有标签

**GET** `/api/tags`

**响应示例**

```json
{
  "code": 0,
  "message": "ok",
  "data": [
    { "id": 1, "name": "Go", "slug": "go", "post_count": 10 },
    { "id": 2, "name": "后端", "slug": "backend", "post_count": 7 },
    { "id": 3, "name": "Vue", "slug": "vue", "post_count": 5 }
  ]
}
```

---

## 二、鉴权接口

### 2.1 登录（管理员与普通用户共用，v1.2 起见第五章，此处为 v1.0 原始说明）

**POST** `/api/auth/login`

**请求体（JSON）**

```json
{
  "username": "admin",
  "password": "admin123"
}
```

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| `username` | string | 是 | 用户名 |
| `password` | string | 是 | 明文密码（后端 bcrypt 验证） |

**响应示例（成功）**

```json
{
  "code": 0,
  "message": "ok",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "expires_in": 86400,
    "user": {
      "id": 1,
      "username": "admin",
      "nickname": "博主",
      "avatar": ""
    }
  }
}
```

**响应示例（失败）**

```json
{
  "code": 401,
  "message": "用户名或密码错误",
  "data": null
}
```

---

## 三、后台管理接口（需 JWT 鉴权，且 role = admin；普通用户访问返回 403）

> 所有请求必须携带 Header：`Authorization: Bearer <token>`

### 3.1 文章管理

#### 3.1.1 获取所有文章列表（含草稿）

**GET** `/api/admin/posts`

**Query 参数**

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| `page` | int | 否 | 页码，默认 1 |
| `page_size` | int | 否 | 每页条数，默认 10 |
| `status` | string | 否 | 筛选状态：`draft` / `published` |
| `keyword` | string | 否 | 关键词搜索 |

**响应示例**（同 1.1，但 `list` 中含草稿文章）

---

#### 3.1.2 获取单篇文章详情（后台用，不增加浏览量）

**GET** `/api/admin/posts/:id`

| 参数 | 类型 | 说明 |
|------|------|------|
| `id` | int | 文章 ID |

**响应**：同 1.2 的 data 结构

---

#### 3.1.3 创建文章

**POST** `/api/admin/posts`

**请求体（JSON）**

```json
{
  "title": "新文章标题",
  "slug": "new-article-slug",
  "summary": "这是摘要...",
  "content": "# 正文\n\n这是 Markdown 内容...",
  "cover_image": "",
  "status": "draft",
  "category_id": 1,
  "tag_ids": [1, 2]
}
```

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| `title` | string | 是 | 文章标题 |
| `slug` | string | 是 | URL slug（需唯一） |
| `summary` | string | 否 | 文章摘要 |
| `content` | string | 是 | Markdown 正文 |
| `cover_image` | string | 否 | 封面图 URL |
| `status` | string | 是 | `draft` 或 `published` |
| `category_id` | int | 否 | 所属分类 ID |
| `tag_ids` | []int | 否 | 标签 ID 数组 |

**响应示例**

```json
{
  "code": 0,
  "message": "ok",
  "data": {
    "id": 10,
    "title": "新文章标题",
    "slug": "new-article-slug",
    "status": "draft",
    "created_at": "2026-11-19T11:00:00+08:00"
  }
}
```

---

#### 3.1.4 更新文章

**PUT** `/api/admin/posts/:id`

| 参数 | 类型 | 说明 |
|------|------|------|
| `id` | int | 文章 ID |

**请求体**：同 3.1.3（所有字段均可选，仅传需要更新的字段）

**响应示例**

```json
{
  "code": 0,
  "message": "ok",
  "data": null
}
```

---

#### 3.1.5 删除文章

**DELETE** `/api/admin/posts/:id`

**响应示例**

```json
{
  "code": 0,
  "message": "ok",
  "data": null
}
```

---

#### 3.1.6 切换文章状态

**PATCH** `/api/admin/posts/:id/status`

**请求体（JSON）**

```json
{
  "status": "published"
}
```

**响应示例**

```json
{
  "code": 0,
  "message": "ok",
  "data": null
}
```

---

### 3.2 分类管理

#### 3.2.1 创建分类

**POST** `/api/admin/categories`

**请求体（JSON）**

```json
{
  "name": "新分类",
  "slug": "new-category"
}
```

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| `name` | string | 是 | 分类名称 |
| `slug` | string | 是 | URL slug（需唯一） |

**响应示例**

```json
{
  "code": 0,
  "message": "ok",
  "data": {
    "id": 3,
    "name": "新分类",
    "slug": "new-category",
    "created_at": "2026-11-19T11:00:00+08:00"
  }
}
```

---

#### 3.2.2 更新分类

**PUT** `/api/admin/categories/:id`

**请求体**：同 3.2.1

**响应**：`{ "code": 0, "message": "ok", "data": null }`

---

#### 3.2.3 删除分类

**DELETE** `/api/admin/categories/:id`

> 若分类下存在文章，返回 `400` 并提示"请先移除该分类下的文章"。

**响应**：`{ "code": 0, "message": "ok", "data": null }`

---

### 3.3 标签管理

#### 3.3.1 创建标签

**POST** `/api/admin/tags`

**请求体（JSON）**

```json
{
  "name": "新标签",
  "slug": "new-tag"
}
```

**响应示例**

```json
{
  "code": 0,
  "message": "ok",
  "data": {
    "id": 4,
    "name": "新标签",
    "slug": "new-tag",
    "created_at": "2026-11-19T11:00:00+08:00"
  }
}
```

---

#### 3.3.2 更新标签

**PUT** `/api/admin/tags/:id`

**请求体**：同 3.3.1

**响应**：`{ "code": 0, "message": "ok", "data": null }`

---

#### 3.3.3 删除标签

**DELETE** `/api/admin/tags/:id`

**响应**：`{ "code": 0, "message": "ok", "data": null }`

---

## 四、数据模型字段类型参考

### Post

| 字段 | Go 类型 | JSON 类型 | 说明 |
|------|---------|-----------|------|
| `id` | uint | number | 自增主键 |
| `title` | string | string | 文章标题（非空） |
| `slug` | string | string | URL 友好标识（唯一） |
| `summary` | string | string | 摘要 |
| `content` | string | string | Markdown 正文（LONGTEXT） |
| `cover_image` | string | string | 封面图 URL |
| `status` | string | string | `draft` / `published` |
| `view_count` | uint | number | 浏览量 |
| `category_id` | uint | number | 分类 ID（可为 null） |
| `created_at` | time.Time | string (ISO8601) | 创建时间 |
| `updated_at` | time.Time | string (ISO8601) | 更新时间 |

### Category

| 字段 | Go 类型 | JSON 类型 | 说明 |
|------|---------|-----------|------|
| `id` | uint | number | 自增主键 |
| `name` | string | string | 分类名称（唯一） |
| `slug` | string | string | URL slug（唯一） |
| `created_at` | time.Time | string (ISO8601) | 创建时间 |

### Tag

| 字段 | Go 类型 | JSON 类型 | 说明 |
|------|---------|-----------|------|
| `id` | uint | number | 自增主键 |
| `name` | string | string | 标签名称（唯一） |
| `slug` | string | string | URL slug（唯一） |
| `created_at` | time.Time | string (ISO8601) | 创建时间 |

### User

| 字段 | Go 类型 | JSON 类型 | 说明 |
|------|---------|-----------|------|
| `id` | uint | number | 自增主键 |
| `username` | string | string | 用户名（唯一） |
| `password_hash` | string | — | bcrypt 哈希（不返回给前端） |
| `nickname` | string | string | 显示昵称 |
| `avatar` | string | string | 头像 URL |
| `created_at` | time.Time | string (ISO8601) | 创建时间 |

---

# v1.1 内容迁移与分类树（本轮实现）

> 背景：把原 VuePress 站点（`Vilivala/`，约 255 篇文章、三级分类）迁入本系统。以下为**对既有接口的变更与新增**，未提及的接口不变。详细迁移规则见 [MIGRATION.md](./MIGRATION.md)。
> 前端首页与主题（暗黑网格背景、三大板块入口等）为纯前端内容，不涉及接口，见 MIGRATION.md 第六节。

## A. 分类树（最多三级）

### A.1 分类对象 `Category`（变更：新增字段）

```json
{
  "id": 7,
  "name": "social",
  "slug": "fun-social",
  "parent_id": 3,
  "depth": 2,
  "sort": 0,
  "post_count": 14,
  "total_count": 14,
  "path": "fun / social",
  "created_at": "2026-11-25T10:00:00+08:00"
}
```

| 字段 | 类型 | 说明 |
|------|------|------|
| `parent_id` | int | 父分类 ID，`0` 表示顶级（板块） |
| `depth` | int | 层级，顶级为 `1`，最大 `3` |
| `sort` | int | 同级排序，升序，相同按 `id` |
| `post_count` | int | **直属**该分类的已发布文章数 |
| `total_count` | int | 该分类及**全部子孙**分类的已发布文章数（去重） |
| `path` | string | 从顶级到自身的名称，以 ` / ` 连接 |

### A.2 分类列表（变更）

**GET** `/api/categories`

| 参数 | 默认 | 说明 |
|------|------|------|
| `tree` | `0` | `1` 返回嵌套树（每个节点多一个 `children` 数组，按 `sort,id` 排序）；`0` 返回扁平数组（按 `sort,id`），向后兼容 |

扁平：`data: [Category, ...]`；树：`data: [Category & { "children": [...] }, ...]`。

### A.3 按分类筛选文章（变更）

`GET /api/posts?category_id=3` 现在**包含该分类全部子孙分类**下的已发布文章。

### A.4 文章中的分类（变更）

文章对象的 `category` 增加 `path` 面包屑数组（顶级 → 当前）：

```json
"category": {
  "id": 7, "name": "social", "slug": "fun-social",
  "path": [ { "id": 3, "name": "fun", "slug": "fun" }, { "id": 7, "name": "social", "slug": "fun-social" } ]
}
```

### A.5 后台分类管理（变更：`POST` / `PUT` 新增字段）

请求体增加：

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| `parent_id` | int | 否 | 默认 `0`（顶级） |
| `sort` | int | 否 | 默认 `0` |

规则与错误（均返回 `400`）：

- 层级超过 3 级 → `分类最多三级`
- 同一父分类下名称重复 → `同级分类名称已存在`（名称**不再全局唯一**，如 `随笔` 可同时存在于不同父级下）
- `slug` 全局唯一 → `slug 已存在`
- `PUT` 时 `parent_id` 指向自己或自己的子孙 → `不能移动到自身或子分类下`；移动后任一子孙层级会超过 3 → `分类最多三级`
- `DELETE` 时存在子分类或文章 → `请先移除子分类和文章`

---

## B. 文章增强

### B.1 文章对象（变更：新增字段）

| 字段 | 类型 | 说明 |
|------|------|------|
| `is_top` | bool | 是否置顶 |
| `word_count` | int | 正文字数（去除空白后的字符数，保存时计算） |
| `reading_time` | int | 预计阅读分钟数 = `max(1, ceil(word_count / 400))` |

### B.2 排序（变更）

所有文章列表（公开与后台）排序改为：`is_top DESC, created_at DESC, id DESC`。

### B.3 创建 / 更新文章（变更：请求体新增）

| 字段 | 类型 | 说明 |
|------|------|------|
| `is_top` | bool | 可选，默认 `false` |
| `created_at` | string (RFC3339) | 可选，**仅管理员可设**，用于迁移时保留原始发布日期；缺省取当前时间。`PUT` 同样允许修改 |

### B.4 文章详情增加上一篇/下一篇（变更）

`GET /api/posts/:slug` 的 `data` 增加：

```json
"prev": { "title": "上一篇标题", "slug": "prev-slug" },
"next": { "title": "下一篇标题", "slug": "next-slug" }
```

按 B.2 的排序，仅在**已发布**文章中取相邻项；`prev` 为更早的一篇，`next` 为更新的一篇；没有则为 `null`。

### B.5 搜索范围（变更）

`keyword` 现在匹配 **标题、摘要、正文**（原为标题、摘要）。

### B.6 归档（新增）

**GET** `/api/archives`（公开）— 已发布文章按年月分组，年份倒序、月份倒序、组内按 B.2 排序。

```json
{
  "code": 0, "message": "ok",
  "data": [
    {
      "year": 2022,
      "count": 9,
      "months": [
        { "month": 2, "posts": [ { "id": 1, "title": "200元内真无线蓝牙耳机对比测评", "slug": "7634a", "created_at": "2022-02-22T13:21:32+08:00" } ] }
      ]
    }
  ]
}
```

---

## C. 独立页面（关于、友情链接、日志、赞助）

页面为一段 Markdown，由管理员维护，**取代**文末 v1.2 规划中的「关于页」。迁移时预置 4 个：`about`、`links`、`changelog`、`sponsor`。

### C.1 页面列表（公开）

**GET** `/api/pages` → `data: [ { "slug": "about", "title": "关于", "updated_at": "..." }, ... ]`（不含 `content`）

### C.2 页面详情（公开）

**GET** `/api/pages/:slug` → `data: { "slug": "about", "title": "关于", "content": "# ...", "updated_at": "..." }`；不存在返回 `404`。

### C.3 后台页面管理（需 JWT）

| 接口 | 说明 |
|------|------|
| **GET** `/api/admin/pages` | 同 C.1 |
| **PUT** `/api/admin/pages/:slug` | 新增或覆盖（upsert），请求 `{ "title": "关于", "content": "..." }`。`slug` 须匹配 `^[a-z0-9-]{1,64}$`，`title` 1–64 字符，`content` 可为空；响应 `data: null` |
| **DELETE** `/api/admin/pages/:slug` | 删除；响应 `data: null` |

---

## D. 图片上传

### D.1 上传（需 JWT，管理员）

**POST** `/api/admin/upload`，`multipart/form-data`，字段 `file`。

| 约束 | 值 |
|------|------|
| 类型 | `image/jpeg`、`image/png`、`image/gif`、`image/webp`，按**文件内容**检测，不信任扩展名 |
| 大小 | 5 MB 以内，超出返回 `413` |
| 存储 | `backend/uploads/YYYYMM/<随机名>.<ext>` |

```json
{ "code": 0, "message": "ok", "data": { "url": "/uploads/202611/3f9a1c2e5b7d.png", "size": 20480, "content_type": "image/png" } }
```

静态访问 **GET** `/uploads/*filepath`（公开，**无** `/api` 前缀）。前端 Vite 需同时代理 `/uploads`。`url` 可直接用于 `cover_image` 或 Markdown 图片。
迁移时同样通过此接口上传文章内引用的本地图片。

---

## E. 数据模型变更

| 表 | 变更 |
|----|------|
| `categories` | 新增 `parent_id`(默认 0, 索引)、`sort`(默认 0)；`name` **去掉全局唯一**，改为 (`parent_id`,`name`) 联合唯一；`slug` 仍全局唯一 |
| `posts` | 新增 `is_top`(bool, 默认 false)、`word_count`(int, 默认 0)；`created_at` 允许由管理员指定；排序索引 (`status`,`is_top`,`created_at`) |
| `pages`（新） | `id`、`slug`(唯一)、`title`(64)、`content`(longtext)、`updated_at` |

> 升级策略：启动时 `AutoMigrate` 自动加列；已有分类的 `parent_id`/`sort` 取默认 `0`，`depth` 由服务端运行时计算，不入库；已有文章的 `word_count` 在启动时一次性回填。

## F. 本轮错误码

沿用既有；新增 `413`（上传过大）。

---

# v1.1.1 体验、搜索与 SEO 升级（本轮实现）

> 对应 `docs/UPGRADE_PLAN.md` 的 P0 + P1（不含 v1.2 评论与账号体系）。

## A. 归档排序修正

**GET** `/api/archives` 组内排序改为严格按发布时间倒序：`created_at DESC, id DESC`（不再因 `is_top` 打乱时间线）。

## B. 登录限流

**POST** `/api/auth/login` 增加按客户端 IP 的固定窗口限流：每 5 分钟最多 15 次请求；超出返回 HTTP `429`：

```json
{ "code": 429, "message": "登录尝试过于频繁，请稍后再试", "data": null }
```

## C. 相关文章推荐（变更：`GET /api/posts/:slug`）

`GET /api/posts/:slug` 的 `data` 新增 `related` 数组（最多 4 篇，按相关性与发布时间降序；仅已发布文章、排除当前文章）：

```json
"related": [
  {
    "id": 12,
    "title": "相关文章标题",
    "slug": "related-slug",
    "summary": "摘要...",
    "reading_time": 4,
    "created_at": "2022-02-20T10:00:00+08:00"
  }
]
```

相关性评分：每个共同标签 `+2`，同一分类 `+1`；若无共同标签或分类则按最新已发布文章补足。

## D. 全文搜索（新增）

**GET** `/api/search`（公开）

| 参数 | 默认 | 说明 |
|------|------|------|
| `q` | — | 搜索关键词（必填，去空白后为空返回空列表） |
| `limit` | `12` | 返回条数，范围 `1–30` |

数据库在 `posts(title, summary, content)` 建立 MySQL `FULLTEXT`（`WITH PARSER ngram`）索引；短词或全文检索无命中时自动降级为 `LIKE` 匹配。

```json
{
  "code": 0,
  "message": "ok",
  "data": {
    "list": [
      {
        "id": 1,
        "title": "Git实际操作",
        "slug": "c6d7d3",
        "snippet": "……使用 git commit 提交更改……",
        "category_name": "博客开发指南",
        "created_at": "2022-02-24T10:00:00+08:00"
      }
    ],
    "total": 1
  }
}
```

## E. SEO、订阅与健康检查端点（公开）

| 路径 | Content-Type | 说明 |
|------|--------------|------|
| **GET** `/api/health` | `application/json` | 健康检查，返回 `{ "code": 0, "message": "ok", "data": { "status": "ok" } }` |
| **GET** `/robots.txt` | `text/plain` | 允许抓取公开页，禁止 `/admin`，声明 `Sitemap` 地址 |
| **GET** `/sitemap.xml` | `application/xml` | 动态站点地图（首页、归档、分类、独立页、全部已发布文章） |
| **GET** `/rss.xml` | `application/rss+xml` | 标准 RSS 2.0 订阅源（最近 30 篇已发布文章） |

> 轻量元信息注入：当后端配置 `FRONTEND_DIST`（指向构建产物 `dist`）时，服务端在返回 `index.html` 时按路径为 `/posts/:slug` 与 `/pages/:slug` 注入动态 `<title>`、`description`、Open Graph / Twitter Card、`canonical` 及 `application/ld+json`（`BlogPosting`）；前端路由同步更新客户端 meta。

## F. 终端起始页（纯前端 + 热榜排序参数）

前端 `/terminal` 交互式终端风格浏览器起始页升级为完整工作台（虚拟目录树挂载三级分类、命令折叠、历史重放、壁纸自定义及 `hot` / `timing` / `varbook` / `shortcut` 内嵌交互组件）。

为此，`GET /api/posts` 新增可选查询参数 `sort`：
- 缺省：保持既有 `is_top DESC, created_at DESC, id DESC`
- `sort=views`：按阅读量倒序 `view_count DESC, id DESC`（供终端 `hot` 热文榜组件拉取 Top 10）

---

# v1.1.2 阅读精修、个人站闭环与后台仪表盘（本轮实现）

## A. 管理员草稿预览（变更：`GET /api/posts/:slug`）

- 游客或普通请求：仅返回 `status = "published"` 的文章（未发布返回 `404`）。
- 若请求头携带**有效管理员 JWT**（`Authorization: Bearer <token>`）：允许读取 `status = "draft"` 的文章用于前台真实样式预览；读取草稿时**不累加 `view_count`**，且响应增加布尔标记 `"is_preview": true`。Token 缺失或过期时不报错，按游客逻辑处理。

## B. 站点公开配置与全站统计（新增）

**GET** `/api/site-config`（公开）

返回站点动态配置（首次未配置时自动返回默认值）以及实时全站文章统计 `stats`：

```json
{
  "code": 0,
  "message": "ok",
  "data": {
    "site_title": "Zervan的小站",
    "site_description": "Zervanのblog",
    "hero_title": "这是存在于某处的幻想世界",
    "hero_tagline": "Be the change you want to see in the World",
    "subtitles": [
      "这是存在于某处的幻想世界",
      "Be the change you want to see in the World",
      "记录技术、思想与未来"
    ],
    "blogger_name": "Zervan",
    "blogger_slogan": "天下最普通的人之一",
    "blogger_avatar": "/img/console.jpg",
    "notice_enabled": true,
    "notice_mode": "custom",
    "notice_text": "Stay hungry, stay foolish · 欢迎来到 Zervan 的小站",
    "reward_enabled": true,
    "reward_wording": "觉得文章有帮助？请我喝杯咖啡吧 ~",
    "copyright_enabled": true,
    "copyright_notice": "本博客所有文章除特别声明外，著作权归作者所有。转载请注明出处！",
    "since_year": 2021,
    "icp_number": "",
    "icp_url": "https://beian.miit.gov.cn/",
    "gongan_number": "",
    "gongan_url": "https://www.beian.gov.cn/",
    "socials": [
      { "title": "微博", "link": "https://weibo.com/zervan" }
    ],
    "stats": {
      "published_posts": 255,
      "total_words": 480000,
      "total_views": 1024
    }
  }
}
```

## C. 后台更新站点配置（新增，需 JWT）

**PUT** `/api/admin/site-config`

- 请求体与 B 中配置字段一致（忽略 `stats`）。
- 校验规则：`site_title` 必填（1–64 字符）；`notice_mode` 仅允许 `custom` 或 `hitokoto`；外链地址仅允许空字符串、`http://`、`https://`、`mailto:` 或以 `/` 开头的站内路径，非法协议返回 `400`。
- 持久化至 `site_settings` 表（`id = 1`），响应返回最新完整配置（同 B）。

## D. 后台数据仪表盘（新增，需 JWT）

**GET** `/api/admin/dashboard`

```json
{
  "code": 0,
  "message": "ok",
  "data": {
    "total_posts": 255,
    "published_posts": 255,
    "draft_posts": 0,
    "categories": 25,
    "tags": 52,
    "pages": 4,
    "total_words": 480000,
    "total_views": 1024,
    "top_posts": [
      { "id": 1, "title": "Git实际操作", "slug": "c6d7d3", "status": "published", "view_count": 120, "word_count": 3200, "created_at": "..." }
    ],
    "recent_posts": [
      { "id": 2, "title": "...", "slug": "...", "status": "published", "updated_at": "..." }
    ],
    "category_distribution": [
      { "id": 1, "name": "code", "slug": "code", "total_count": 36 }
    ]
  }
}
```

---

# v1.2 规划（⏸ 暂缓，尚未实现，勿按此开发）

> 该部分已评审但**搁置**，待内容迁移与界面重做完成后再启动。其中「关于页」「图片上传」已被上面的 v1.1 独立页面 / 上传接口取代。

> 既有接口除下列「变更」外保持不变。统一响应格式、分页参数、`Authorization: Bearer <token>` 方式同前文。
> **暗黑模式**与**界面主题**为纯前端内容，见第十一章。

## 变更总览

| 项 | 变更 |
|----|------|
| 账号体系 | 所有人共用 `users` 表。`role` 为 `admin` / `user`。注册为**图片验证码 + 用户名密码**，无需邮箱验证，注册即登录 |
| 登录 | `POST /api/auth/login` 路径不变，返回的 `user` 增加 `role`、`bio`、`created_at`；失败多次后需验证码 |
| JWT | Claims 增加 `role`。`/api/admin/*` 要求 `role = admin`，普通用户返回 `403` |
| 鉴权约定 | **可选登录**：公开接口若带有效 Token，会额外返回「我的状态」字段（`liked` / `favorited`）；无 Token 或无效 Token 视为游客，**不报 401** |
| 上传 | 路径由 `/api/admin/upload` 改为 **`/api/upload`**，登录用户均可用（用于头像、文章图片） |
| 文章对象 | 列表和详情增加 `like_count`、`comment_count`、`favorite_count`；详情还可能有 `liked`、`favorited`（仅登录时） |
| 初始管理员 | 启动时若 `admin` 用户存在但 `role` 为空，自动设为 `admin` |

### 公共用户对象 `UserPublic`

任何接口里出现的用户信息均使用此结构（**不含**密码哈希、IP）：

```json
{ "id": 1, "username": "admin", "nickname": "博主", "avatar": "/uploads/202611/a.png", "bio": "写代码的人", "role": "admin", "created_at": "2026-11-19T10:00:00+08:00" }
```

---

## 五、账号与认证

### 5.1 获取图片验证码（公开）

**GET** `/api/auth/captcha`

```json
{
  "code": 0, "message": "ok",
  "data": {
    "captcha_id": "Xk3dQ9a2LmPz",
    "image": "data:image/png;base64,iVBORw0KGgo...",
    "expires_in": 300
  }
}
```

| 规则 | 说明 |
|------|------|
| 内容 | 4 位数字+字母，带干扰线，**不区分大小写** |
| 有效期 | 300 秒，**一次性**：校验一次（无论对错）即作废，需重新获取 |
| 限流 | 同一 IP 每分钟最多获取 20 次，超出返回 `429` |

### 5.2 注册（公开）

**POST** `/api/auth/register`

```json
{
  "username": "alice",
  "password": "secret123",
  "nickname": "爱丽丝",
  "captcha_id": "Xk3dQ9a2LmPz",
  "captcha_code": "a7K2"
}
```

| 字段 | 类型 | 必填 | 规则 |
|------|------|------|------|
| `username` | string | 是 | 3–20 位，仅字母、数字、下划线；不区分大小写唯一 |
| `password` | string | 是 | 6–64 位 |
| `nickname` | string | 否 | 1–32 字符，缺省取 `username` |
| `captcha_id` / `captcha_code` | string | 是 | 来自 5.1 |

**响应**：注册成功即登录，结构同 5.3 登录成功，`user.role` 固定为 `user`。

| 错误 | code | message 示例 |
|------|------|------|
| 参数不合法 | 400 | `用户名需为 3-20 位字母数字下划线` |
| 验证码错误或过期 | 400 | `验证码错误或已过期` |
| 用户名已存在 | 409 | `用户名已被注册` |
| 触发限流（同一 IP 每小时最多注册 10 个） | 429 | `操作过于频繁` |

### 5.3 登录（公开）

**POST** `/api/auth/login`

```json
{ "username": "alice", "password": "secret123", "captcha_id": "", "captcha_code": "" }
```

`captcha_id` / `captcha_code` 默认可不传。**同一 IP + 用户名在 10 分钟内连续失败 3 次后**，必须携带验证码。

**成功**

```json
{
  "code": 0, "message": "ok",
  "data": { "token": "eyJhbGciOi...", "expires_in": 86400, "user": { "id": 2, "username": "alice", "nickname": "爱丽丝", "avatar": "", "bio": "", "role": "user", "created_at": "2026-11-25T10:00:00+08:00" } }
}
```

**失败**

| 情况 | HTTP / code | data |
|------|------|------|
| 用户名或密码错误 | 401 / 401 | `{ "captcha_required": false }` |
| 连续失败达到阈值，需要验证码 | 401 / 401 | `{ "captcha_required": true }`（前端据此显示验证码输入框） |
| 已要求验证码但未填或填错 | 400 / 400 | `{ "captcha_required": true }` |
| 账号被禁用 | 403 / 403 | `null`，message 为 `账号已被禁用` |

### 5.4 当前用户（需登录）

**GET** `/api/me` → `data` 为 `UserPublic`。Token 无效返回 `401`；账号被禁用返回 `403`。

**PUT** `/api/me` 修改资料

```json
{ "nickname": "新昵称", "avatar": "/uploads/202611/a.png", "bio": "一句话介绍" }
```

| 字段 | 规则 |
|------|------|
| `nickname` | 1–32 字符 |
| `avatar` | 空字符串，或以 `/uploads/` 开头的本站路径 |
| `bio` | 0–200 字符 |

响应 `data` 为更新后的 `UserPublic`。

**PUT** `/api/me/password` 修改密码

```json
{ "old_password": "secret123", "new_password": "newsecret456" }
```

旧密码错误返回 `400`；成功后原 Token 仍有效至过期。响应 `data: null`。

---

## 六、评论（需登录才能发表）

> 游客可阅读评论，**发表、回复、点赞、编辑、删除均需登录**。被禁用账号不能操作。

### 评论对象 `Comment`

```json
{
  "id": 1,
  "content": "写得真好！",
  "user": { "id": 2, "nickname": "爱丽丝", "avatar": "", "role": "user" },
  "reply_to": null,
  "like_count": 3,
  "liked": false,
  "edited": false,
  "created_at": "2026-11-25T10:00:00+08:00",
  "replies": []
}
```

- `replies` 仅顶层评论有，是同一顶层评论下的所有回复（**两级结构**），按时间升序。
- 回复对象的 `reply_to` 为 `{ "id": 2, "nickname": "爱丽丝" }`，顶层评论为 `null`。
- `liked`：登录时表示我是否点赞，游客恒为 `false`。
- 已被删除用户或隐藏的评论不会出现在公开列表。

### 6.1 评论列表（公开，可选登录）

**GET** `/api/posts/:slug/comments`

| 参数 | 默认 | 说明 |
|------|------|------|
| `page` / `page_size` | 1 / 10（最大 50） | 对**顶层评论**分页 |
| `sort` | `latest` | `latest` 最新在前 / `hot` 按点赞数降序 |

`data`：`{ "list": [Comment], "total": 顶层评论数, "page": 1, "page_size": 10, "all_count": 含回复的评论总数 }`。文章不存在或未发布返回 `404`。

### 6.2 发表评论（需登录）

**POST** `/api/posts/:slug/comments`

```json
{ "content": "写得真好！", "parent_id": 0 }
```

| 字段 | 规则 |
|------|------|
| `content` | 1–1000 字符，纯文本，前端转义显示、不渲染 HTML；纯空白视为空 |
| `parent_id` | `0`/不传为顶层评论；否则必须是**同一文章**下未隐藏的评论。回复回复时服务端自动挂到其顶层评论下，并设置 `reply_to` |

限流：**每个用户每分钟最多 5 条**，超出 `429`。

响应 `data`：新建的 `Comment`（`replies: []`）。回复他人评论时，被回复者收到一条通知（见 8.2）；回复自己的评论不通知。

### 6.3 编辑评论（需登录，仅作者本人）

**PUT** `/api/comments/:id` — 请求 `{ "content": "..." }`，规则同上。成功后 `edited = true`。只能编辑自己的评论，否则 `403`。响应 `data` 为更新后的 `Comment`。

### 6.4 删除评论（需登录，作者本人或管理员）

**DELETE** `/api/comments/:id` — 级联删除其回复及相关点赞、通知。他人评论且非管理员返回 `403`。响应 `data: null`。

### 6.5 评论点赞开关（需登录）

**POST** `/api/comments/:id/like` — 已点赞则取消，未点赞则点赞。不能给自己的评论点赞（`400`）。

```json
{ "code": 0, "message": "ok", "data": { "liked": true, "like_count": 4 } }
```

给他人评论点赞时，对方收到一条通知（同一用户对同一评论只在**首次**点赞时通知）。

---

## 七、文章互动与个人中心（需登录）

### 7.1 文章点赞开关

**POST** `/api/posts/:slug/like` → `data`: `{ "liked": true, "like_count": 12 }`

### 7.2 文章收藏开关

**POST** `/api/posts/:slug/favorite` → `data`: `{ "favorited": true, "favorite_count": 5 }`

> 文章不存在或未发布均返回 `404`。

### 7.3 我的收藏

**GET** `/api/me/favorites`（分页）→ `data` 同 1.1 的文章列表结构（不含 `content`），按收藏时间倒序。已下架（草稿）的文章不显示。

### 7.4 我的评论

**GET** `/api/me/comments`（分页）

```json
{
  "list": [
    { "id": 1, "content": "写得真好！", "like_count": 3, "created_at": "2026-11-25T10:00:00+08:00",
      "post": { "id": 1, "title": "Go 语言入门指南", "slug": "go-getting-started" } }
  ],
  "total": 1, "page": 1, "page_size": 10
}
```

---

## 八、通知（需登录）

### 通知对象 `Notification`

```json
{
  "id": 9,
  "type": "reply",
  "actor": { "id": 3, "nickname": "鲍勃", "avatar": "" },
  "post": { "id": 1, "title": "Go 语言入门指南", "slug": "go-getting-started" },
  "comment_id": 12,
  "content": "谢谢分享，学到了",
  "is_read": false,
  "created_at": "2026-11-25T10:30:00+08:00"
}
```

| `type` | 触发条件 | `content` |
|--------|----------|-----------|
| `reply` | 有人回复了我的评论 | 回复内容（截取前 100 字） |
| `comment_like` | 有人给我的评论点赞（首次） | 被点赞的评论内容（截取前 100 字） |

### 8.1 通知列表

**GET** `/api/me/notifications`（分页）— 参数 `unread=1` 仅看未读。按时间倒序。

### 8.2 未读数量

**GET** `/api/me/notifications/unread-count` → `data`: `{ "count": 3 }`（前端导航栏红点轮询，间隔不低于 30 秒）

### 8.3 标记已读

**POST** `/api/me/notifications/read`

```json
{ "ids": [9, 10] }
```

或 `{ "all": true }` 全部已读。响应 `data: null`。

---

## 九、后台管理新增（需 JWT，role = admin）

### 9.1 统计面板

**GET** `/api/admin/stats`

```json
{
  "code": 0, "message": "ok",
  "data": {
    "posts": 12, "published_posts": 10, "draft_posts": 2,
    "comments": 48, "users": 15, "total_views": 3210,
    "recent_users": 4,
    "top_posts": [ { "id": 1, "title": "Go 语言入门指南", "slug": "go-getting-started", "view_count": 520 } ],
    "daily": [ { "date": "2026-11-19", "comments": 3, "users": 1 } ]
  }
}
```

`recent_users` 为近 7 天新注册数；`top_posts` 为浏览量前 5；`daily` 为近 7 天（含今天，升序）每日新增评论与新增用户。

### 9.2 用户管理

**GET** `/api/admin/users`（分页）— 参数 `keyword`（匹配用户名/昵称）、`role`、`status`（`active` / `banned`）

```json
{
  "list": [
    { "id": 2, "username": "alice", "nickname": "爱丽丝", "avatar": "", "bio": "", "role": "user",
      "status": "active", "comment_count": 5, "created_at": "2026-11-25T10:00:00+08:00", "last_login_at": "2026-11-25T11:00:00+08:00" }
  ],
  "total": 1, "page": 1, "page_size": 10
}
```

**PATCH** `/api/admin/users/:id/status` — `{ "status": "banned" }`（取值 `active` / `banned`）。不能禁用管理员或自己（`400`）。被禁用用户已发出的评论保留，但其 Token 在下次请求时返回 `403`。

**DELETE** `/api/admin/users/:id` — 删除用户并级联删除其评论、点赞、收藏、通知。不能删除管理员或自己（`400`）。

### 9.3 评论管理

**GET** `/api/admin/comments`（分页）— 参数 `status`（`approved` / `hidden`）、`post_id`、`user_id`、`keyword`（匹配内容）

```json
{
  "list": [
    { "id": 1, "post": { "id": 1, "title": "Go 语言入门指南", "slug": "go-getting-started" }, "parent_id": 0,
      "user": { "id": 2, "username": "alice", "nickname": "爱丽丝" }, "content": "写得真好！",
      "status": "approved", "like_count": 3, "ip": "127.0.0.1", "created_at": "2026-11-25T10:00:00+08:00" }
  ],
  "total": 1, "page": 1, "page_size": 10
}
```

**PATCH** `/api/admin/comments/:id/status` — `{ "status": "hidden" }`。隐藏顶层评论时其回复一并不展示。
**DELETE** `/api/admin/comments/:id` — 与 6.4 效果相同（管理员可删任意评论）。

---

## 十、图片上传与关于页

### 10.1 上传图片（需登录，变更：`/api/upload`）

**POST** `/api/upload`，`multipart/form-data`，字段 `file`。

| 约束 | 值 |
|------|------|
| 类型 | `image/jpeg`、`image/png`、`image/gif`、`image/webp`，按**文件内容**检测 |
| 大小 | 5 MB 以内；超出 `413` |
| 限流 | 普通用户每小时最多 20 张；管理员不限 |
| 存储 | `backend/uploads/YYYYMM/<随机名>.<ext>` |

```json
{ "code": 0, "message": "ok", "data": { "url": "/uploads/202611/3f9a1c2e5b7d.png", "size": 20480, "content_type": "image/png" } }
```

静态访问：**GET** `/uploads/*filepath`（无 `/api` 前缀，公开）。前端 Vite 需代理 `/uploads`。

### 10.2 获取关于页（公开）

**GET** `/api/about` → `data`: `{ "title": "关于我", "content": "# 你好...", "updated_at": "..." }`。首次启动自动写入默认内容，始终有数据。

### 10.3 更新关于页（需 role = admin）

**PUT** `/api/admin/about` — `{ "title": "关于我", "content": "# ..." }`，`title` 1–64 字符，`content` 为 Markdown，可空。响应 `data: null`。

---

## 十一、前端界面与主题（纯前端，无接口）

参考 **Vue 官网（vuejs.org）** 的视觉语言，要求：

| 项 | 要求 |
|----|------|
| 主色 | 品牌绿 `#42b883`（暗色下用更亮的 `#4fc08d`），深绿辅助色 `#35495e` |
| 布局 | 顶部固定导航栏（Logo、导航、搜索、通知铃铛、主题切换、用户菜单）；文章页采用 **左目录/右内容** 的文档式布局；留白克制、圆角 8px、细边框、淡阴影 |
| 字体 | 正文无衬线（系统字体栈 + Inter），代码用等宽字体；正文行高约 1.75 |
| 暗黑模式 | 三态：**跟随系统 / 浅色 / 深色**；首次默认跟随系统；偏好存 `localStorage`（key `theme`）；切换无闪烁（`<head>` 内联脚本在渲染前设置 `html.dark`）；Markdown、代码块、编辑器、表格、评论区、后台全部适配 |
| 代码块 | 浅色与深色各一套高亮主题，保留语言标签和复制按钮 |
| 响应式 | 桌面 / 平板 / 手机均可用；手机端导航折叠为抽屉 |
| 无障碍 | 对比度满足 WCAG AA；可键盘操作；尊重 `prefers-reduced-motion` |

---

## 十二、新增数据模型

### User（扩展 `users`）

| 字段 | 类型 | 说明 |
|------|------|------|
| `role` | string(16) | `admin` / `user`，默认 `user` |
| `bio` | string(200) | 个人简介 |
| `status` | string(16) | `active` / `banned`，默认 `active` |
| `last_login_at` | datetime（可空） | 最近登录时间 |

`username` 唯一索引（比较时统一转小写）。

### Comment（表 `comments`）

`id`、`post_id`（索引）、`user_id`（索引）、`parent_id`（`0` 为顶层）、`root_id`（顶层评论 ID，顶层自身为 `0`）、`reply_to_user_id`、`content`(text)、`status`（`approved`/`hidden`，默认 `approved`）、`like_count`、`edited`(bool)、`ip`、`created_at`、`updated_at`

### 关联表

| 表 | 字段 | 唯一约束 |
|----|------|----------|
| `post_likes` | `post_id`、`user_id`、`created_at` | (`post_id`,`user_id`) |
| `favorites` | `post_id`、`user_id`、`created_at` | (`post_id`,`user_id`) |
| `comment_likes` | `comment_id`、`user_id`、`created_at` | (`comment_id`,`user_id`) |
| `notifications` | `id`、`user_id`（接收者）、`actor_id`、`type`、`post_id`、`comment_id`、`content`、`is_read`、`created_at` | — |
| `pages` | `id`、`slug`（唯一，目前仅 `about`）、`title`、`content`(longtext)、`updated_at` | `slug` |

> 验证码存放于**服务端内存**（带过期），不入库；重启后未使用的验证码失效，属预期。

---

## 十三、错误码补充

| code | 含义 |
|------|------|
| 409 | 资源冲突（如用户名已被注册） |
| 413 | 上传文件过大 |
| 429 | 请求过于频繁 |
