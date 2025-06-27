---
title: "Linux 下安装 Composer"
author: "BroQiang"
github_url: "https://broqiang.com"
created_at: 2017-02-21T02:28:49
updated_at: 2017-02-21T02:28:49
description: ""
tags: ["linux", "php"]
---

## 获取软件

- [官网下载地址](https://getcomposer.org/download)

- [中文官方](https://www.phpcomposer.com/)

## 安装

```bash
# 官方下载,这个连接不一定是最新版本，可以去上面的官网下载地址找到最新版本
wget https://getcomposer.org/download/1.6.3/composer.phar -O composer

chmod +x composer

sudo mv composer /usr/local/bin
```

## 配置国内全镜像

因为网络和墙的原因，官方的镜像经常访问不到，导致安装不成功或者很慢，所以配置国内源

现在有两个镜像，随意配置哪个都可以，看个人喜好，一个是 bootcss 组织搞的，
一个是 laravel-china 搞的，都算是比较靠谱的社区

### packagist.phpcomposer.com 镜像

[官网](https://www.phpcomposer.com)

- 全局（推荐）

```bash
composer config -g repo.packagist composer https://packagist.phpcomposer.com
```

- 修改当前项目

  **进入到项目根目录**

```bash
composer config repo.packagist composer https://packagist.phpcomposer.com
```

### packagist.laravel-china.org 镜像

[官网](https://laravel-china.org/composer)

> 我一般将这个镜像作为备份，只有当上面的无法使用的时候才会使用此镜像，这两个镜像选择一个即可。

- 全局（推荐）

```bash
composer config -g repo.packagist composer packagist.laravel-china.org
```

- 修改当前项目

  **进入到项目根目录**

```bash
composer config repo.packagist composer packagist.laravel-china.org
```
