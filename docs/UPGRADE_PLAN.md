# Vilivala 升级改造方案（第一阶段：只读调研，待确认）

> 调研范围：`other/` 下 4 个参考项目。所有参考目录只读，实施时一律原生重写，不 import 参考目录。

## 0. 结论先行

**4 个参考项目整体弱于现有主站**，主站已具备其大部分能力（三态暗黑模式、左侧 TOC + 滚动高亮、复制代码、字数/阅读时间、上一篇/下一篇、归档、置顶、三级分类）。真正值得借鉴的集中在：**myblog 的工程化与安全做法**、**Bilibala 的功能清单与中文字数算法**、**aircloud 的搜索浮层 UX**。VusaIndex 不是博客，基本无可借鉴。

## 1. 技术栈与项目画像

| 项目 | 框架 | 样式 | 内容管理 | 构建/部署 | 可复用度 |
|---|---|---|---|---|---|
| **主站 vilivala** | Vue 3.5 + TS + Pinia + Router | Tailwind 3 + typography | Go(Gin/GORM)+MySQL，marked + highlight.js | Vite；docker-compose 仅含 MySQL | — |
| **myblog** | Vue 3.5 + TS + Vite 7 | 手写 CSS（`frontend/src/style.css` ~1700 行，`--blog-c-*` 变量）；后台 Element Plus | Go(Gin/GORM)+PostgreSQL；markdown-it + DOMPurify | Docker Compose 三服务 + nginx（`deploy/nginx.conf`） | **高**：栈几乎一致，思路可直接移植 |
| **Bilibala** | Hexo 主题（EJS） | Stylus（`source-src/css/`） | Hexo Markdown | rollup（`rollup.config.js`），CDN 依赖多 | 低：仅功能清单与算法 |
| **bookmagisk-1** | Hexo 主题 aircloud（EJS） | LESS（`source/_less/`） | Hexo Markdown | 无构建 | 低：仅交互思路 |
| **VusaIndex** | Vue 3 + Vite 2 + antd-vue | scoped CSS | Node 后端，无内容模型 | Dockerfile（单阶段） | 几乎无（终端式起始页，非博客） |

**可直接复用**：myblog 的后端分层/限流/Compose/nginx、`useTheme` 思路、TOC 实现。
**需跨栈改写**：Hexo/EJS/jQuery 的全部模板与脚本（仅取逻辑）。

## 2. 参考项目亮点（带源码路径）

### 视觉与阅读
- **myblog** `外观需求文档.md`：VitePress 风格设计规范（克制动效、细边框、内容优先）；`frontend/src/style.css` 的 `--blog-c-*` 语义化变量；`style.css:1339` `prefers-reduced-motion` 降级；`@supports` 包裹的 `backdrop-filter` 导航栏。
- **myblog** `frontend/src/composables/useTheme.ts`：同时更新 `meta[name=theme-color]` 与 `colorScheme`（主站尚未做）。
- **myblog** `frontend/src/utils/markdown.ts`：标题 slug 支持 Unicode 且去重（`-2` 后缀），主站目前用 `h-0/h-1` 序号 id，锚点不稳定。
- **Bilibala** `_darkmode.styl`、`float-btns.ejs`：浮动按钮组（暗黑/回顶）布局；`apple.styl`：macOS 风格代码窗。
- **aircloud** `layout/_partial/toc.ejs` + `source/js/index.js`：滚动联动 TOC（主站已有，无需移植）。
- 不存在：阅读进度条（4 个项目均没有真正的阅读进度条，Bilibala 的 pace.js 是加载条）。

### 核心功能
- **Bilibala** `scripts/helpers/wordcount.js`：中英文分别计数（中文 300 字/分、英文 160 词/分）。主站 `CountWords/ReadingMinutes` 可对照校验。
- **myblog** `backend/internal/middleware/comment_rate_limit.go`：内存 IP 固定窗口限流（含可注入时钟便于测试）；`router.go` + 评论审核状态（approved/待审）。
- **aircloud** `source/js/index.js` L38-220：全屏搜索浮层，Esc/Enter 键盘处理，无结果态。
- **Bilibala** `_config.yml` + `scripts/events/lib/merge-configs.js`：功能开关集中配置（主站 `site.config.ts` 同思路）。
- **Bilibala** 友链页、`topping.ejs` 置顶、`nav.ejs` 上下篇：主站均已实现。
- **myblog** 首页 CMS（`model/homepage_config.go`、`cms_validation.go`）：草稿/发布/预览，对个人博客过度设计，仅借鉴"严格校验、不允许任意 HTML"的思路。
- 不存在：相关文章、RSS、Sitemap（4 个项目均无）。

