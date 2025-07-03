# MDBlog

一个简单高效的 Markdown 博客系统，使用 Go 语言开发。支持单文件部署、实时搜索、响应式设计和 Webhook 自动同步。

> **V2 版本说明**：重构后的版本，将 Markdown 文件迁移到 Gitee。旧版本请查看 [V1 分支](https://github.com/broqiang/mdblog/releases/tag/v1.1.0)

## ✨ 核心特性

- 📝 **Markdown 驱动** - 基于文件系统，无需数据库
- 🔍 **实时搜索** - 快捷键搜索（双击 Cmd/Ctrl）
- 📱 **响应式设计** - 移动端优化，现代化 UI
- ⚡ **高性能** - 内存缓存，毫秒级响应
- 📦 **单文件部署** - 静态资源嵌入，零依赖
- 🔄 **自动同步** - Gitee Webhook 自动更新
- 🎨 **代码高亮** - 支持多种编程语言
- 🚀 **一键部署** - 自动化编译、上传、重启

## 📸 界面预览

### 首页展示

![首页展示](https://image.broqiang.com/mdblog/mdblog_index.png)

### 文章详情

![文章详情](https://image.broqiang.com/mdblog/mdblog_detail.png)

## 🚀 快速开始

### 环境要求

- Go 1.24.4+
- Git
- 开发环境：macOS（测试）
- 生产环境：Debian 6.1（测试）

### 本地开发

```bash
# 1. 克隆项目
git clone https://github.com/broqiang/mdblog.git
cd mdblog

# 2. 安装依赖
go mod tidy

# 3. 本地运行
make dev
# 或者
go run main.go

# 4. 访问应用
# http://localhost:8091
```

### 创建文章

在 `posts/` 目录下创建 Markdown 文件：

```markdown
---
title: "文章标题"
author: "BroQiang"
github_url: "https://github.com/broqiang/mdblog"
created_at: "2024-01-01T10:00:00Z"
updated_at: "2024-01-15T15:30:00Z"
description: "文章描述"
---

# 文章内容

这里是文章的具体内容...
```

> **注意**：生产环境中，posts 目录是独立的 Gitee 仓库。目前仅支持 Gitee Webhook，如需其他 Git 平台，请修改 `internal/server/webhook.go` 中的签名验证算法。

## 📁 项目结构

```
mdblog/
├── main.go              # 主程序入口
├── embed.go             # 静态资源嵌入
├── go.mod               # Go 依赖管理
├── go.sum               # Go 依赖校验
├── Makefile             # 构建脚本
├── posts/               # Markdown 文章目录
├── internal/            # 内部模块
│   ├── config/          # 配置管理
│   ├── data/            # 数据管理和缓存
│   ├── markdown/        # Markdown 解析器
│   ├── server/          # Web 服务器
│   └── assets/          # 内部静态资源
├── web/                 # 前端资源
│   ├── static/          # 静态文件（CSS/JS/图片等）
│   └── templates/       # HTML 模板
├── cmd/                 # 命令行工具
├── deploy/              # 部署配置
└── docs/                # 文档
```

## 📝 文章格式说明

### Front Matter 字段

```yaml
---
title: "文章标题" # 必填
author: "作者名" # 可选，默认 "BroQiang"
github_url: "GitHub链接" # 可选
created_at: "2024-01-01T10:00:00Z" # 可选
updated_at: "2024-01-01T12:00:00Z" # 可选
description: "文章描述" # 可选
---
```

### 分类规则

- 基于目录结构自动分类
- `posts/go/article.md` → 分类 "go"
- `posts/article.md` → 分类 "其他"

### 支持的 Markdown 语法

- GitHub Flavored Markdown (GFM)
- 代码语法高亮
- 表格、任务列表、删除线等扩展语法
- 数学公式、脚注等

## 🛠️ 生产部署

### 1. 创建用户和目录

```bash
# 创建专用用户（/bro 为家目录，可根据实际情况修改）
sudo useradd -m -d /bro -s /bin/bash bro
sudo mkdir -p /bro/mdblog/{posts,logs}
sudo chown -R bro:bro /bro/mdblog
```

### 2. 安装 systemd 服务

```bash
# 上传服务配置文件
scp deploy/mdblog.service user@server:/tmp/

# 安装服务
sudo cp /tmp/mdblog.service /etc/systemd/system/
sudo systemctl daemon-reload
sudo systemctl enable mdblog
```

### 3. 配置文章仓库

```bash
# 用 bro 用户克隆文章仓库
sudo -u bro git clone git@gitee.com:your-username/posts.git /bro/mdblog/posts

# 确保目录权限
sudo chown -R bro:bro /bro/mdblog/posts
```

### 4. 部署应用

编辑 `Makefile` 配置：

```makefile
SERVER_HOST=your-server-ip
SERVER_USER=root
SERVER_PATH=/bro/mdblog
SSH_PORT=22
APP_PORT=8091
```

执行部署：

```bash
# 编译并部署
make build
make scp
```

## 🔄 Webhook 自动同步

### 配置 SSH 密钥

```bash
# 在服务器上生成密钥
sudo -u bro ssh-keygen -t rsa -b 4096 -C "bro@your-server"

# 查看公钥并添加到 Gitee
sudo -u bro cat /bro/.ssh/id_rsa.pub
```

### 配置 Gitee Webhook

1. 在 Gitee 仓库设置中添加 Webhook：

   - **URL**: `https://your-domain.com/webhook/gitee`
   - **密码**: 设置一个安全密钥
   - **触发事件**: Push
   - **分支过滤**: main

2. 更新应用配置 `internal/config/config.go`：

```go
const (
    WebhookSecret  = "your-webhook-secret"
    WebhookBranch  = "main"
    WebhookDevMode = false  // 生产环境设为 false
)
```

3. 重新部署应用：

```bash
make scp
```

### 测试 Webhook

```bash
# 本地测试推送
cd posts
echo "# 测试文章" > test.md
git add . && git commit -m "测试" && git push origin main

# 检查同步状态
curl https://your-domain.com/health
sudo tail -f /bro/mdblog/logs/mdblog.log | grep -i webhook
```

## 📊 服务管理

### 常用命令

```bash
# 服务控制
sudo systemctl start mdblog       # 启动
sudo systemctl stop mdblog        # 停止
sudo systemctl restart mdblog     # 重启
sudo systemctl status mdblog      # 状态

# 日志查看
sudo tail -f /bro/mdblog/logs/mdblog.log    # 应用日志
sudo journalctl -u mdblog -f                # 系统日志

# 健康检查
curl http://localhost:8091/health
netstat -tlnp | grep 8091
```

### 日志轮转

创建 `/etc/logrotate.d/mdblog`：

```bash
/bro/mdblog/logs/*.log {
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

### 性能监控

```bash
# 资源使用
ps aux | grep mdblog
du -sh /bro/mdblog/

# 文章统计
find /bro/mdblog/posts -name "*.md" | wc -l
```

## 🚨 故障排查

### 服务启动问题

```bash
# 检查服务状态
sudo systemctl status mdblog
sudo journalctl -u mdblog --no-pager

# 检查端口占用
sudo netstat -tlnp | grep 8091
```

### Webhook 同步问题

```bash
# 检查 Git 配置
cd /bro/mdblog/posts
sudo -u bro git status
sudo -u bro git remote -v

# 测试 SSH 连接
sudo -u bro ssh -T git@gitee.com

# 检查文件权限
ls -la /bro/mdblog/posts/
```

### 权限问题解决

```bash
# 重置权限
sudo chown -R bro:bro /bro/mdblog
sudo chmod 755 /bro/mdblog
sudo chmod 755 /bro/mdblog/posts/
sudo chmod 644 /bro/mdblog/posts/*.md
```

## 🔧 配置选项

### 命令行参数

```bash
# 默认配置
./mdblog

# 自定义 posts 目录
./mdblog -posts /path/to/posts
```

### 应用配置

主要配置在 `internal/config/config.go`：

```go
const (
    DefaultPort      = 8091
    DefaultHost      = "0.0.0.0"
    SummaryLines     = 3
    PageSize         = 10
    WebhookSecret    = "secret"
    WebhookBranch    = "main"
    WebhookDevMode   = false
)
```

## 📡 API 接口

### 文章接口

```http
GET /api/posts?page=1&size=10          # 获取文章列表
GET /api/posts/{id}                     # 获取单篇文章
GET /api/categories                     # 获取分类列表
GET /api/search?q=关键词&page=1&size=10  # 搜索文章
```

### 系统接口

```http
GET /health                             # 健康检查
POST /webhook/gitee                     # Gitee Webhook
```

## 🎨 功能特性

### 搜索功能

- **快捷键**：双击 `Cmd`（macOS）或 `Ctrl`（Windows/Linux）
- **搜索范围**：文章标题和内容
- **实时结果**：输入即搜索

### 响应式设计

- **桌面端**：侧边导航，宽屏布局
- **移动端**：汉堡菜单，触摸优化

## 🔒 安全建议

1. **用户权限**：使用专用用户 `bro` 运行服务，避免 root 权限
2. **防火墙**：只开放必要端口（8091、SSH 端口）
3. **HTTPS**：使用 SSL 证书加密传输
4. **备份策略**：定期备份 posts 目录和配置文件
5. **日志监控**：定期检查访问日志，发现异常行为

## 📄 许可证

MIT License - 详见 [LICENSE](LICENSE) 文件

## 🔗 相关链接

- [产品需求文档](docs/prd.md)
- [项目仓库](https://github.com/broqiang/mdblog)

---

**MDBlog** - 让 Markdown 写作更简单，让博客部署更轻松！
