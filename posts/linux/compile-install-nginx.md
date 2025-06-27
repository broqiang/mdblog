---
title: "Linux 编译安装 Nginx"
author: "BroQiang"
github_url: "https://broqiang.com"
created_at: 2017-04-18T02:21:44
updated_at: 2017-04-18T02:21:44
description: "记录了当前最新稳定版的 nginx 1.12.2 的编译安装过程，并进行了 PHP 的配置。"
tags: ["linux", "nginx"]
---

当前最新稳定版本: 1.12.2

[官网](https://nginx.org/)

## 获取软件

```bash
# 下载当前最新稳定版
$ wget https://nginx.org/download/nginx-1.12.2.tar.gz

# 解压到src目录
$ sudo tar xzvf nginx-1.12.2.tar.gz -C /usr/local/src/
```

## 解决依赖关系

```bash
# Ubuntu 执行下面命令
sudo apt install -y libpcre3-dev zlib1g-dev libssl-dev

# CentOS 执行下面命令
sudo yum install -y pcre-devel openssl-devel zlib-devel

# Fedora
sudo dnf install -y pcre-devel openssl-devel zlib-devel
```

## 创建守护进程用户

如果安装 PHP 的时候已经创建了,此步骤忽略

```bash
sudo useradd -M -s /sbin/nologin www
```

## 编译安装

参数多数用了默认的, 如果有特殊需求,可以查看 [官方参数配置](https://nginx.org/en/docs/configure.html)

```bash
# 进入到解压后的源码目录
cd /usr/local/src/nginx-1.12.2/

# 配置
sudo ./configure --prefix=/usr/local/nginx --user=www --group=www \
--with-select_module --with-poll_module --with-http_ssl_module --with-pcre \
--with-pcre-jit --with-zlib= --pid-path=/usr/local/nginx/run/nginx.pid

# 如果是使用的 gcc 7，如 Fedora 26 （我就是这个版本出现的问题）
# 加上这个参数:  --with-cc-opt="-Wno-error"
# 简单解释下，查找的资料说 gcc 7 检察的更加严格，所以原本的编译就编译不过去了

# 编译 & 安装
sudo make
sudo make install
```

## 配置自动启动

```bash
# 新建 service 文件
sudo vim /lib/systemd/system/nginx.service

# 写入下面内容
[Unit]
Description=nginx - high performance web server
Documentation=https://nginx.org/en/docs/
After=network.target remote-fs.target nss-lookup.target

[Service]
Type=forking
PIDFile=/usr/local/nginx/run/nginx.pid
ExecStartPre=/usr/local/nginx/sbin/nginx -t -c /usr/local/nginx/conf/nginx.conf
ExecStart=/usr/local/nginx/sbin/nginx -c /usr/local/nginx/conf/nginx.conf
ExecReload=/bin/kill -s HUP $MAINPID
ExecStop=/bin/kill -s QUIT $MAINPID
PrivateTmp=true

[Install]
WantedBy=multi-user.target

# 配置开机自动启动
sudo systemctl enable nginx

# 启动服务
sudo systemctl start nginx

# 停止服务
sudo systemctl stop nginx
```

## 修改配置文件

### nginx.conf 主配置文件

简单配置了下,进行了一些基本的优化，将配置文件拆分开，方便分开配置

这个只是最初步的优化,真实环境要根据实际情况去优化

实际生产中各种原因都会影响性能,有时候并不是参数给的高性能就会提高

这个里面的参数要根据实际的服务器配置,并发访问情况等去优化,很难一步到位的

```bash
# 修改配置前先将原本的备份一份, 要养成良好习惯,因为谁也不能保证自己的修改就是 100% 正确, 必要时候要能够还原
cd /usr/local/nginx/conf/
sudo mv nginx.conf nginx.conf_`date +%Y-%m-%d`_bak

sudo vim /usr/local/nginx/conf/nginx.conf
```

将默认的内容全部删除,写入下面内容

```nginx
# 守护进程用户
user  www;

#启动进程,通常设置成和cpu的数量相等
worker_processes  4;

events {
    # epoll是多路复用IO(I/O Multiplexing)中的一种方式,
    # 仅用于linux2.6以上内核,可以大大提高nginx的性能
    use   epoll;

    # 单个worker process进程的最大并发链接数
    # 并发总数理论值是 worker_processes 和 worker_connections 的乘积
    worker_connections  8192;
}

include conf.d/http.main;

```

### http 模块主配置

```bash
# 创建自定义配置文件存放目录 目录名称就是上面配置文件中 include 的目录
sudo mkdir /usr/local/nginx/conf/conf.d

sudo vim /usr/local/nginx/conf/conf.d/http.main
```

写入下面内容：

```nginx
http {
    include       mime.types;
    default_type  application/octet-stream;

    sendfile on;
    directio 4m;
    directio_alignment 512;

    tcp_nopush on;
    tcp_nodelay off;

    keepalive_timeout 180s;
    keepalive_requests 100;
    keepalive_disable msie6;
    send_timeout 60s;
    client_body_timeout 60s;
    client_header_timeout 30s;

    gzip on;
    gzip_min_length 1024;
    gzip_buffers 4 16k;
    gzip_comp_level 5;
    gzip_proxied any;
    gzip_types text/plain application/json application/x-javascript text/css appl
    ication/xml;
    gzip_vary on;
    gzip_disable "MSIE [1-6]\.";

    fastcgi_buffer_size 64k;
    fastcgi_buffers 4 64k;
    fastcgi_busy_buffers_size 128k;
    fastcgi_temp_file_write_size 256k;
    fastcgi_connect_timeout 300;
    fastcgi_send_timeout 300;
    fastcgi_read_timeout 300;

    include conf.d/*.conf;
}
```

## 配置 server (虚拟主机)

```bash
# 有新的虚拟主机就再创建一个 xxx.conf 即可
$ sudo vim /usr/local/nginx/conf/conf.d/test.conf
```

写入下面内容

```nginx
server {
    listen       80; # 端口,一般http的是80
    server_name  localhost; # 一般是域名,本机就是localhost
    index index.php index.html;  # 默认可以访问的页面,按照写入的先后顺序去寻找
    root  /home/bro/test; # 项目根目录

    # 防止访问版本控制内容
    location ~ .*.(svn|git|cvs) {
        deny all;
    }

    # 此处不是必须的,需要时候配置
    location / {
        # Laravel 5.4 Url 重写
        try_files $uri $uri/ /index.php?$query_string;
    }


    # 下面是所有关于 PHP 的请求都转给 php-fpm 去处理
    location ~ \.php {

        # 注意： unix sock 和 ip ，两种方式只能选择一种

        # 基于unix sock 访问,Ubuntu Apt 方式安装的PHP默认是以sock方式启动
        # fastcgi_pass    unix:/run/php/php7.0-fpm.sock;

        # 基于IP访问
        fastcgi_pass    127.0.0.1:9000;

        fastcgi_split_path_info ^(.+\.php)(.*)$;
        fastcgi_param PATH_INFO $fastcgi_path_info;
        fastcgi_param PATH_TRANSLATED $document_root$fastcgi_path_info;
        fastcgi_param SCRIPT_FILENAME $document_root$fastcgi_script_name;
        include         fastcgi_params;
    }

    fastcgi_intercept_errors on;
    # 日志保存目录,一般按照项目单独保存, 开发环境可以关闭
    #access_log  logs/localhost_access.log access;
    access_log off;
}
```

## 其他配置

### 隐藏 Nginx 或 Tengine 的服务器及版本信息

这样可以稍微安全一点，叫功能者不能快速定位 Web 服务，修改也很简单，直接在编译前修改下源码即可（一定是编译前修改，如果已经编译完的话，可以重新编译一次，反正 Nginx 编译速度也很快）。如果是个人环境或者开发环境这个修改就省略吧，反正也不对外提供服务，不涉及到安全的问题。

修改文件 `源码包中/src/core/nginx.h`

（适用于 nginx 和 Tengine ，位置和文件名称相同）

修改里面的内容：

`Nginx` 修改方式

将这两行修改成想要显示的内容

```c
#define NGINX_VERSION      "1.12.2"
#define NGINX_VER          "nginx/" NGINX_VERSION
```

修改后：

```c
#define NGINX_VERSION      "3.4.5"
#define NGINX_VER          "MyServer/" NGINX_VERSION
```

`Tengine` 修改方式：

修改前

```c
#define TENGINE            "Tengine"
#define tengine_version    2002002
#define TENGINE_VERSION    "2.2.2"
```

修改后

```c
#define TENGINE            "MyServer"
#define tengine_version    3003003
#define TENGINE_VERSION    "3.3.3"
```
