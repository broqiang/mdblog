# MDBlog - Markdown 博客系统 PRD

## 🎯 项目概述

MDlog 是一个基于 Go 语言开发的现代化 Markdown 博客系统，专注于简单、高效、易部署的博客解决方案。系统支持实时搜索、分类导航、响应式设计，并实现了真正的单文件部署。

### 核心优势

- **零依赖**: 无需数据库，纯文件系统驱动
- **高性能**: 内存缓存，毫秒级响应
- **易部署**: 单文件部署，一键自动化
- **现代化**: 响应式设计，移动端优化
- **开发友好**: 热重载，Markdown 驱动

## 🔧 开发过程与技术要点

### 核心技术决策

#### 1. 架构设计原则

**单一职责模块化设计**：

- `internal/config/` - 全局配置常量，避免硬编码
- `internal/data/` - 数据管理和缓存逻辑
- `internal/markdown/` - Markdown 解析专用模块
- `internal/server/` - Web 服务器和路由处理

**内存优先策略**：

- 启动时一次性加载所有文章到内存（适合中小型博客）
- 使用 `sync.RWMutex` 保证并发安全
- 缓存文章内容、分类索引、搜索索引

#### 2. 关键技术选型

**Markdown 处理链**：

```go
// 使用 Goldmark + Chroma 组合
markdown := goldmark.New(
    goldmark.WithExtensions(
        extension.GFM,              // GitHub Flavored Markdown
        extension.Footnote,         // 脚注支持
        highlighting.NewHighlighting(
            highlighting.WithStyle("github"),
            highlighting.WithFormatOptions(),
        ),
    ),
    goldmark.WithParserOptions(
        parser.WithAutoHeadingID(), // 自动生成标题ID
    ),
    goldmark.WithRendererOptions(
        html.WithHardWraps(),       // 硬换行
        html.WithXHTML(),           // XHTML 兼容
    ),
)
```

**前端技术栈**：

- 原生 JavaScript（无框架依赖，减少复杂度）
- CSS3 Flexbox/Grid（响应式布局）
- SVG 图标（矢量，内嵌，性能好）

#### 3. 静态资源嵌入方案

**关键实现**：

```go
//go:embed web
var EmbeddedAssets embed.FS

// 服务静态文件的关键代码
func (s *Server) setupStaticRoutes() {
    // 嵌入的文件系统需要去掉前缀
    webFS, _ := fs.Sub(s.assets, "web")
    s.router.StaticFS("/static", http.FS(webFS))
}
```

**优势**：

- 真正的单文件部署
- 避免静态文件丢失问题
- 版本一致性保证

### 开发过程中的关键问题与解决方案

#### 1. **Gitee Webhook 签名验证实现**

**问题描述**：
Gitee 的 Webhook 签名机制与 GitHub/GitLab 标准类似，但实现细节稍有不同。

**技术分析**：

- GitHub: `X-Hub-Signature-256: sha256=xxxxx`（HMAC-SHA256）
- Gitee: `X-Gitee-Token: xxxxx` + `X-Gitee-Timestamp: xxxxx`（也使用 HMAC-SHA256，但算法稍有不同）

**解决方案**：

```go
// Gitee 签名验证实现
func verifyGiteeSignature(signature, timestamp, secret string) bool {
    timestampInt, err := strconv.ParseInt(timestamp, 10, 64)
    if err != nil {
        return false
    }
    // 构造待签名字符串：timestamp + "\n" + secret
    stringToSign := fmt.Sprintf("%d\n%s", timestampInt, secret)
    mac := hmac.New(sha256.New, []byte(secret))
    mac.Write([]byte(stringToSign))
    signData := mac.Sum(nil)
    encodedSign := base64.StdEncoding.EncodeToString(signData)
    urlEncodedSign := url.PathEscape(encodedSign)
    return urlEncodedSign == signature
}

// 验证逻辑
giteeToken := r.Header.Get("X-Gitee-Token")
giteeTimestamp := r.Header.Get("X-Gitee-Timestamp")
if !verifyGiteeSignature(giteeToken, giteeTimestamp, config.WebhookSecret) {
    return false
}
```

