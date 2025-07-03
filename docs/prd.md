# MDBlog - Markdown 博客系统 PRD

## 🎯 项目概述

MDlog 是一个基于 Go 语言开发的现代化 Markdown 博客系统，专注于简单、高效、易部署的博客解决方案。系统支持实时搜索、分类标签、响应式设计，并实现了真正的单文件部署。

### 核心优势

- **零依赖**: 无需数据库，纯文件系统驱动
- **高性能**: 内存缓存，毫秒级响应
- **易部署**: 单文件部署，一键自动化
- **现代化**: 响应式设计，移动端优化
- **开发友好**: 热重载，Markdown 驱动

## ✅ 已实现功能

### 1. Markdown 文档处理系统

**功能描述**: 完整的 Markdown 文档解析和渲染系统

**技术实现**:

- 使用 Goldmark 解析器，支持 CommonMark 规范
- 集成 Chroma 语法高亮，支持 100+ 编程语言
- 支持 GFM 扩展（表格、任务列表、删除线等）
- 自定义 YAML Front Matter 解析
- HTML 安全渲染，防止 XSS 攻击

**支持格式**:

```yaml
---
title: "文章标题"
author: "作者名称"
github_url: "https://github.com/username"
created_at: 2024-01-01T10:00:00
updated_at: 2024-01-01T12:00:00
description: "文章描述"
tags: ["Go", "Web开发", "技术"]
---
```

### 2. 高性能内存数据管理

**功能描述**: 全内存数据存储和管理系统

**技术特性**:

- 启动时一次性加载所有文章到内存
- 使用 `sync.RWMutex` 保证并发安全
- 多维度索引：分类索引、标签索引、搜索索引
- 智能摘要生成（前 3 行文本内容）
- 实时数据更新机制（为 Webhook 预留）

**性能指标**:

- 文章加载：< 100ms（1000 篇文章）
- 搜索响应：< 10ms
- 内存占用：约 1MB/100 篇文章

### 3. 实时搜索系统

**功能描述**: 快捷键驱动的实时搜索功能

**交互设计**:

- 快捷键：双击 `Cmd`（macOS）/ `Ctrl`（Windows/Linux）
- 搜索范围：文章标题、内容、标签
- 搜索算法：模糊匹配，大小写不敏感
- 实时结果：输入即搜索，无需点击

**技术实现**:

- 前端：原生 JavaScript，防抖优化
- 后端：内存搜索，字符串包含匹配
- UI：弹窗式搜索界面，键盘导航

### 4. 响应式 UI 设计

**功能描述**: 多端适配的现代化用户界面

**设计规范**:

- **桌面端（>768px）**: 侧边导航，多列布局
- **平板端（768px-480px）**: 折叠导航，双列布局
- **移动端（<480px）**: 汉堡菜单，单列布局

**技术实现**:

- CSS3 Flexbox/Grid 布局
- 媒体查询响应式设计
- 触摸友好的交互设计
- SVG 矢量图标

### 5. 分类标签系统

**功能描述**: 基于目录结构的自动分类和多标签支持

**分类规则**:

- 自动识别：`posts/go/article.md` → 分类 "go"
- 根目录文章：默认分类 "其他"
- 特殊处理：`about.md` 仅在关于页面显示

**标签功能**:

- Front Matter 定义：`tags: ["Go", "Web开发"]`
- 标签页面：`/tag/Go`
- 多标签过滤支持

### 6. 静态资源嵌入系统

**功能描述**: 所有静态资源嵌入到可执行文件

**技术实现**:

```go
//go:embed web
var EmbeddedAssets embed.FS
```

**优势**:

- 真正的单文件部署
- 无需 web 目录依赖
- 版本一致性保证
- 部署简化

### 7. 自动化部署系统

**功能描述**: 一键编译、部署、监控的完整 DevOps 流程

**部署流程**:

```bash
make scp  # 一键部署
```

**执行步骤**:

1. 交叉编译 Linux 可执行文件
2. SSH 连接服务器停止服务
3. SCP 上传新的可执行文件
4. SSH 启动服务
5. 健康检查（端口 8091）
6. 显示服务状态和日志

