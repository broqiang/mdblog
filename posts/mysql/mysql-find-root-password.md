---
title: "Mysql root 密码找回"
author: "BroQiang"
created_at: 2017-02-12T02:33:32
updated_at: 2017-02-12T02:33:32
---

> MySQL 从 5.7.6 开始，user 表结构发生了变化，重置密码的方式不太相同，请根据 MySQL 的版本来操作。

## 查询 MySQL 版本

```bash
mysql --version
```

## 停止 mysql 服务

根据 MySQL 的安装方式来停止

### init 方式

```bash
sudo /etc/init.d/mysqld stop
```

### systemd 方式

```bash
sudo systemctl stop mysqld
```

## 跳过权限认证启动

### 方法一

```bash
mysqld_safe --skip-grant-tables --skip-networking &
```

如果是 apt 方式安装的 MySQL 请使用下面方式，这种方式好像不太好用，也没去纠结过为什么。

### 方法二

修改 my.cnf 配置文件

```bash
# 根据配置文件的实际位置去修改

# 一般在这个位置
sudo vi /etc/my.cnf

# apt 方式安装的在下面位置
sudo vi /etc/mysql/mysql.conf.d/mysqld.cnf

# 自己编辑安装的话看初始化时的配置
```

在[mysqld]的段中加上一句：`skip-grant-tables`,如下：

```bash
[mysqld]
#
# * Basic Settings
#
skip-grant-tables
user        = mysql
pid-file    = /var/run/mysqld/mysqld.pid
socket      = /var/run/mysqld/mysqld.sock
```

## 重置密码

此时 MySQL 的登录可以不需要密码了，直接不输入密码登录即可。

```bash
mysql -uroot
```

登录后重置密码即可。

### 5.7.6 之前版本

```sql
mysql> UPDATE user SET Password = password ( 'new-password' ) WHERE User = 'root' ;
mysql> flush privileges;
```

### 5.7.6 及以后版本

设置成空密码，稍后就可以不使用密码登陆，然后重置密码。

```sql
UPDATE mysql.user SET authentication_string="" WHERE user="root" and host = "localhost";
```

## 重新启动 MySQL 服务

将之前配置文件中加的 `skip-grant-tables` 去掉，正常启动服务。

此时再登录 MySQL 就可以用重置后的密码去登录了。

## 再次修改账号

使用新的密码登录 MySQL 后查看下权限是否是完整的，有些版本（忘了是哪个版本，遇到过）修改完的账号功能还是不完全了，还需要再次修改，用刚刚重置完的密码登录，如果没有问题可以忽略此步。

```bash
mysql -uroot -p`
```

执行下面 sql

```sql
mysql> alter user 'root'@'localhost' identified by 'password';
```