**关键技术点**：

1. **签名算法**：`HMAC-SHA256(timestamp + "\n" + secret)` → Base64 → URL 编码
2. **时间戳验证**：防止重放攻击，设置时间窗口（如 1 小时）
3. **头部字段**：同时需要 `X-Gitee-Token` 和 `X-Gitee-Timestamp`

**生产环境安全策略**：

- 设置复杂的 WebhookSecret（64 位随机字符串）
- 启用 HTTPS 加密传输
- 配置时间戳验证窗口，防重放攻击
- 生产环境使用 Nginx 作为反向代理，所有请求均经由 Nginx 转发到 Go 服务，HTTPS 证书统一配置在 Nginx 层

#### 2. **Git 操作权限问题**

**问题描述**：
服务器上使用 `bro` 用户运行服务，确保 Git 操作具有正确的权限和认证。

**实际部署架构**：

- **服务运行用户**: `bro`（普通用户，非 root）
- **用户家目录**: `/bro`
- **SSH 密钥配置**: 本地电脑与服务器间已配置公钥认证
- **Git 仓库认证**: Gitee SSH 公钥使用相同密钥对

**技术解决方案**：

```go
// Git 操作前检查权限
func (dm *DataManager) executeGitPull() error {
    // 检查目录权限
    if err := dm.checkGitPermissions(); err != nil {
        return fmt.Errorf("权限检查失败: %v", err)
    }

    // 执行 git pull 带超时控制
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()

    cmd := exec.CommandContext(ctx, "git", "pull", "origin", "main")
    cmd.Dir = dm.postsDir

    output, err := cmd.CombinedOutput()
    if err != nil {
        return fmt.Errorf("git pull 失败: %v, 输出: %s", err, output)
    }

    return nil
}
```

**SSH 密钥配置详解**：

1. **本地密钥生成**：

   在本地电脑上使用如下命令生成 SSH 密钥对：

   ```bash
   ssh-keygen -t rsa -b 4096 -C "your-email@example.com"
   ```

   生成后，可通过如下命令查看密钥文件：

   ```bash
   ls ~/.ssh
   # 常见输出：id_rsa  id_rsa.pub  known_hosts
   ```

2. **服务器端配置** (`/bro/.ssh/authorized_keys`):

   ```bash
   # 将本地生成的公钥内容（id_rsa.pub）添加到服务器，允许免密码 SSH 连接
   ssh-rsa AAAAB3NzaC1yc2E... your-local-computer
   ```

3. **Gitee SSH 公钥配置**:

   ```bash
   # 在 Gitee 设置中添加相同的公钥
   # 路径：个人设置 → SSH 公钥 → 添加公钥
   ```

**服务启动脚本准备**：

在进行权限验证前，需要提前在服务器上配置好 Systemd 的 service 启动脚本（如 /etc/systemd/system/mdblog.service），用于以 bro 用户身份启动 mdblog 服务，并实现开机自启、自动重启等功能。

配置完成后，使用如下命令启动和管理服务：

```bash
sudo systemctl daemon-reload
sudo systemctl enable mdblog
sudo systemctl start mdblog
sudo systemctl status mdblog
```

**权限验证步骤**：

```bash
# 1. 验证 SSH 连接
ssh bro@your-ip -p your-port

# 2. 验证 Git SSH 连接
ssh -T git@gitee.com

# 3. 验证 Git 仓库访问
cd /bro/mdblog && git pull origin main

# 4. 检查文件权限
ls -la /bro/mdblog/posts/
```

**最佳实践**：

- **用户一致性**: 服务运行用户 `bro` 与 Git 仓库所有者保持一致
- **SSH 密钥认证**: 使用相同的 SSH 密钥对（本地 ↔ 服务器，服务器 ↔Gitee）
- **文件权限**: 确保 `bro` 用户对 posts 目录有完整的读写权限
- **安全策略**: 使用普通用户运行服务，避免 root 权限风险

