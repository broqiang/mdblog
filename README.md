# MDBlog

一个简单的基于 Markdown 的博客系统，使用 Go 语言编写，支持单文件部署、实时搜索、响应式设计。

好久没有处理这个博客项目了，今天新写了个文章， 发现阿里云不能 pull Github 的仓库了，
将 markdown 文件改到 gitee.com 了， 顺便将代码重构了一下

目前这个是重构的 V2 分支，旧的版本请查看 [V1 分支](https://github.com/broqiang/mdblog/releases/tag/v1.1.0)

## ✨ 特性

- 📝 **Markdown 驱动**：基于 Markdown 文件的博客系统，支持 Front Matter
- 🏷️ **分类标签**：自动识别目录结构生成分类，支持多标签系统
- 🔍 **实时搜索**：快捷键搜索（双击 Cmd/Ctrl），支持标题和内容模糊匹配
- 📱 **响应式设计**：移动端优化，汉堡菜单，现代化 UI
- ⚡ **高性能**：内存缓存，单文件部署，无数据库依赖
- 🎨 **代码高亮**：基于 Chroma 的语法高亮，支持多种编程语言
- 📦 **嵌入资源**：静态文件嵌入到可执行文件，真正的单文件部署
- 🔧 **零配置**：开箱即用，最小化配置需求
- 🚀 **自动部署**：一键编译上传重启，支持自定义 SSH 端口

## 🚀 快速开始

### 安装依赖

```bash
git clone https://github.com/your-username/mdblog.git
cd mdblog
go mod tidy
```

### 开发模式

```bash
make dev
# 或者直接运行
go run main.go
```

访问 http://localhost:8091

### 编译部署

```bash
# 本地编译
make build

# 编译并部署到服务器
make scp
```

## 📖 使用方式

### 命令行参数

```bash
# 使用默认配置（自动检测posts目录）
./mdblog

# 指定自定义posts目录
./mdblog -posts /path/to/custom/posts
```

### 服务器配置

- **监听地址**: 0.0.0.0:8091（固定）
- **Posts 目录**: 可执行文件同级的 `posts` 目录（可通过 `-posts` 参数自定义）

如需修改监听地址和端口，请编辑 `internal/config/config.go` 文件后重新编译。

## 📁 项目结构

```
mdblog/
├── main.go                 # 主程序入口（嵌入静态资源）
├── embed.go                # 静态资源嵌入声明
├── internal/
│   ├── config/             # 配置常量
│   ├── data/               # 数据管理和缓存
│   ├── markdown/           # Markdown解析器
│   └── server/             # Web服务器和路由
├── posts/                  # Markdown文章目录
│   ├── about.md           # 关于页面
│   ├── category1/         # 分类目录
│   └── category2/         # 分类目录
├── web/                   # 前端资源（已嵌入）
│   ├── static/            # CSS、JS、图片
│   └── templates/         # HTML模板
├── deploy/                # 部署相关文件
├── docs/                  # 文档
├── Makefile              # 构建和部署脚本
└── README.md
```

## 📝 文章格式

在 `posts/` 目录中创建 Markdown 文件，支持 YAML Front Matter：

```markdown
---
title: "文章标题"
author: "作者名称"
github_url: "https://github.com/username"
created_at: 2024-01-01T10:00:00
updated_at: 2024-01-01T12:00:00
tags: ["Go", "Web开发", "技术"]
---

# 文章内容

这里是 Markdown 格式的文章内容...

\`\`\`go
func main() {
fmt.Println("Hello, World!")
}
\`\`\`
```

### Front Matter 字段说明

- `title`: 文章标题（必填）
- `author`: 作者（可选，默认"BroQiang"）
- `github_url`: GitHub 链接（可选）
- `created_at`: 创建时间（可选，支持多种时间格式）
- `updated_at`: 更新时间（可选）
- `description`: 文章描述（可选）
- `tags`: 标签数组（可选）

### 分类规则

- 分类基于文件所在目录自动生成
- 例如：`posts/go/article.md` → 分类为 "go"
- 根目录文章分类为 "其他"

## 🛠️ 构建命令

| 命令         | 说明                                 |
| ------------ | ------------------------------------ |
| `make dev`   | 开发模式运行                         |
| `make build` | 交叉编译 Linux 可执行文件            |
| `make scp`   | 编译、停止服务、上传、重启、检测状态 |

### 部署流程详解

`make scp` 执行的完整流程：

1. **编译** - 交叉编译 Linux 版本可执行文件
2. **停止服务** - SSH 连接服务器停止 mdblog 服务
3. **上传文件** - 使用 SCP 上传新的可执行文件
4. **启动服务** - SSH 启动 mdblog 服务
5. **状态检测** - 检测端口 8091 是否可访问
6. **日志反馈** - 显示服务状态和最新日志

### 服务器配置

编辑 `Makefile` 中的服务器配置：

```makefile
# 服务器配置
SERVER_HOST=123.56.186.148    # 服务器IP
SERVER_USER=root              # SSH用户名
SERVER_PATH=/bro/mdblog       # 部署路径
SSH_PORT=21345               # SSH端口
APP_PORT=8091                # 应用端口
```

## 🔧 服务器部署

### 首次部署

1. **准备服务器环境**：

```bash
# 连接服务器
ssh -p 21345 root@your-server-ip

# 创建目录
sudo mkdir -p /bro/mdblog/posts /bro/mdblog/logs

# 复制服务文件（需要先上传 deploy/mdblog.service）
sudo cp /bro/mdblog/deploy/mdblog.service /etc/systemd/system/

# 设置权限
sudo chown -R root:root /bro/mdblog

# 启用服务
sudo systemctl daemon-reload
sudo systemctl enable mdblog
```

2. **同步文章目录**：

```bash
# 将本地 posts 目录同步到服务器
scp -P 21345 -r posts/ root@your-server-ip:/bro/mdblog/
```

3. **首次部署**：

```bash
make scp
```

### 日常更新

```bash
# 更新代码和重新部署
make scp

# 仅更新文章内容
scp -P 21345 -r posts/ root@your-server-ip:/bro/mdblog/
ssh -p 21345 root@your-server-ip "sudo systemctl restart mdblog"
```

## 📊 服务管理

### Systemd 服务命令

```bash
# 查看服务状态
sudo systemctl status mdblog

# 启动/停止/重启服务
sudo systemctl start mdblog
sudo systemctl stop mdblog
sudo systemctl restart mdblog

# 查看服务日志
sudo journalctl -u mdblog -f

# 查看应用日志
tail -f /bro/mdblog/logs/mdblog.log

# 开机自启动
sudo systemctl enable mdblog
sudo systemctl disable mdblog
```

### 日志管理

应用日志位置：`/bro/mdblog/logs/mdblog.log`

```bash
# 实时查看日志
tail -f /bro/mdblog/logs/mdblog.log

# 查看最新100行日志
tail -100 /bro/mdblog/logs/mdblog.log

# 清空日志（谨慎操作）
sudo truncate -s 0 /bro/mdblog/logs/mdblog.log
```

## 🔍 功能特性详解

### 搜索功能

- **快捷键**: 双击 `Cmd`（macOS）或 `Ctrl`（Windows/Linux）
- **搜索范围**: 文章标题和内容
- **搜索算法**: 模糊匹配，大小写不敏感
- **实时结果**: 输入即搜索，无需点击按钮

### 响应式设计

- **桌面端**: 侧边导航，宽屏布局
- **平板端**: 折叠导航，适中布局
- **移动端**: 汉堡菜单，单列布局，触摸优化

### 性能优化

- **内存缓存**: 所有文章内容加载到内存
- **嵌入资源**: CSS/JS/图片嵌入可执行文件
- **单文件部署**: 无外部依赖，部署简单
- **并发安全**: 读写锁保护数据结构

## ⚙️ 配置选项

所有配置位于 `internal/config/config.go`：

```go
const (
    DefaultPort         = 8091    // 服务端口
    DefaultHost         = "0.0.0.0"  // 监听地址
    SummaryLines        = 3       // 摘要行数
    PageSize           = 10      // 分页大小
    MaxSearchResults   = 100     // 最大搜索结果
)
```

## 🎨 自定义主题

CSS 文件位于 `web/static/css/main.css`，支持：

- 深色/浅色主题切换
- 自定义颜色方案
- 响应式断点调整
- 代码高亮主题

## 🤝 贡献指南

1. Fork 项目
2. 创建功能分支 (`git checkout -b feature/amazing-feature`)
3. 提交更改 (`git commit -m 'Add amazing feature'`)
4. 推送到分支 (`git push origin feature/amazing-feature`)
5. 创建 Pull Request

## 📄 许可证

MIT License - 详见 [LICENSE](LICENSE) 文件

## 🔗 相关链接

- [项目文档](docs/)
- [产品需求文档](docs/prd.md)
- [部署指南](deploy/)

---

**MDlog** - 让 Markdown 写作更简单，让博客部署更轻松！