**配置参数**:

```makefile
SERVER_HOST=123.56.186.148    # 服务器IP
SSH_PORT=21345               # SSH端口
APP_PORT=8091                # 应用端口
```

### 8. 服务监控系统

**功能描述**: Systemd 集成的服务管理和日志监控

**服务配置**:

- Systemd 服务文件：`deploy/mdblog.service`
- 日志输出：`/bro/mdblog/logs/mdblog.log`
- 自动重启：服务异常时自动恢复
- 开机自启：系统重启后自动启动

## 🚧 开发中功能

### Webhook 自动同步（✅ 95% 完成）

**目标**: 实现 Git 仓库变更时自动更新博客内容

**已完成**:

- ✅ HTTP 接口：`POST /webhook/gitee`
- 🟡 签名验证（已实现多种算法，开发模式可跳过）
- ✅ Git 操作：`git pull` 执行
- ✅ 内存数据重新加载
- ✅ 错误处理和日志记录
- ✅ 健康检查接口：`GET /health`
- ✅ Webhook 测试工具
- ✅ 开发模式配置

**已知问题**:

- ⚠️ Gitee 签名验证算法与标准不同，目前采用开发模式跳过验证
- 📋 生产环境建议使用 IP 白名单 + HTTPS 替代签名验证

**技术实现**:

```go
// Webhook 配置
const (
    WebhookSecret = "c85b7544d37c89280afcf912ec70f4083b18065d9e89966cae6059798f0dadf5"
    WebhookBranch = "main"     // 监听的分支
)

// posts 目录路径动态确定：
// 1. 命令行参数：./mdblog -posts /path/to/posts
// 2. 默认位置：可执行文件同级的 posts 目录

// Gitee Webhook 载荷结构
type GiteeWebhookPayload struct {
    Ref        string `json:"ref"`
    Repository struct {
        Name        string `json:"name"`
        FullName    string `json:"full_name"`
        CloneURL    string `json:"clone_url"`
        SSHURL      string `json:"ssh_url"`
        GitHTTPURL  string `json:"git_http_url"`
        GitSSHURL   string `json:"git_ssh_url"`
    } `json:"repository"`
    Commits []struct {
        ID        string    `json:"id"`
        Message   string    `json:"message"`
        Timestamp time.Time `json:"timestamp"`
        Author    struct {
            Name  string `json:"name"`
            Email string `json:"email"`
        } `json:"author"`
        Added    []string `json:"added"`
        Removed  []string `json:"removed"`
        Modified []string `json:"modified"`
    } `json:"commits"`
    HeadCommit struct {
        ID        string    `json:"id"`
        Message   string    `json:"message"`
        Timestamp time.Time `json:"timestamp"`
        Author    struct {
            Name  string `json:"name"`
            Email string `json:"email"`
        } `json:"author"`
    } `json:"head_commit"`
    Pusher struct {
        Name  string `json:"name"`
        Email string `json:"email"`
    } `json:"pusher"`
}
```

**安全机制**:

1. **HMAC-SHA256 签名验证**: 使用密钥验证请求来源的合法性
2. **分支过滤**: 只处理来自 `main` 分支的推送事件
3. **输入验证**: 验证请求方法、载荷格式和必需头部
4. **错误处理**: 全面的错误捕获和日志记录

**工作流程**:

1. 接收 Gitee Webhook POST 请求
2. 验证 `X-Gitee-Token` 签名
3. 解析 JSON 载荷并检查目标分支
4. 执行 `git pull origin main` 更新 posts 仓库
5. 清空内存缓存并重新加载所有文章
6. 返回同步结果和统计信息

**使用说明**:

```bash
# 测试本地 Webhook
make webhook-test-local

# 测试生产环境 Webhook
make webhook-test-remote

# 健康检查
curl http://localhost:8091/health

# 查看 Webhook 日志
tail -f /bro/mdblog/logs/mdblog.log | grep -i webhook
```

**Gitee 配置**:

