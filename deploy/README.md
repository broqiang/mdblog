# 部署文件

本目录包含 MDlog 项目的部署相关文件。

## 文件说明

- `mdblog.service` - systemd 服务配置文件

## 部署步骤

### 1. 首次部署

在服务器上执行以下命令：

```bash
# 创建应用目录
sudo mkdir -p /bro/mdblog/posts /bro/mdblog/logs

# 安装systemd服务
sudo cp deploy/mdblog.service /etc/systemd/system/

# 设置权限
sudo chown -R root:root /bro/mdblog

# 重载并启用服务
sudo systemctl daemon-reload
sudo systemctl enable mdblog
```

### 2. 上传程序文件

使用项目根目录的 Makefile：

```bash
# 编译、上传并重启服务
make scp
```

### 3. 服务管理

```bash
# 查看服务状态
sudo systemctl status mdblog

# 重启服务
sudo systemctl restart mdblog

# 查看日志
tail -f /bro/mdblog/logs/mdblog.log
```

## 目录结构

服务器上的目录结构：

```
/bro/mdblog/
├── mdblog          # 可执行文件
├── posts/          # Markdown文章目录
└── logs/           # 日志目录
    └── mdblog.log  # 应用日志
```

## 注意事项

- 服务以 root 用户运行
- 日志文件位于 `/bro/mdblog/logs/mdblog.log`
- 服务会在系统启动时自动启动
- 进程异常退出时会自动重启