**部署权限检查清单**：

- ✅ `bro` 用户可以正常 SSH 登录服务器
- ✅ 服务器可以通过 SSH 连接到 Gitee (`ssh -T git@gitee.com`)
- ✅ `bro` 用户对 `/bro/mdblog/posts/` 目录有读写权限
- ✅ Git 仓库配置正确（origin 指向 Gitee SSH URL）
- ✅ Systemd 服务配置为 `bro` 用户运行

### 并发安全设计

**核心数据结构**：

```go
type DataManager struct {
    mu          sync.RWMutex           // 读写锁
    posts       map[string]*Post       // 文章缓存
    categories  map[string][]string    // 分类索引
    searchIndex map[string][]string    // 搜索索引
}

// 读操作（多个可并发）
func (dm *DataManager) GetPost(id string) *Post {
    dm.mu.RLock()
    defer dm.mu.RUnlock()
    return dm.posts[id]
}

// 写操作（独占锁）
func (dm *DataManager) ReloadData() error {
    dm.mu.Lock()
    defer dm.mu.Unlock()

    // 清空现有数据
    dm.posts = make(map[string]*Post)
    // 重新加载...

    return nil
}
```

**关键点**：

- 读多写少的场景，优先使用 `sync.RWMutex`
- 避免在锁内进行耗时操作（如文件 I/O）
- Webhook 更新时使用完整重载而非增量更新（确保数据一致性）

#### 4. **前端搜索优化**

**防抖实现**：

```javascript
// 搜索防抖，避免频繁请求
function debounce(func, wait) {
  let timeout;
  return function executedFunction(...args) {
    const later = () => {
      clearTimeout(timeout);
      func(...args);
    };
    clearTimeout(timeout);
    timeout = setTimeout(later, wait);
  };
}

// 搜索实现
const debouncedSearch = debounce(async (query) => {
  if (query.length < 2) return;

  try {
    const response = await fetch(`/api/search?q=${encodeURIComponent(query)}`);
    const results = await response.json();
    displaySearchResults(results);
  } catch (error) {
    console.error("搜索失败:", error);
  }
}, 300);
```

**性能优化点**：

- 输入长度限制（至少 2 个字符）
- 防抖延迟 300ms
- 结果数量限制
- 错误处理

#### 5. **路径处理的跨平台兼容性**

**问题描述**：
Windows 和 Linux 路径分隔符不同，需要统一处理。

**解决方案**：

```go
import "path/filepath"

// 始终使用 filepath 包进行路径操作
func (dm *DataManager) getPostsDir() string {
    // 获取可执行文件目录
    execPath, _ := os.Executable()
    execDir := filepath.Dir(execPath)

    // 跨平台路径拼接
    return filepath.Join(execDir, "posts")
}

// 生成文章 ID 时统一使用斜杠
func generatePostID(filePath, postsDir string) string {
    relPath, _ := filepath.Rel(postsDir, filePath)
    // 将路径分隔符统一为斜杠（用于 URL）
    return strings.ReplaceAll(relPath, string(filepath.Separator), "/")
}
```

### 编码规范与最佳实践

#### 1. **错误处理模式**

**统一错误处理**：

```go
// 定义错误类型
var (
    ErrPostNotFound = errors.New("文章不存在")
    ErrInvalidPath  = errors.New("无效路径")
)

// 错误包装
func (dm *DataManager) LoadPost(filePath string) (*Post, error) {
    if !filepath.IsAbs(filePath) {
        return nil, fmt.Errorf("加载文章失败: %w", ErrInvalidPath)
    }

    content, err := os.ReadFile(filePath)
    if err != nil {
        return nil, fmt.Errorf("读取文件 %s 失败: %w", filePath, err)
    }

    // 继续处理...
    return post, nil
}

// 调用方处理
post, err := dm.LoadPost(path)
if err != nil {
    if errors.Is(err, ErrPostNotFound) {
        // 特殊处理文章不存在
    }
    log.Printf("加载文章失败: %v", err)
    return
}
```

