# Makefile for mdblog
# 
# 使用方法：
#   make dev   - 开发模式运行
#   make build - 交叉编译 Linux 可执行文件  
#   make scp   - 停止服务、编译并上传到服务器、重启服务并检测启动状态
#
# 服务器配置说明：
#   SERVER_HOST: 服务器IP地址
#   SERVER_USER: SSH用户名
#   SERVER_PATH: 服务器上的部署路径
#   SSH_PORT: SSH连接端口
#   APP_PORT: 应用监听端口
#
# 注意：
#   - 服务器地址和端口已硬编码在程序中 (0.0.0.0:8091)
#   - 静态资源已嵌入到可执行文件中，无需单独上传web目录
#   - 如需修改应用配置，请编辑 internal/config/config.go 并重新编译

# 构建变量
BINARY_NAME=mdblog
BUILD_DIR=build
MAIN_PATH=.

# 服务器配置
SERVER_HOST=123.56.186.148
SERVER_USER=root
SERVER_PATH=/bro/mdblog
SSH_PORT=21345
APP_PORT=8091

# Go 变量
GOCMD=go

.PHONY: dev build scp

# 开发模式运行
dev:
	@echo "开发模式运行..."
	$(GOCMD) run main.go

# 交叉编译 Linux 可执行文件
build:
	@echo "交叉编译 Linux 可执行文件..."
	@mkdir -p $(BUILD_DIR)
	GOOS=linux GOARCH=amd64 $(GOCMD) build -o $(BUILD_DIR)/$(BINARY_NAME) main.go
	@echo "编译完成: $(BUILD_DIR)/$(BINARY_NAME)"

# 编译、停止服务、上传并重启服务
scp: build
	@echo "停止服务..."
	ssh -p $(SSH_PORT) $(SERVER_USER)@$(SERVER_HOST) "sudo systemctl stop mdblog"
	@echo "上传文件到服务器..."
	scp -P $(SSH_PORT) build/mdblog $(SERVER_USER)@$(SERVER_HOST):$(SERVER_PATH)/
	@echo "重启服务..."
	ssh -p $(SSH_PORT) $(SERVER_USER)@$(SERVER_HOST) "sudo systemctl start mdblog"
	@echo "等待服务启动..."
	@sleep 5
	@echo "检测服务状态..."
	@if ssh -p $(SSH_PORT) $(SERVER_USER)@$(SERVER_HOST) "nc -z localhost $(APP_PORT)" >/dev/null 2>&1; then \
		echo "✅ 服务启动成功，端口 $(APP_PORT) 可访问"; \
		ssh -p $(SSH_PORT) $(SERVER_USER)@$(SERVER_HOST) "sudo systemctl status mdblog --no-pager -l"; \
		echo ""; \
		echo "📊 最新日志:"; \
		ssh -p $(SSH_PORT) $(SERVER_USER)@$(SERVER_HOST) "tail -n 5 /bro/mdblog/logs/mdblog.log"; \
	else \
		echo "❌ 服务启动失败，端口 $(APP_PORT) 不可访问"; \
		ssh -p $(SSH_PORT) $(SERVER_USER)@$(SERVER_HOST) "sudo systemctl status mdblog --no-pager -l"; \
		echo ""; \
		echo "📊 错误日志:"; \
		ssh -p $(SSH_PORT) $(SERVER_USER)@$(SERVER_HOST) "tail -n 10 /bro/mdblog/logs/mdblog.log"; \
		exit 1; \
	fi 