### 交互细节
- **Bilibala** `source-src/js/bili.js` ~L93-116：回到顶部（阈值淡入）；`post/clipboard.ejs` 复制按钮（主站已有）；`viewer.ejs` PhotoSwipe 图片灯箱。
- **myblog** `LoadingState.vue / ErrorState.vue / PaginationControls.vue`：统一加载/错误态组件。
- 不存在：页面过渡、全局快捷键（均无）。

### SEO 与性能
- **myblog** `ArticleDetailView.vue`：前端更新 `document.title` 与 `meta description`（仅客户端，爬虫不可见，**反面教材**）；Vite 路由级分包；`deploy/nginx.conf`：SPA fallback、`server_tokens off`、`X-Robots-Tag`。
- 其余三者无 OG/Sitemap/JSON-LD，且 Bilibala 大量同步 CDN 脚本，**不要照搬**。

### 架构与可维护性
- **myblog** `backend/internal/{config,database,handler,middleware,model,router,service}` 分层 + 各层 Go 测试 + 前端 vitest（`*.test.ts` 与组件同目录）；`docker-compose.yml` healthcheck + `depends_on: service_healthy`；`测试需求.md` 验收清单。
- 反面：myblog 无 eslint/CI，`frontend/dist` 入库；Bilibala 无 lint/schema/test。

## 3. 主站短板诊断与升级清单

主站当前缺口：无 SEO（OG/canonical/sitemap/RSS）、搜索仅 `LIKE`、`v-html` 渲染未做净化、无登录限流、无测试/lint/CI、无完整容器化部署，以及上轮 E2E 遗留问题。

### P0（高收益、低风险）

| # | 改进点 | 借鉴来源 | 预期收益 | 改动范围与风险 | 新增依赖 |
|---|---|---|---|---|---|
| P0-1 | **Markdown 渲染净化（XSS）**：对 `marked` 输出做白名单净化 | myblog `utils/markdown.ts`（DOMPurify） | 消除 `v-html` 注入风险（导入内容含原始 `<br>/<img>`） | 仅 `utils/markdown.ts`；需放行 hljs 的 class 与复制按钮 | `dompurify`（小） |
| P0-2 | **修复上轮 E2E 遗留**：后台 401 自动跳登录；归档排序不应让置顶干扰时间线；`mailto:` 链接；正文重复 `# 标题`；分类页标题统一；md-editor 工具栏溢出/暗黑预览空白 | 自测报告 | 稳定性与观感 | 前端少量文件 + 归档接口 1 个排序参数（需先改 `API.md`）；低 | 无 |
| P0-3 | **稳定标题锚点**：Unicode slug + 去重，替代 `h-0` 序号 | myblog `markdown.ts` | 可分享的章节链接 | `utils/markdown.ts`；低 | 无 |
| P0-4 | **回到顶部 + 阅读进度条**（Tailwind 原生，尊重 `prefers-reduced-motion`） | Bilibala `bili.js`、`float-btns.ejs`（仅回顶；进度条为自研） | 长文阅读体验 | 新增 2 个组件；低 | 无 |
| P0-5 | **主题细节**：同步 `theme-color` meta 与 `colorScheme` | myblog `useTheme.ts` | 移动端地址栏/滚动条协调 | `stores/theme.ts`；极低 | 无 |
| P0-6 | **登录接口限流** | myblog `comment_rate_limit.go` | 防暴力破解 | 新增 middleware 约 60 行；低 | 无 |
| P0-7 | **基础 SEO 元信息**：`index.html` 补 OG/Twitter/canonical；文章页前端更新 meta | myblog（改进其仅客户端的缺陷） | 分享卡片可用 | 低；爬虫仍看不到 SPA 内容（见 P1-1） | 无 |

### P1（核心功能/体验升级）