#### 2. **配置管理模式**

**集中配置管理**：

```go
// internal/config/config.go
package config

const (
    // 服务配置
    DefaultPort = 8091
    DefaultHost = "0.0.0.0"

    // 业务配置
    SummaryLines     = 3
    PageSize         = 10
    MaxSearchResults = 100

    // Webhook 配置
    WebhookSecret   = "your-secret"
    WebhookBranch   = "main"
    WebhookDevMode  = false  // 生产环境设为 false
)
```

**环境配置支持**：

```go
// 支持环境变量覆盖
func GetPort() int {
    if port := os.Getenv("MDBLOG_PORT"); port != "" {
        if p, err := strconv.Atoi(port); err == nil {
            return p
        }
    }
    return DefaultPort
}
```

#### 3. **日志记录规范**

**结构化日志**：

```go
import "log"

// 统一日志格式
func logInfo(format string, args ...interface{}) {
    log.Printf("[INFO] "+format, args...)
}

func logError(format string, args ...interface{}) {
    log.Printf("[ERROR] "+format, args...)
}

func logWebhook(format string, args ...interface{}) {
    log.Printf("[WEBHOOK] "+format, args...)
}

// 使用示例
logInfo("服务启动在端口 %d", port)
logError("加载文章失败: %v", err)
logWebhook("收到推送事件，分支: %s", branch)
```

#### 4. **测试驱动开发**

**单元测试结构**：

```go
func TestMarkdownParser(t *testing.T) {
    tests := []struct {
        name     string
        input    string
        expected *Post
        wantErr  bool
    }{
        {
            name: "有效的文章",
            input: `---
title: "测试文章"
author: "测试作者"
---
# 内容`,
            expected: &Post{
                Title: "测试文章",
                Author: "测试作者",
            },
            wantErr: false,
        },
        {
            name:    "无效的 Front Matter",
            input:   "invalid yaml",
            wantErr: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            parser := NewMarkdownParser()
            result, err := parser.Parse(tt.input)

            if tt.wantErr {
                assert.Error(t, err)
                return
            }

            assert.NoError(t, err)
            assert.Equal(t, tt.expected.Title, result.Title)
            assert.Equal(t, tt.expected.Author, result.Author)
        })
    }
}
```

#### 5. **性能监控点**

**关键性能指标**：

```go
// 启动时间监控
func (dm *DataManager) LoadAllPosts() error {
    start := time.Now()
    defer func() {
        logInfo("加载文章耗时: %v", time.Since(start))
    }()

    // 加载逻辑...
    return nil
}

// 内存使用监控
func logMemoryUsage() {
    var m runtime.MemStats
    runtime.ReadMemStats(&m)
    logInfo("内存使用: Alloc=%d KB, Sys=%d KB", m.Alloc/1024, m.Sys/1024)
}

// 搜索性能监控
func (dm *DataManager) Search(query string) ([]*Post, error) {
    start := time.Now()
    defer func() {
        logInfo("搜索 '%s' 耗时: %v", query, time.Since(start))
    }()

    // 搜索逻辑...
    return results, nil
}
```

### 部署与运维注意事项

#### 1. **服务配置最佳实践**

**Systemd 服务配置要点**：

```ini
[Unit]
Description=MDBlog Service
After=network.target

[Service]
Type=simple
User=mdblog          # 避免使用 root
Group=mdblog
WorkingDirectory=/opt/mdblog
ExecStart=/opt/mdblog/mdblog
Restart=always      # 自动重启
RestartSec=10       # 重启间隔

# 安全设置
NoNewPrivileges=true
ProtectSystem=strict
ProtectHome=true
ReadWritePaths=/opt/mdblog

