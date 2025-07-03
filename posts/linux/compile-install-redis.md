---
title: "Linux 编译安装 Redis"
author: "BroQiang"
created_at: 2017-05-20T02:15:46
updated_at: 2017-05-20T02:15:46
---

[官方下载地址](https://redis.io/download)

## 安装

下面介绍的是手动安装，当然也可以通过 apt 或者 yum 安装，就不需要看本文档了

```bash
# 下载安装包
wget https://download.redis.io/releases/redis-3.2.8.tar.gz

# 解压到 src 目录
sudo tar xzvf redis-3.2.8.tar.gz -C /usr/local/src/

# 进入到解压后的目录
cd /usr/local/src/redis-3.2.8/

# 编译安装 redis 到 /usr/local/redis 目录
sudo make PREFIX=/usr/local/redis install
```

## 配置启动脚本

```bash
cd /usr/local/src/redis-3.2.8/utils/

sudo ./install_server.sh
```

### 开始配置启动脚本

下面的内容是一个测试的内容，具体的参数，要根据实际情况去配置

```bash
Welcome to the redis service installer
This script will help you easily set up a running redis server

# 填入端口
Please select the redis port for this instance: [6379] 6379
# 填写配置文件位置
Please select the redis config file name [/etc/redis/6379.conf] /usr/local/redis/etc/redis.conf
#  配置日志文件保存位置
Please select the redis log file name [/var/log/redis_6379.log] /usr/local/redis/log/redis.log
# 配置库的位置
Please select the data directory for this instance [/var/lib/redis/6379] /data/redis/redis
# 配置redis-server的位置
Please select the redis executable path [] /usr/local/redis/bin/redis-server
Selected config:
Port           : 6379
Config file    : /usr/local/redis/etc/redis.conf
Log file       : /usr/local/redis/log/redis.log
Data dir       : /data/redis/redis
Executable     : /usr/local/redis/bin/redis-server
Cli Executable : /usr/local/redis/bin/redis-cli
Is this ok? Then press ENTER to go on or Ctrl-C to abort.

```

### 启动服务

```bash
sudo /etc/init.d/redis_6379 start

# 配置开机自动启动
sudo systemctl enable redis_6379
```

## 配置环境变量

```bash
# 创建一个 redis 用的profile文件
sudo vim /etc/profile.d/redis.sh

# 写入下面内容
export REDIS_HOME=/usr/local/redis
export PATH=$PATH:$REDIS_HOME/bin

# 使配置在当前shell下生效
source /etc/profile.d/redis.sh
```

## 简单的配置

此处只介绍几个简单的配置，更多详细配置请去官网查看，具体的参数调整需要配合业务进行优化

打开我们刚刚生成的配置文件

```bash
sudo vim /usr/local/redis/etc/redis.conf
```

可以看到里面有一些默认的配置信息，下面简单介绍几个

```bash
# 默认只有本机可以访问
# 如果想要任何地方都可以访问, 将 127.0.0.1 改成 0.0.0.0
# 如果只允许私网网段访问，可以设置成 172.16.136.222 (本机的私网地址)
bind 127.0.0.1

# 监听端口，默认是6379，根据实际需要去配置
port 6379

# 客户端空闲多少秒后关闭连接，0是不关闭
timeout 0

# 是否为长链接，如果值为0，表示禁用
# 300 表示间隔300秒去和客户端确认连接
# 300是官方的默认值，一般这个值就可以，不用修改
tcp-keepalive 300

# 默认是否当做一个守护进程执行(也就是当做一个服务去执行)
daemonize yes
```

## 编写启动脚本

可以给 redis 配置一个 service 脚本，这样就可以通过 systemd 来管理进程

```bash
#表示基础信息
[Unit]
#描述
Description=Redis
#在哪个服务之后启动
After=syslog.target network.target remote-fs.target nss-lookup.target

#表示服务信息
[Service]
Type=forking
#注意：需要和redis.conf配置文件中的信息一致
PIDFile=/var/run/redis_6379.pid
#启动服务的命令
#redis-server安装的路径 和 redis.conf配置文件的路径
ExecStart=/usr/local/redis/bin/redis-server /usr/local/redis/etc/6379.conf
#重新加载命令
ExecReload=/bin/kill -s HUP $MAINPID
#停止服务的命令
ExecStop=/bin/kill -s QUIT $MAINPID
PrivateTmp=true

#安装相关信息
[Install]
#以哪种方式启动
WantedBy=multi-user.target
#multi-user.target表明当系统以多用户方式（默认的运行级别）启动时，这个服务需要被自动运行。
```

## 简单测试

如果环境变量生效了，就可以直接用 cli 去执行，终端下输入下面命令

```bash
redis-cli
```

会出现 redis 的命令行窗口，在此处写入命令即可，如下面

```bash
127.0.0.1:6379> set key1 value1
OK
127.0.0.1:6379> get key1
"value1"
127.0.0.1:6379> set test1 'This is test value'
OK
127.0.0.1:6379> get test1
"This is test value"
```