| # | 改进点 | 借鉴来源 | 预期收益 | 改动范围与风险 | 新增依赖 |
|---|---|---|---|---|---|
| P1-1 | **服务端 SEO**：Go 提供 `/sitemap.xml`、`/rss.xml`、`/robots.txt`，并对文章路由注入 meta + JSON-LD（爬虫 UA 或统一注入到 `index.html` 模板） | 三者均缺，自研；nginx 思路参考 myblog `deploy/nginx.conf` | RSS 订阅 + 搜索引擎收录 | 后端新增 handler 与 API.md 条目；中（SPA 注入方案需选型） | 无 |
| P1-2 | **全文搜索升级**：MySQL `FULLTEXT (ngram)` + 高亮片段 + 搜索浮层（Esc/↑↓/Enter）与 `Ctrl/⌘+K` | aircloud `index.js` 搜索浮层 UX | 255 篇内容下的检索质量与速度 | 后端索引与 `/api/search`、前端新组件；中 | 无 |
| P1-3 | **相关文章**（同分类/同标签评分） | 参考项目均无，自研 | 提升停留 | 后端新增 `related` 字段/接口；低-中 | 无 |
| P1-4 | **图片灯箱 + 图片懒加载 + 外链图片本地化**（迁移时 CSDN 图片被拦截） | Bilibala `viewer.ejs`（PhotoSwipe） | 图文阅读、修复图片失效 | 渲染层改造 + 导入脚本下载图片；中 | `photoswipe`（或自研轻量实现） |
| P1-5 | **工程化**：ESLint + Prettier、vitest（markdown/format/theme 单测）、Go 各层测试、GitHub Actions（type-check/build/test） | myblog 测试布局 | 重构安全网 | 新增配置与测试文件；低 | `eslint`、`prettier`、`vitest`（dev） |
| P1-6 | **容器化部署**：前端 nginx 多阶段镜像 + 后端镜像 + compose healthcheck；密钥改环境变量、提供 `.env.example` | myblog `docker-compose.yml`、`deploy/nginx.conf` | 一键上线 | 新增 Dockerfile/nginx.conf；中 | 无 |
| P1-7 | **统一加载/错误/空态组件**与路由级骨架屏 | myblog `LoadingState/ErrorState.vue` | 一致体验 | 抽组件替换现有零散写法；低 | 无 |

### P2（重度重构/可选增强）

| # | 改进点 | 借鉴来源 | 说明 | 新增依赖 |
|---|---|---|---|---|
| P2-1 | 评论系统 + 账号体系（即暂缓的 v1.2：验证码注册、二级回复、审核、限流） | myblog 评论审核状态、限流 | 工作量大，API.md 已有规划 | 图片验证码库（Go） |
| P2-2 | 页面切换过渡（View Transitions API，带降级） | 参考项目均无 | 锦上添花 | 无 |
| P2-3 | 预渲染/SSR（vite-ssg 或 Go 侧渲染）以彻底解决 SPA SEO | 自研 | 改动大，与 P1-1 二选一深化 | 视方案 |
| P2-4 | 首页模块可配置（Hero/精选/最新开关与排序），仅取校验思路 | myblog 首页 CMS | 对个人站过度设计 | 无 |
| P2-5 | 代码块增强：行号、macOS 窗口风格、语言标签 | Bilibala `apple.styl` | 观感 | 无 |
| P2-6 | 全局快捷键帮助面板、阅读设置（字号/宽度） | 自研（VusaIndex 命令面板仅作灵感） | 可选 | 无 |

### 明确不借鉴
jQuery/EJS/Stylus/LESS 代码、CDN 同步脚本（Bilibala 的 pace/sweetalert）、VusaIndex 全部、myblog 的 Element Plus 与首页 CMS 全量实现、任何提交 `dist` 的做法。

## 4. 约束承诺
- 第一阶段未修改任何业务代码，本文件为唯一新增文件。
- 实施时凡涉及接口字段变更，先更新 `docs/API.md`。
- 所有移植均原生重写；完成后删除 `other/` 主站仍可独立构建（可在 P1-5 CI 中加入校验）。
- 沿用 Tailwind，不引入第二套 CSS 框架。

## 5. 待你确认
1. 执行范围：仅 P0？P0 + P1？还是逐条点名？
2. 是否保留当前视觉基调（Vue 官网绿 `#42b883` + 深色网格 hero）？我的建议：保留，仅增量优化。
3. P1-1 的 SEO 方案取向：轻量元信息注入（建议）还是 P2-3 预渲染？
4. P1-4 外链图片：下载到本地 `/uploads`（建议），还是保留外链？
