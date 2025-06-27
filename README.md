# MDlog

一个简单的基于 Markdown 的博客系统，使用 Go 语言编写。

## 特性

- 📝 基于 Markdown 文件的博客系统
- 🏷️ 支持分类和标签
- 🔍 实时搜索功能
- 📱 响应式设计，移动端优化
- ⚡ 快速、轻量、单文件部署

## 快速开始

### 安装依赖

```bash
go mod tidy
```

### 开发模式

```bash
make dev
```

### 编译

```bash
make build
```

### 部署到服务器

1. **编辑 Makefile 中的服务器配置**：

   ```makefile
   SERVER_HOST=your-server-ip
   SERVER_USER=root
   SERVER_PATH=/bro/mdblog
   SERVER_PORT=8080
   ```

2. **首次部署**：

   ```bash
   # 在服务器上创建目录和安装服务
   ssh your-user@your-server
   sudo mkdir -p /bro/mdblog/posts /bro/mdblog/logs
   sudo cp deploy/mdblog.service /etc/systemd/system/
   sudo chown -R root:root /bro/mdblog
   sudo systemctl daemon-reload
   sudo systemctl enable mdblog
   ```

3. **日常更新**：
   ```bash
   make scp  # 编译、上传、重启服务并检测启动状态
   ```

## 命令行参数

```bash
./mdblog -h
```

参数说明：

- `-host string`: 服务器地址（默认：0.0.0.0）
- `-port int`: 服务器端口（默认：8080）
- `-posts string`: posts 目录路径（默认：可执行文件同级的 posts 目录）

## 使用示例

```bash
# 使用默认配置
./mdblog

# 指定端口
./mdblog -port 9000

# 指定posts目录
./mdblog -posts /path/to/posts

# 组合使用
./mdblog -host 127.0.0.1 -port 9000 -posts ./content
```

## 项目结构

```
mdblog/
├── cmd/mdblog/          # 主程序入口
├── internal/
│   ├── config/          # 配置常量
│   ├── data/            # 数据管理
│   ├── markdown/        # Markdown解析
│   └── server/          # Web服务器
├── posts/               # Markdown文章目录
├── web/                 # 静态资源和模板
├── Makefile            # 构建脚本
└── README.md
```

## 配置说明

所有配置都硬编码在 `internal/config/config.go` 中，如需修改请直接编辑该文件后重新编译。

默认配置：

- 服务器端口：8080
- 分页大小：10
- 摘要行数：3
- 搜索结果数：100

## 文章格式

在 `posts/` 目录中创建 Markdown 文件，支持 Front Matter：

```markdown
---
title: "文章标题"
author: "作者"
date: 2024-01-01T00:00:00Z
category: "分类"
tags: ["标签1", "标签2"]
---

文章内容...
```

## 构建命令

- `make dev` - 开发模式运行
- `make build` - 交叉编译 Linux 可执行文件
- `make scp` - 编译、上传到服务器、重启服务并检测启动状态

### scp 命令详细流程

`make scp` 命令会自动执行以下步骤：

1. **编译** - 交叉编译 Linux 版本的可执行文件
2. **上传** - 通过 scp 上传到服务器指定目录
3. **重启** - 通过 ssh 执行 `sudo systemctl restart mdblog`
4. **等待** - 等待 5 秒让服务完全启动
5. **检测** - 使用 nc 命令检测端口是否可访问
6. **反馈** - 显示服务状态，失败时显示错误日志

成功示例输出：

```
✅ 服务启动成功，端口 8080 可访问
● mdblog.service - MDlog - Markdown Blog Server
   Active: active (running)
```

失败时会显示详细的错误信息和日志。

## 服务管理

### Systemd 服务命令

```bash
# 查看服务状态
sudo systemctl status mdblog

# 启动服务
sudo systemctl start mdblog

# 停止服务
sudo systemctl stop mdblog

# 重启服务
sudo systemctl restart mdblog

# 查看系统日志
sudo journalctl -u mdblog -f

# 查看应用日志
tail -f /bro/mdblog/logs/mdblog.log

# 开机自启
sudo systemctl enable mdblog

# 禁用开机自启
sudo systemctl disable mdblog
```

### 日志管理

应用日志统一输出到 `/bro/mdblog/logs/mdblog.log` 文件中，包含：

- **标准输出**: 应用的正常运行日志
- **标准错误**: 应用的错误信息和异常

常用日志查看命令：

```bash
# 实时查看日志
tail -f /bro/mdblog/logs/mdblog.log

# 查看最新100行日志
tail -100 /bro/mdblog/logs/mdblog.log

# 查看日志文件大小
ls -lh /bro/mdblog/logs/mdblog.log

# 清空日志文件（谨慎操作）
sudo truncate -s 0 /bro/mdblog/logs/mdblog.log
```

**注意**: 日志文件会随着时间增长，建议定期清理或使用 logrotate 进行日志轮转。

## License

MIT License
