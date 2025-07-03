# 🔗 Webhook 自动同步使用指南

## 📖 功能概述

Webhook 自动同步功能允许你在 Gitee 仓库发生推送时，自动更新博客内容。当你向 posts 仓库提交新文章或修改现有文章时，博客会自动拉取最新内容并刷新缓存。

## 🚀 快速开始

### 1. 配置 Gitee Webhook

1. 登录 Gitee，进入你的 posts 仓库
2. 点击 **管理** → **WebHooks**
3. 点击 **添加 WebHook**
4. 填写以下配置：

```
URL: https://broqiang.com/webhook/gitee
密码: c85b7544d37c89280afcf912ec70f4083b18065d9e89966cae6059798f0dadf5
触发事件: ✅ Push
分支过滤: main
```

5. 点击 **添加** 完成配置

### 2. 测试 Webhook

使用内置的测试工具验证配置：

```bash
# 测试本地开发环境
make webhook-test-local

# 测试生产环境
make webhook-test-remote
```

### 3. 监控同步状态

```bash
# 查看服务健康状态
curl https://broqiang.com/health

# 查看 Webhook 日志
tail -f /bro/mdblog/logs/mdblog.log | grep -i webhook
```

## 🔧 工作原理

### 同步流程

1. **触发**: 向 posts 仓库 main 分支推送代码
2. **验证**: Gitee 发送带签名的 Webhook 请求
3. **签名验证**: 服务器验证请求的合法性
4. **Git Pull**: 执行 `git pull origin main` 更新代码
5. **数据重载**: 清空缓存，重新解析所有 Markdown 文件
6. **响应**: 返回同步结果和统计信息

### 安全机制

- **HMAC-SHA256 签名**: 防止恶意请求
- **分支过滤**: 只处理 main 分支的推送
- **请求验证**: 验证方法、头部和载荷格式
- **错误处理**: 全面的异常捕获和日志记录

## 📝 使用示例

### 添加新文章

1. 在本地 posts 仓库创建新文章：

```bash
cd posts
mkdir -p go
cat > go/new-article.md << 'EOF'
---
title: "我的新文章"
author: "BroQiang"
created_at: 2024-01-15T10:00:00
tags: ["Go", "技术"]
---

# 文章内容

这是一篇新文章的内容...
EOF
```

2. 提交并推送：

```bash
git add .
git commit -m "添加新文章: 我的新文章"
git push origin main
```

3. 等待几秒钟，访问博客查看更新

### 修改现有文章

1. 编辑文章文件
2. 提交并推送更改
3. Webhook 自动触发，文章立即更新

## 🐛 故障排查

### Webhook 未触发

1. **检查 Gitee 配置**:

   ```bash
   # 确认 URL 正确
   curl https://broqiang.com/health
   ```

2. **检查分支**:

   ```bash
   # 确保推送到 main 分支
   git branch
   git push origin main
   ```

3. **查看 Gitee Webhook 日志**:
   - 在 Gitee 仓库的 WebHook 页面查看发送日志
   - 检查返回状态码和错误信息

### 同步失败

1. **查看服务器日志**:

   ```bash
   tail -n 50 /bro/mdblog/logs/mdblog.log
   ```

2. **检查 Git 仓库状态**:

   ```bash
   cd ../posts
   git status
   git remote -v
   ```

3. **手动测试**:

   ```bash
   # 测试 Git 操作
   cd ../posts
   git pull origin main

   # 测试 Webhook
   make webhook-test-remote
   ```

### 常见错误

| 错误信息           | 原因               | 解决方案                |
| ------------------ | ------------------ | ----------------------- |
| `签名验证失败`     | Webhook 密钥不匹配 | 检查 Gitee 配置中的密钥 |
| `git pull 失败`    | Git 仓库问题       | 检查仓库状态和权限      |
| `posts 目录不存在` | 路径配置错误       | 确认 posts 目录位置     |
| `已忽略非目标分支` | 推送到非 main 分支 | 推送到 main 分支        |

## ⚙️ 高级配置

### 修改监听分支

编辑 `internal/config/config.go`:

```go
const (
    WebhookBranch = "develop" // 改为其他分支
)
```

### 修改 posts 目录路径

posts 目录路径通过以下方式确定：

1. **命令行参数**（优先级最高）：

   ```bash
   ./mdblog -posts /path/to/posts
   ```

2. **默认位置**：
   - 正常编译的可执行文件：可执行文件同级的 `posts` 目录
   - 开发模式（`go run`）：当前工作目录下的 `posts` 目录

