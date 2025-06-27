---
title: "Ubuntu 18.04 设置 ulimit"
author: "BroQiang"
github_url: "https://broqiang.com"
created_at: 2019-05-04T22:35:03
updated_at: 2019-05-04T22:35:03
description: "Ubuntu 18.04 设置 ulimit ，解决 too many open files 的问题。"
tags: ["linux"]
---

一直用 Ubuntu 做桌面，本来以为很简单就可以设置了。

## CentOS7 设置 ulimit

将这个也记录在这里就是为了做个对比，之前设置 CentOS 的时候非常简单，只需要编辑
`/etc/security/limits.conf` ，加入下面内容即可：

```bash
* soft nproc 65535
* hard nproc 65535
* soft nofile 65535
* hard nofile 65535
```

如果是 root 用户的话，这样就 OK 了， 如果是普通用户，还需要将
`/etc/security/limits.d/20-nproc.conf` 修改或注释，如：

```bash
# 我选择的是将此行注释，因为它会覆盖 limit.conf 中的配置，或者直接修改此文件也可以
# *          soft    nproc     4096
root       soft    nproc     unlimited
```

## ulimit 配置的解释

- nproc 是操作系统级别对每个用户创建的进程数的限制。

- nofile 最大打开文件描述符数, too many open files 的时候一般就是这个小了（在系统资源空闲的时候）。

- 第一列的 \* 或 root ，代表为所有用户或 root 用户设置这个配置。

- 第二列的 soft 软规则限制， hard 硬规则限制， soft 指的是当前系统生效的设置值,
  hard 表明系统中所能设定的最大值，soft 的限制不能比 hard 限制高。

- 65535 是设置的值，根据实际的需求和系统的配置去配置。

## Ubuntu 设置 ulimit

这个就有点坑（或者说对 Ubuntu 不够了解）， 一直用 Ubuntu 做的桌面，服务器用的是 CentOS ，
今天做测试的时候遇到这个问题了，按照 CentOS 的方式配置完后竟然不生效，后来搜索了下，
发现竟然比 CentOS 复杂了很多，要多修改好几个文件，所以记录下来。

参考： [https://askubuntu.com/questions/1049058/how-to-increase-max-open-files-limit-on-ubuntu-18-04](https://askubuntu.com/questions/1049058/how-to-increase-max-open-files-limit-on-ubuntu-18-04)

### 修改 limit.conf 文件

写入下面内容

```bash
* soft     nproc          65535
* hard     nproc          65535
* soft     nofile         65535
* hard     nofile         65535

root soft     nproc          65535
root hard     nproc          65535
root soft     nofile         65535
root hard     nofile         65535

bro soft     nproc          65535
bro hard     nproc          65535
bro soft     nofile         65535
bro hard     nofile         65535
```

> 注意，这里是每一个用户都要配置，网上有人说只有 root 用户才需要单独写出来，不过我试的时候发现，
> bro （普通用户） 没有写在这里的时候也没生效。

### 添加 common-session

在 `/etc/pam.d/common-session` 和 `/etc/pam.d/common-session` 文件中，都添加下面内容：

```bash
session required pam_limits.so
```

### 修改 /etc/systemd/system.conf 文件

如果不修改这个文件的话，重启后当前登录的 bro 用户不会超过这个限制（影响的是登录桌面的用户）。

修改完成后，如果是远程链接，重新链接即可生效； 如果是桌面环境，重启后生效。
