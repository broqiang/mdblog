# Makefile for mdblog
# 
# 使用方法：
#   make dev   - 开发模式运行
#   make build - 交叉编译 Linux 可执行文件
#   make scp   - 编译并上传到服务器，重启服务并检测启动状态
#
# 配置服务器信息：
#   请修改下面的 SERVER_HOST、SERVER_USER、SERVER_PATH 变量
#   示例：
#     SERVER_HOST=192.168.1.100
#     SERVER_USER=root
#     SERVER_PATH=/bro/mdblog

# 构建变量
BINARY_NAME=mdblog
BUILD_DIR=build
MAIN_PATH=cmd/mdblog

# 服务器配置（请根据实际情况修改）
SERVER_HOST=your-server-host
SERVER_USER=your-username
SERVER_PATH=/bro/mdblog
SERVER_PORT=8080

# Go 变量
GOCMD=go

.PHONY: dev build scp

# 开发模式运行
dev:
	@echo "开发模式运行..."
	$(GOCMD) run ./$(MAIN_PATH)

# 交叉编译 Linux 可执行文件
build:
	@echo "交叉编译 Linux 可执行文件..."
	@mkdir -p $(BUILD_DIR)
	GOOS=linux GOARCH=amd64 $(GOCMD) build -o $(BUILD_DIR)/$(BINARY_NAME) ./$(MAIN_PATH)
	@echo "编译完成: $(BUILD_DIR)/$(BINARY_NAME)"

# 编译并上传到服务器，重启服务并检测
scp: build
	@echo "上传文件到服务器..."
	scp $(BUILD_DIR)/$(BINARY_NAME) $(SERVER_USER)@$(SERVER_HOST):$(SERVER_PATH)/
	@echo "上传完成"
	@echo "重启服务..."
	ssh $(SERVER_USER)@$(SERVER_HOST) "sudo systemctl restart mdblog"
	@echo "等待服务启动..."
	@sleep 5
	@echo "检测服务状态..."
	@if ssh $(SERVER_USER)@$(SERVER_HOST) "nc -z localhost $(SERVER_PORT)"; then \
		echo "✅ 服务启动成功，端口 $(SERVER_PORT) 可访问"; \
		ssh $(SERVER_USER)@$(SERVER_HOST) "sudo systemctl status mdblog --no-pager"; \
		echo "查看最新日志:"; \
		ssh $(SERVER_USER)@$(SERVER_HOST) "tail -5 $(SERVER_PATH)/logs/mdblog.log"; \
	else \
		echo "❌ 服务启动失败，端口 $(SERVER_PORT) 不可访问"; \
		echo "查看服务状态:"; \
		ssh $(SERVER_USER)@$(SERVER_HOST) "sudo systemctl status mdblog --no-pager"; \
		echo "查看错误日志:"; \
		ssh $(SERVER_USER)@$(SERVER_HOST) "tail -10 $(SERVER_PATH)/logs/mdblog.log"; \
		exit 1; \
	fi 