如果需要修改默认行为，可以编辑 `main.go` 中的路径逻辑：

```go
// 正常编译的可执行文件，使用可执行文件同级目录
actualPostsDir = filepath.Join(execDir, "posts")
```

### 自定义 Webhook 密钥

1. 生成新密钥：

   ```bash
   openssl rand -hex 32
   ```

2. 更新配置：

   ```go
   const (
       WebhookSecret = "your-new-secret-key"
   )
   ```

3. 更新 Gitee Webhook 配置

## 🔍 监控和日志

### 日志格式

Webhook 相关日志使用以下格式：

```
2024-01-15 10:30:45 [INFO] 收到 Gitee Webhook 请求
2024-01-15 10:30:45 [INFO] 开始执行 Git Pull...
2024-01-15 10:30:46 [INFO] Git Pull 成功: Already up to date.
2024-01-15 10:30:46 [INFO] 开始重新加载数据...
2024-01-15 10:30:46 [INFO] 重新加载文章: 我的新文章
2024-01-15 10:30:46 [INFO] 数据重新加载完成，共加载 25 篇文章
2024-01-15 10:30:46 [INFO] Webhook 处理成功完成
```

### 性能监控

```bash
# 查看内存使用
ps aux | grep mdblog

# 查看端口状态
netstat -tlnp | grep 8091

# 查看最近的 Webhook 请求
grep "Webhook" /bro/mdblog/logs/mdblog.log | tail -10
```

## 📊 最佳实践

### 1. 文章管理

- 使用有意义的文件名
- 遵循分类目录结构
- 确保 Front Matter 格式正确
- 定期检查文章的 Git 状态

### 2. 安全建议

- 定期更换 Webhook 密钥
- 监控异常的同步请求
- 备份重要文章内容
- 使用分支保护策略

### 3. 性能优化

- 避免频繁的小提交
- 合并相关更改为单次提交
- 监控同步时间和资源使用
- 定期清理无用的文章文件

## 🆘 技术支持

如果遇到问题，可以：

1. 查看本指南的故障排查部分
2. 检查项目的 GitHub Issues
3. 查看详细的错误日志
4. 使用测试工具验证配置

## ⚠️ 重要说明：签名验证问题

### 当前状态

Gitee Webhook 的签名验证算法与标准的 GitHub/GitLab 不同，经过测试发现多种签名算法都无法匹配。

### 解决方案

#### 方案一：开发模式（推荐用于测试）

在 `internal/config/config.go` 中设置：

```go
WebhookDevMode = true  // 跳过签名验证
```

#### 方案二：生产环境解决方案

1. **网络安全**: 使用防火墙或 Nginx 限制只允许 Gitee 服务器 IP 访问
2. **HTTPS**: 确保使用 HTTPS 加密传输
3. **内网部署**: 如果可能，将 Webhook 端点部署在内网

#### 方案三：自定义验证

你可以在 `verifyGiteeSignature` 函数中实现自己的验证逻辑，比如：

- IP 白名单验证
- 请求头验证
- 时间戳验证（防重放攻击）

### 生产环境配置

生产环境使用时，请：

1. **修改配置**：

```go
// internal/config/config.go
WebhookDevMode = false  // 启用签名验证
WebhookSecret = "你的实际密钥"  // 替换为你的密钥
```

2. **Nginx 反向代理配置**：

```nginx
# 限制 Gitee IP 访问 webhook
location /webhook/gitee {
    # Gitee 服务器 IP 段（需要查询最新的）
    allow 116.211.167.0/24;
    allow 117.184.188.0/24;
    deny all;

    proxy_pass http://localhost:8091;
    proxy_set_header Host $host;
    proxy_set_header X-Real-IP $remote_addr;
    proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
}
```

### 调试签名验证

如果你想解决签名验证问题，可以查看日志输出：

```bash
# 查看详细的签名验证日志
tail -f /var/log/your-app.log | grep "签名验证"
```

当前实现会尝试多种签名算法：

1. 直接密钥比较
2. `base64(sha256(timestamp + secret + body))`
3. `base64(sha256(secret + timestamp + body))`
4. `base64(sha256(body + timestamp + secret))`
5. `base64(sha256(secret))`

### 修改 posts 目录路径

---

**更新时间**: 2024-01-15  
**版本**: v1.1.0

## ⚠️ 重要发现：Gitee Token 验证机制

经过深入研究多个开源项目和实际测试，发现了 Gitee Webhook 的关键特点：

### Gitee vs GitHub 签名验证对比

