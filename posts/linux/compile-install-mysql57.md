---
title: "Linux 编译安装 Mysql5.7"
author: "BroQiang"
github_url: "https://broqiang.com"
created_at: 2017-04-18T02:26:09
updated_at: 2017-04-18T02:26:09
description: "记录了 Linux（Ubuntu、CentOS、Fedora）下编译安装 MySQL 5.7，当前已经更新到 5.7.21"
tags: ["linux", "mysql"]
---

此文档支持下面系统，其他发行版本未测试过

- Fedora26
- Fedora27
- CentOS7
- Ubuntu 16.04
- Ubuntu 16.10
- Ubuntu 17.10

## 安装前准备

### 获取 Mysql

当前的最新版本是 5.7.21 ,有了集成 boost 的版本, 这个比较友善 5.7.17 的时候还要单独去下载 80M 左右的 boost

```bash
# 下载源码包
wget https://dev.mysql.com/get/Downloads/MySQL-5.7/mysql-boost-5.7.21.tar.gz

# 解压到src
sudo tar xzvf mysql-boost-5.7.21.tar.gz -C /usr/local/src/

# 修改权限
sudo chown bro:bro /usr/local/src/mysql-5.7.21

```

### 安装依赖关系

```bash
# 安装开发工具包, 一般默认就会已经安装了
sudo apt install -y build-essential cmake libncurses5-dev bison

# CentoOS 用下面方式安装依赖关系
# sudo yum install -y gcc gcc-c++ cmake bison bison-devel ncurses ncurses-devel autoconf

# Fedora
# sudo dnf install -y gcc gcc-c++ cmake bison bison-devel ncurses ncurses-devel autoconf
```

### 创建守护进程用户

```bash
sudo useradd -M -s /sbin/nologin -r mysql
```

### 关闭 SELinux

> 需要注意，如果是 Fedora 或者 CentOS 系统一定要将这个关闭

```bash
# 永久关闭，需要重启才可以生效
sudo sed -i 's/SELINUX=enforcing/SELINUX=disabled/g' /etc/selinux/config
# 临时关闭
setenforce 0
```

## 编译安装