- **Webhook URL**: `https://broqiang.com/webhook/gitee`
- **密码**: `c85b7544d37c89280afcf912ec70f4083b18065d9e89966cae6059798f0dadf5`
- **触发事件**: Push
- **分支过滤**: main

**性能优化**:

- 使用读写锁保证并发安全
- 完整重载而非增量更新（确保数据一致性）
- 异步处理，不阻塞 Webhook 响应
- Git 操作超时控制和错误恢复

## 📋 计划功能

### 1. 主题系统（v2.0）

**功能描述**: 支持亮色/暗色主题切换

**设计方案**:

- CSS 变量主题系统
- 浏览器本地存储偏好
- 一键切换按钮
- 默认跟随系统主题

### 2. 评论系统（v2.1）

**功能描述**: 集成第三方评论系统

**技术选型**:

- Disqus 集成
- Gitalk 集成（基于 GitHub Issues）
- 自定义评论系统（可选）

### 3. RSS 订阅（v2.2）

**功能描述**: 自动生成 RSS Feed

**实现要点**:

- XML 格式 RSS 2.0
- 自动更新机制
- 分类订阅支持

### 4. 文章归档（v2.3）

**功能描述**: 按时间归档文章列表

**设计方案**:

- 年度归档页面
- 月份归档页面
- 时间轴展示

### 5. 站点地图（v2.4）

**功能描述**: SEO 优化的站点地图

**实现要点**:

- XML Sitemap 生成
- 搜索引擎提交
- 自动更新机制

## 🏗️ 技术架构

### 后端技术栈

| 组件          | 技术选型      | 版本要求 | 用途说明              |
| ------------- | ------------- | -------- | --------------------- |
| 编程语言      | Go            | 1.21+    | 高性能后端开发        |
| Web 框架      | Gin           | v1.9+    | HTTP 路由和中间件     |
| Markdown 解析 | Goldmark      | v1.5+    | Markdown 到 HTML 转换 |
| 语法高亮      | Chroma        | v2.0+    | 代码块语法高亮        |
| 模板引擎      | html/template | 内置     | HTML 模板渲染         |
| 静态资源      | embed         | 内置     | 资源嵌入              |

### 前端技术栈

| 组件       | 技术选型 | 说明                |
| ---------- | -------- | ------------------- |
| HTML       | HTML5    | 语义化标签          |
| CSS        | CSS3     | Flexbox/Grid 布局   |
| JavaScript | ES6+     | 原生 JS，无框架依赖 |
| 图标       | SVG      | 矢量图标，内嵌      |
| 字体       | Inter    | 现代化字体          |

### 数据流架构

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Markdown      │    │   Data Manager  │    │   Web Server    │
│   Files         │───▶│   (Memory)      │───▶│   (Gin Router)  │
│   (/posts/)     │    │                 │    │                 │
└─────────────────┘    └─────────────────┘    └─────────────────┘
                              │                          │
                              ▼                          ▼
                       ┌─────────────────┐    ┌─────────────────┐
                       │   Search Index  │    │   HTTP Response │
                       │   Tag Index     │    │   (HTML/JSON)   │
                       │   Category Index│    │                 │
                       └─────────────────┘    └─────────────────┘
```

### API 接口设计

#### 页面路由

```
GET  /                     # 首页（文章列表）
GET  /post/:id            # 文章详情
GET  /category/:category  # 分类页面
GET  /tag/:tag           # 标签页面
GET  /about              # 关于页面
```

#### API 接口

```
GET  /api/posts          # 文章列表API
GET  /api/posts/:id      # 文章详情API
GET  /api/categories     # 分类列表API
GET  /api/tags          # 标签列表API
GET  /api/search        # 搜索API
```

#### 管理接口

```
POST /webhook/gitee     # Gitee Webhook
GET  /health           # 健康检查
```

### 数据结构设计

#### 核心数据模型

```go
// Post 文章数据结构
type Post struct {
    ID          string    `json:"id"`          // 文章ID（路径）
    Title       string    `json:"title"`       // 标题
    Author      string    `json:"author"`      // 作者
    GitHubURL   string    `json:"github_url"`  // GitHub链接
    Content     string    `json:"content"`     // Markdown内容
    HTML        string    `json:"html"`        // HTML内容
    Summary     string    `json:"summary"`     // 摘要
    Category    string    `json:"category"`    // 分类
    CreateTime  time.Time `json:"created_at"`  // 创建时间
    UpdateTime  time.Time `json:"updated_at"`  // 更新时间
    Description string    `json:"description"` // 描述
    Tags        []string  `json:"tags"`        // 标签
    FilePath    string    `json:"file_path"`   // 文件路径
}

