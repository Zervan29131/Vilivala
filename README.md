Vue+Go 博客系统

gin框架


结合你提供的Vue+Gin博客系统后端代码，完整开发流程可分为**技术栈确认、项目架构设计、分端开发、联调测试、部署上线**5大阶段，具体步骤如下：

### 一、明确技术栈与核心依赖
#### 1. 前端（Vue）核心技术
- 基础框架：Vue 3（推荐，搭配Composition API）
- 路由：Vue Router（管理页面跳转）
- 状态管理：Pinia（替代Vuex，轻量简洁）
- 网络请求：Axios（对接Gin后端API）
- UI组件库：Element Plus（快速搭建后台界面）
- 富文本编辑器：TinyMCE/CKEditor（用于文章编辑）

#### 2. 后端（Gin）核心技术
- Web框架：Gin（高性能HTTP框架，你提供的代码已基于此）
- ORM：GORM（数据库操作，代码中已集成）
- 认证：JWT（用户登录态管理，middleware/jwt.go已实现）
- 数据库：MySQL（存储业务数据）
- 工具依赖：bcrypt（密码加密，util/crypto.go）、Viper（配置读取，config/config.go）

### 二、项目整体架构设计
采用**前后端分离架构**，职责清晰：
- 前端：负责页面渲染、用户交互、请求发送与响应处理
- 后端：提供RESTful API、数据校验、业务逻辑处理、数据库操作、权限控制
- 通信方式：JSON格式数据交互，JWT令牌传递登录状态

### 三、分阶段开发实施
#### 阶段1：需求梳理与数据库设计（核心前置）
1. **明确核心功能**
   - 基础功能：用户注册/登录、文章CRUD、分类/标签管理、评论互动
   - 扩展功能：文章搜索、阅读量统计、头像/封面图上传（util/file.go可扩展）、管理员权限（AdminMiddleware已支持）
2. **数据库设计（对应model目录）**
   - 复用你代码中的4张核心表：`users`（用户）、`articles`（文章）、`categories`（分类）、`tags`（标签）
   - 补充关联表：`article_tags`（文章-标签多对多关联）、`comments`（评论，若需）
   - 执行数据库迁移：用GORM的`AutoMigrate`自动创建表（可在main.go中初始化）

#### 阶段2：后端开发（基于你提供的代码扩展）
1. **环境搭建**
   - 初始化Go项目：`go mod init my-blog-backend`
   - 安装依赖：`go get github.com/gin-gonic/gin gorm.io/gorm gorm.io/driver/mysql`等
   - 配置`config.yaml`：填写实际MySQL密码、JWT密钥
2. **完善核心模块**
   - 初始化入口：编写`main.go`，依次初始化配置（config.InitConfig）、数据库（model.InitDB）、路由（router.InitRouter）
   - 补充缺失代码：如`util/file.go`文件上传逻辑、`controller`中未完成的删除/发布接口、`service`层权限校验完善
   - 接口测试：用Postman测试所有API（如注册`/api/v1/user/register`、登录`/api/v1/user/login`）
3. **关键功能实现要点**
   - 认证：登录成功返回JWT令牌，前端存储后，每次请求在Header携带`Authorization: Bearer 令牌`
   - 权限：普通用户仅操作自己的文章，管理员可管理所有用户/文章（依赖AdminMiddleware）
   - 文件上传：用Gin接收文件，存储到`uploads`目录，返回文件URL（配合`r.Static("/uploads", "./uploads")`访问）

#### 阶段3：前端开发（Vue）
1. **项目初始化**
   - 创建Vue项目：`npm create vue@latest my-blog-frontend`，选择Router、Pinia等依赖
   - 安装插件：`npm install axios element-plus @tinymce/tinymce-vue`
2. **核心模块开发**
   - 全局配置：Axios封装（设置BaseURL、请求拦截器添加JWT令牌、响应拦截器处理错误）
   - 页面拆分：
     - 公开页面：首页（文章列表）、文章详情页、登录/注册页
     - 私有页面：个人中心、文章编辑页、分类管理页
     - 管理员页面：用户管理页、全量文章管理页
   - 状态管理：用Pinia存储用户信息、登录状态（避免重复请求）
   - 交互实现：调用后端API，如登录调用`/api/v1/user/login`，文章列表调用`/api/v1/article/list`

#### 阶段4：前后端联调
1. **跨域解决**：后端已通过`CorsMiddleware`支持跨域，前端直接请求后端地址即可
2. **接口联调**：逐一验证功能，如：
   - 登录：前端提交用户名密码，获取令牌后存储到localStorage
   - 文章发布：携带令牌调用`/api/v1/article`接口，传递标题/内容等参数
   - 权限控制：测试普通用户访问管理员接口（如`/api/v1/user/list`）是否返回403
3. **问题排查**：用浏览器F12查看网络请求，后端打印日志（Gin默认Logger）定位错误

#### 阶段5：部署上线
1. **前端部署**
   - 打包：`npm run build`生成dist目录
   - 部署：上传到Nginx根目录，或用Vercel、Netlify等平台
2. **后端部署**
   - 编译：`GOOS=linux GOARCH=amd64 go build -o my-blog`生成Linux可执行文件
   - 部署：上传到云服务器（如阿里云ECS），配合MySQL服务启动，用Systemd管理进程
3. **域名与HTTPS**：配置域名解析到服务器，通过Let's Encrypt申请免费HTTPS证书

### 四、参考资源与优化方向
1. **核心参考**：你提供的后端代码已覆盖配置、模型、中间件、路由等核心模块，可直接复用并补充接口
2. **优化建议**：
   - 前端：添加加载动画、表单验证、错误提示提升体验
   - 后端：增加接口限流（如Gin Rate Limit）、日志切割、数据库索引优化
   - 安全：敏感配置用环境变量（Viper支持），避免硬编码密码