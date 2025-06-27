---
title: "Git 安装"
author: "BroQiang"
github_url: "https://broqiang.com"
created_at: 2017-02-18T02:31:03
updated_at: 2017-02-18T02:31:03
description: "记录了 Ubuntu、Fedora、CentOS 及源码编译安装，并做了安装后的简单配置。"
tags: ["linux", "git"]
---

## 获取最新安装包

- [官方下载列表](https://www.kernel.org/pub/software/scm/git/)

- [当前最新版本 2.19.1](https://mirrors.edge.kernel.org/pub/software/scm/git/git-2.19.1.tar.gz)

- [Github 仓库](https://github.com/git/git)

- [官方中文文档](https://git-scm.com/book/zh/v2)

## Ubuntu 安装

Ubuntu 软件仓库自带的版本是比较新的，可以直接在线安装。

```bash
sudo apt install git -y
```

## Fedora 安装

Fedora 软件仓库的版本也是很新的，比 Ubuntu 还要新一些。

```bash
sudo dnf install git
```

## CentOS 安装

CentOS 也可以通过 yum 安装，不过 yum 安装的版本有些低，`不推荐`使用此方法，和最新版本差的有些多，会有一些不兼容。

```bash
sudo yum install git
```

## 源码编译安装

如果在线安装的版本过低或者想用最新版本，可以直接编译 git 源码安装，`推荐 CentOS 使用这种方式` 。

下面的方法支持 CentOS 和 Ubuntu，区别只是依赖关系的安装方式不同，其他全部相同

### 安装依赖关系

```bash
# CentOS
sudo yum install curl-devel expat-devel gettext-devel openssl-devel zlib-devel  perl-ExtUtils-MakeMaker autoconf gcc gcc-c++

# Fedora
sudo dnf install curl-devel expat-devel gettext-devel openssl-devel zlib-devel  perl-ExtUtils-MakeMaker

# Ubuntu
sudo apt install libcurl4-gnutls-dev libexpat1-dev gettext libz-dev libssl-dev
```

### 安装编译 doc 和 info 的依赖关系

这个是用来编译的时候将 man 和 info 手册的内容编译出来，个人不建议编译，因为下面几个依赖安装的包比较多，一般 git 的 man 也用不到，所以个人推荐 `不安装` 下面依赖，后面编译的时候也`不去编译` doc 和 info

```bash
# CentOS
sudo yum install asciidoc xmlto docbook2X
# CentOS 有个 bug，docbook2X 安装完成后的包名叫 db2x_docbook2texi ，所以需要创建一个 软链接
sudo ln -s /usr/bin/db2x_docbook2texi /usr/bin/docbook2x-texi

# Fedota
sudo dnf install asciidoc xmlto docbook2X

# Ubuntu
sudo apt install asciidoc docbook2x
```

### 下载并解压

```bash
# 下载解压
$ wget https://mirrors.edge.kernel.org/pub/software/scm/git/git-2.19.1.tar.gz

$ sudo tar xzvf git-2.19.1.tar.gz -C /usr/local/src/
$ cd /usr/local/src/git-2.19.1
```

### 配置并编译

```bash
# 如果是 github 上下载的版本，需要执行下面命令，官方下载的可以忽略
sudo make configure

# 配置
sudo ./configure --prefix=/usr/local/git

# 编译，一般情况下直接执行 sudo make all 就可以了
# 如果需要编译 man 和 info 帮助手册，执行 make all doc info
sudo make all

# 安装，一般直接安装就可以了
# 如果前面 make 的时候编译了 doc 和 info
# 也可以执行 sudo make install install-doc install-html install-info 安装man 和 info 的手册
sudo make install
```

## 配置环境变量

```bash
sudo vim /etc/profile.d/git.sh
# 写入下面两行到 git.sh
export GIT_HOME=/usr/local/git
export PATH=$GIT_HOME/bin:$PATH

# 生效下
source /etc/profile.d/git.sh

# 测试下
git --version
```

## 配置账号

安装完成后最好配置下用户和邮箱，否则 commit 的时候也无法提交

全局配置

```bash
# 邮箱改成实际的
$ git config --global user.email "broqiang@qq.com"

# 用户名改成实际的
$ git config --global user.name "Bro Qiang"

```

## 配置记住用户名密码

比如提交 github ，如果不记住，每次都要输一次，个人电脑，还是记住比较方便

```bash
# 长期保存, 一般要是个人开发电脑，配置这个就行
$ git config --global credential.helper store

# 临时保存，默认15分钟
$ git config --global credential.helper cache

# 指定临时保存时间
$ git config credential.helper 'cache --timeout=3600'
```

## 别名设置

因为命令会很长，提交频繁的时候很影响工作效率，因此利用 Linux 的别名特性，给 Git 创建一些别名，可以根据个人的喜好创建别名。

```bash
$ vim ~/.bashrc
# 在默认配置文件中写入下面内容
# 为了尽量不破坏系统原本配置文件的结构，新建了一个 .bash_alias文件，也可以将别名直接写在下面，效果没有任何区别
# 为了防止配置错误启动的时候报错，做了个判断
# 如果家目录下有.bash_alias文件将bash_alias文件引入, .bash_alias可以根据个人喜好随意起名字
# 为了尽量不破坏系统原本配置文件的结构，
if [ -f ~/.bash_alias ]; then
    . ~/.bash_alias
fi

# 然后新建 ~/.bash_alias文件
$ vim ~/.bash_alias
#写入下面内容，此处只是我常用的配置，可以根据个人习惯去配置
# git
alias gs='git status'
alias gaa='git add .'
alias ga='git add '
alias gp='git push'
alias gc='git commit -m '
alias gl='git log'
alias grao='git remote add origin '
alias gpo='git push origin '
```

## 更新日志

### 2018-10-7

更新版本到 git-2.19.1