// BlogData 全局数据结构
type BlogData struct {
    Posts       map[string]*Post    // 文章索引
    Categories  map[string][]string // 分类索引
    Tags        map[string][]string // 标签索引
    SearchIndex map[string][]string // 搜索索引
    LastUpdate  time.Time           // 最后更新时间
}
```

## 🚀 性能指标

### 性能目标

| 指标       | 目标值 | 当前值 | 说明              |
| ---------- | ------ | ------ | ----------------- |
| 启动时间   | <3s    | ~2s    | 冷启动到可用      |
| 首页响应   | <50ms  | ~20ms  | 不含网络延迟      |
| 搜索响应   | <100ms | ~30ms  | 1000 篇文章内搜索 |
| 内存占用   | <100MB | ~50MB  | 1000 篇文章       |
| 可执行文件 | <50MB  | ~25MB  | 包含所有资源      |

### 压力测试结果

```bash
# 并发测试（100并发，持续1分钟）
wrk -t10 -c100 -d60s http://localhost:8091/

Results:
  Requests/sec:    2847.32
  Transfer/sec:    892.45KB
  Avg Latency:     35.12ms
  Max Latency:     187.89ms
```

## 🔒 安全考虑

### 已实现安全措施

1. **HTML 安全渲染**: 防止 XSS 攻击
2. **输入验证**: 搜索参数验证
3. **路径遍历保护**: 文件路径安全检查
4. **CORS 配置**: 跨域请求控制

### 计划安全增强

1. **Webhook 签名验证**: HMAC 签名校验
2. **Rate Limiting**: API 请求频率限制
3. **HTTPS 强制**: TLS 加密传输
4. **访问日志**: 安全审计日志

## 📊 监控指标

### 应用监控

- **服务状态**: Systemd 服务监控
- **端口检查**: TCP 8091 端口可用性
- **日志监控**: 应用错误日志
- **内存使用**: 进程内存占用

### 业务指标

- **文章数量**: 加载的文章总数
- **搜索频率**: 搜索 API 调用统计
- **热门文章**: 访问频率统计（计划中）
- **用户行为**: 页面访问统计（计划中）

## 🎯 发展路线图

### V1.0 - 核心功能（已完成）

- ✅ Markdown 博客基础功能
- ✅ 搜索和分页
- ✅ 响应式设计
- ✅ 单文件部署
- ✅ 自动化部署

### V1.1 - 增强功能（开发中）

- 🚧 Webhook 自动同步
- 📋 性能优化
- 📋 错误处理完善

### V2.0 - 用户体验（计划中）

- 📋 主题系统
- 📋 评论系统
- 📋 RSS 订阅
- 📋 文章归档

### V3.0 - 高级功能（远期）

- 📋 多语言支持
- 📋 插件系统
- 📋 管理后台
- 📋 统计分析

## 🤝 贡献指南

### 开发环境搭建

```bash
# 克隆项目
git clone https://github.com/your-username/mdblog.git
cd mdblog

# 安装依赖
go mod tidy

# 启动开发服务器
make dev
```

### 代码规范

- 遵循 Go 官方代码规范
- 使用 `gofmt` 格式化代码
- 编写单元测试
- 提交前运行 `go vet` 检查

### 提交流程

1. Fork 项目
2. 创建功能分支
3. 提交代码
4. 创建 Pull Request
5. 等待 Code Review

## 📄 许可证

MIT License - 详见 [LICENSE](LICENSE) 文件

---

**MDlog PRD** - 详细的产品需求文档，指导项目开发和迭代方向。