[Install]
WantedBy=multi-user.target
```

#### 2. **监控与日志轮转**

**日志轮转配置**：

```bash
# /etc/logrotate.d/mdblog
/opt/mdblog/logs/*.log {
    daily
    missingok
    rotate 30
    compress
    delaycompress
    notifempty
    sharedscripts
    postrotate
        systemctl reload mdblog
    endscript
}
```

### 未来技术债务

1. **数据库支持**：当文章数量超过 1000 篇时，考虑引入轻量级数据库
2. **缓存策略**：实现增量更新而非全量重载
3. **多实例支持**：为高可用部署做准备
4. **插件系统**：模块化扩展机制
5. **国际化支持**：多语言界面和内容

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
description: "技术博客文章"
---
```

### 2. 高性能内存数据管理

**功能描述**: 全内存数据存储和管理系统

**技术特性**:

- 启动时一次性加载所有文章到内存
- 使用 `sync.RWMutex` 保证并发安全
- 多维度索引：分类索引、搜索索引
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
- 搜索范围：文章标题、内容
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

### 5. 分类系统

**功能描述**: 基于目录结构的自动分类

**分类规则**:

- 自动识别：`posts/go/article.md` → 分类 "go"
- 根目录文章：默认分类 "其他"
- 特殊处理：`about.md` 仅在关于页面显示

**分类功能**:

- 分类页面：`/category/Go`
- 文章按分类组织展示

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

### 9. Webhook 自动同步

**功能描述**: Git 仓库变更时自动更新博客内容

**核心功能**:

- HTTP 接口：`POST /webhook/gitee`
- 签名验证：HMAC-SHA256 安全校验
- Git 操作：自动执行 `git pull` 更新内容
- 内存数据重新加载：实时更新博客数据
- 错误处理和日志记录：完整的错误追踪
- 健康检查接口：`GET /health`
- Webhook 测试工具：开发和生产环境测试

**技术实现**:

```go
// Webhook 配置
const (
    WebhookSecret = "c85b7544d37c89280afcf912ec70f4083b18065d9e89966cae6059798f0dadf5"
    WebhookBranch = "main"     // 监听的分支
)
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
                       │   Category Index│    │   (HTML/JSON)   │
                       └─────────────────┘    └─────────────────┘
```

### API 接口设计

#### 页面路由

```
GET  /                  # 首页
GET  /post/:id          # 文章详情页
GET  /category/:category # 分类页面
GET  /search            # 搜索页面
GET  /about             # 关于页面
```

#### API 接口

```
GET  /api/posts         # 文章列表API
GET  /api/posts/:id     # 文章详情API
GET  /api/categories    # 分类列表API
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
    ID          string    `json:"id"`          // 文章ID
    Title       string    `json:"title"`       // 标题
    Author      string    `json:"author"`      // 作者
    GitHubURL   string    `json:"github_url"`  // GitHub链接
    Content     string    `json:"content"`     // 原始内容
    HTML        string    `json:"html"`        // HTML内容
    Summary     string    `json:"summary"`     // 摘要
    Category    string    `json:"category"`    // 分类
    CreateTime  time.Time `json:"created_at"`  // 创建时间
    UpdateTime  time.Time `json:"updated_at"`  // 更新时间
    Description string    `json:"description"` // 描述
    FilePath    string    `json:"file_path"`   // 文件路径
}

// BlogData 全局数据结构
type BlogData struct {
    Posts       map[string]*Post    // 文章索引
    Categories  map[string][]string // 分类索引
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
5. **Webhook 签名验证**: HMAC-SHA256 签名校验
6. **分支过滤**: 只处理指定分支的推送事件

## 📊 监控指标

### 应用监控

- **服务状态**: Systemd 服务监控
- **端口检查**: TCP 8091 端口可用性
- **日志监控**: 应用错误日志
- **内存使用**: 进程内存占用

### 业务指标

- **文章数量**: 加载的文章总数
- **搜索频率**: 搜索 API 调用统计

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