此处只配置了常规参数，更多参数请看 [官方参数说明](https://dev.mysql.com/doc/refman/5.7/en/source-configuration-options.html)

```bash
# 进入到之前解压的源码包目录
cd /usr/local/src/mysql-5.7.21/

# 创建编译后的代码保存位置
sudo mkdir build
cd build

# cmake 创建make需要的文件（makefile等）
# 参数说明
# -DCMAKE_INSTALL_PREFIX= 指定安装目录
# -DEFAULT_CHARSET= 指定默认字符集，如果不设置，安装完成后也可以配置
sudo cmake .. -DCMAKE_INSTALL_PREFIX=/usr/local/mysql -DDEFAULT_CHARSET=utf8mb4 \
-DWITH_BOOST=../boost/boost_1_59_0 -DMYSQL_DATADIR=/data/mysql/data -DWITH_SYSTEMD=1

# 编译，小内存云主机会出问题，看后面的处理办法
sudo make

# 如果配置高可以添加 -j number 参数，增加编译速度
# 如： CPU 是8核的，可以加上 make -j 6

# 安装到 /usr/local/mysql
sudo make install
```

## 配置环境变量

非必须，不过不配置的话以后就都要绝对路径执行 mysql 命令……

可以配置全局, 也可以配置当前用户, 我一般配置全局, 因为习惯平时使用普通用户, 会和 root 用户 sudo 或者 su

```bash
# 在 /etc/profile.d 下创建一个mysql用的配置文件
# 建议不要直接写在/etc/profile 中, 功能上没有区别, 不过配置多了之后可读性差点
sudo vim /etc/profile.d/mysql.sh

# 写入下面内容
export MYSQL_PATH=/usr/local/mysql
export PATH=$PATH:$MYSQL_PATH/bin

# 使配置生效
source /etc/profile.d/mysql.sh
```

## 编辑配置文件

可以直接放在`/etc/my.cnf`, 也可以 `/etc/mysql/my.cnf`

此处选择的第三种方式,个人比较喜欢所有和 mysql 相关的都放在一个目录下, 方便查找和备份

```bash
# 在mysql根目录下新建一个etc目录
sudo mkdir /usr/local/mysql/etc

# 创建my.cnf文件
sudo vim /usr/local/mysql/etc/my.cnf

# 如果是CentOS7, 需要将 /etc/my.cnf 删除，或者将配置文件直接写到这个文件中
# 如果直接写在这个文件中，后面初始化的时候也要指定配置文件为这个位置，或者就不指定，默认会到 /etc/my.cnf 去找

# 写入下面内容
[client]

[mysqld]
basedir=/usr/local/mysql
datadir=/data/mysql/data
socket=/tmp/mysql.sock
character_set_filesystem = utf8mb4
character_set_server = utf8mb4
symbolic-links=0

log-error=/data/mysql/log/mysqld.log
pid-file=/data/mysql/run/mysqld.pid
```

## 初始化数据库

需要注意:

- MySQL 5.7.6 之前, 可以使用 mysql_install_db 初始化数据库

- MySQL 5.7.6 之后, 就需要使用 mysqld with the --initialize or --initialize-insecure 初始化

```bash
# 创建数据仓库目录,日志目录,PID 目录
sudo mkdir -p /data/mysql/{data,log,run}

# 讲目录拥有者改成 mysql 用户
sudo chown -R mysql:mysql /data/mysql

# 初始化, 需要注意,配置文件及目录一定要正确
sudo /usr/local/mysql/bin/mysqld --defaults-file=/usr/local/mysql/etc/my.cnf  --initialize --user=mysql

```

## 配置 Mysql 自启动服务

```bash
# 复制启动脚本到system目录,CentOS 和 Ubuntu 编译后这个文件所在目录不太一样
# CentOS 执行下面这行命令
# sudo cp /usr/local/mysql/usr/lib/systemd/system/mysqld.service /lib/systemd/system
# Ubuntu 执行下面这行命令
sudo cp /usr/local/mysql/lib/systemd/system/mysqld.service /lib/systemd/system

# 修改启动脚本，需要注意 mysqld.pid 的位置要和配置文件中指定的相同
sudo sed -i 's/\/var\/run\/mysqld\/mysqld.pid/\/data\/mysql\/run\/mysqld.pid/g' \
/lib/systemd/system/mysqld.service

# 启动服务
sudo systemctl start mysqld

# 如果是CentOS，默认会存在一个 mariadb 的配置文件，将它删除即可，删除后再启动就不会出现问题
# sudo rm -rf /etc/my.cnf

# 停止服务
sudo systemctl stop mysqld

# 配置开机自动启动
sudo systemctl enable mysqld
```

## 配置 root@localhost 用户

Mysql 5.7 开始,初始化数据库之后不再是空的 root 密码,而是在日志文件中写入一个随机密码

安装完之后需要想默认给的密码修改

```bash
# 查询默认密码, root@localhost: 后面是密码
sudo grep 'root'@'localhost' /data/mysql/log/mysqld.log

# 用查询到的用户登录密码 mysql
/usr/local/mysql/bin/mysql -uroot -p

mysql> ALTER USER root@localhost IDENTIFIED BY '1';

# 如果是5.7.6以前版本是下面方式
# mysql> SET PASSWORD FOR user = PASSWORD('new_password');

# 测试,执行下面语句,可以查询出用户,配置正确
mysql> select user,host from mysql.user;

# 顺便查询一下字符集, 看看配置文件中的配置是否生效
mysql> show variables like '%char%';
# 正常应该显示下面结果

+--------------------------+----------------------------------+
| Variable_name            | Value                            |
+--------------------------+----------------------------------+
| character_set_client     | utf8                             |
| character_set_connection | utf8                             |
| character_set_database   | utf8mb4                          |
| character_set_filesystem | utf8                             |
| character_set_results    | utf8                             |
| character_set_server     | utf8mb4                          |
| character_set_system     | utf8                             |
| character_sets_dir       | /usr/local/mysql/share/charsets/ |
+--------------------------+----------------------------------+
8 rows in set (0.01 sec)

```

## 错误的处理方式

### 阿里云小内存服务器编译 45% 左右时候报错

错误：

```bash
[ 45%] Building CXX object sql/CMakeFiles/sql.dir/item_geofunc.cc.o
c++: internal compiler error: Killed (program cc1plus)
Please submit a full bug report,
...
```

云主机是 1G 内存的，编译到 45%的时候会报错，因为内存不足(阿里云默认没有启用虚拟内存)，

编译 `sql.dir` 比较消耗内存，正常服务器一般不会配置 1G 内存，都会够用

处理方法，自己设置虚拟内存，我的云主机是 1G 的内存，所以我配置 2G 的虚拟内存

```bash
# 创建一个 2G 的虚拟文件
sudo dd if=/dev/zero of=/swapfile bs=1k count=2048000
# 创建swap文件
sudo mkswap /swapfile
# 激活
sudo swapon /swapfile
# 查看是否配置成功
sudo swapon -s
free -m
# 如果需要下次开机自动执行，将配置配置到fstab中，此处只是为了编译临时配置，就不添加了
```

配置完虚拟内存后观察了下，用了大概 800M 的虚拟内存，所以服务器如果是 2G 的话应该是不会出现这个问题

以前虚拟机和服务器配置的时候都是 4G 起步，最近给人做项目，碰到了个奇葩的 1 核心 1G 阿里云主机，才出现了这个问题……

不过这样处理后编译的速度会慢一些，硬件配置就这么多的时候也没办法了，享受慢速吧