| 平台       | 验证方式   | 算法        | 传递方式                          |
| ---------- | ---------- | ----------- | --------------------------------- |
| **GitHub** | HMAC 签名  | HMAC-SHA256 | `X-Hub-Signature-256: sha256=...` |
| **Gitee**  | 明文 Token | 无加密      | `X-Gitee-Token: 原文密码`         |

### 关键发现

1. **Gitee 不使用 HMAC 签名验证**

   - 与 GitHub 的 HMAC-SHA256 完全不同
   - Gitee 直接传递明文 Token，不进行加密

2. **Token 生成方式**

   - 用户在 Gitee 设置的"密码"会被系统处理
   - 实际发送的 Token 可能与用户设置的密码不同
   - Token 似乎是由 Gitee 系统生成的 Base64 编码值

3. **实际请求示例**
   ```http
   POST /webhook/gitee
   X-Gitee-Token: 2ovtNdnFjPV6yXQynqrzdBR0Rpa4w9cHDvHspg3MFAA=
   X-Gitee-Timestamp: 1751520876267
   Content-Type: application/json
   ```

## 解决方案

### 方案 1：开发模式（推荐用于测试）

在 `internal/config/config.go` 中设置：

```go
WebhookDevMode = true  // 跳过Token验证
```

**优点**：立即可用，适合开发测试
**缺点**：没有安全验证

### 方案 2：Token 适配（推荐用于生产）

1. **获取真实 Token**：
   - 临时启用开发模式
   - 触发一次 Webhook 推送
   - 查看日志中的 `X-Gitee-Token` 值
2. **更新配置**：
   ```go
   WebhookSecret = "2ovtNdnFjPV6yXQynqrzdBR0Rpa4w9cHDvHspg3MFAA="  // 使用实际Token
   WebhookDevMode = false  // 启用验证
   ```

### 方案 3：生产环境安全策略

由于 Gitee 的 Token 验证机制限制，建议生产环境采用多层安全策略：

1. **网络层安全**：

   - 配置 Nginx 反向代理
   - 限制访问来源 IP（Gitee 服务器 IP）
   - 使用 HTTPS 加密传输

2. **应用层安全**：
   - 设置复杂的 Webhook URL 路径
   - 启用访问日志监控
   - 设置失败次数限制

## 配置步骤

### 1. 系统配置

确保在 `internal/config/config.go` 中有正确的设置：

```go
// Webhook 配置
WebhookSecret  = "你的Token"    // 根据上述方案设置
WebhookBranch  = "main"         // 监听的分支
WebhookDevMode = false          // 生产环境设为false
```

### 2. Gitee 仓库配置

1. 访问：`https://gitee.com/你的用户名/仓库名/hooks`
2. 点击"添加 Webhook"
3. 填写信息：
   - **URL**: `https://broqiang.com/webhook/gitee`
   - **密码**: 任意设置（注意后续需要适配）
   - **请求内容类型**: `application/json`
4. 选择触发事件：**Push**
5. 点击"添加"

### 3. 测试验证

1. **启用开发模式**进行初始测试
2. **提交代码**到 main 分支
3. **检查日志**：
   ```bash
   tail -f /var/log/mdblog.log
   ```
4. **查看响应**：应该看到"同步成功"消息

## 故障排除

### 常见问题

1. **签名验证失败**

   ```
   Webhook 错误: Token 验证失败
   ```

   **解决**：启用开发模式或获取正确的 Token

2. **Git Pull 失败**

   ```
   Git pull 失败: exit status 1
   ```

   **解决**：检查 Git 仓库权限和网络连接

3. **无响应**
   - 检查 URL 是否正确
   - 确认服务器运行状态
   - 验证防火墙设置

### 调试日志

查看详细的 Webhook 调试信息：

```bash
# 查看实时日志
tail -f /var/log/mdblog.log | grep -i webhook

# 查看Token信息
tail -f /var/log/mdblog.log | grep -i "X-Gitee-Token"
```

## 参考资料

- [GitHub Webhook 文档](https://docs.github.com/zh/webhooks)
- [第三方 Gitee Webhook 实现](https://github.com/CloudnuY/gitee-webhook-handler)
- [Gitee 与其他平台对比](https://github.com/notzheng/Tapd-Git-Hooks)

## 注意事项

⚠️ **重要提醒**：

1. 生产环境建议使用 HTTPS 和 IP 白名单
2. 定期更换 Webhook 密钥
3. 监控异常访问和失败请求
4. 备份重要配置文件
