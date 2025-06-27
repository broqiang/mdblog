---
title: "Samba 服务器配置"
author: "BroQiang"
github_url: "https://broqiang.com"
created_at: 2018-01-05T01:56:11
updated_at: 2018-01-05T01:56:11
description: "这里的示例使用的是 Ubuntu，CentOS 和 Fedora 也类似。"
tags: ["linux"]
---

## 服务器端配置

### 安装软件

```bash
# 服务器端
sudo apt install samba
```

### 创建准备共享的目录和用户

如果已经存在准备共享使用的目录和用户就不需要了

```bash
# 共享目录（可以是其他目录，根据需要去创建使用）
mkdir ~/share

# 创建用户，如果共享目录是公开的不需要验证，也可以不创建
sudo useradd bro
```

### 修改配置文件

```bash
sudo vim /etc/samba/smb.conf
```

打开配置文件，修改下面几点：

### 将不相关的内容注释

```bash
[printers]
   comment = All Printers
   browseable = no
   path = /var/spool/samba
   printable = yes
   guest ok = no
   read only = yes
   create mask = 0700

# 和下面的内容
[printers]
   comment = All Printers
   browseable = no
   path = /var/spool/samba
   printable = yes
   guest ok = no
   read only = yes
   create mask = 0700
```

改为：

```bash
#[printers]
#   comment = All Printers
#   browseable = no
#   path = /var/spool/samba
#   printable = yes
#   guest ok = no
#   read only = yes
#   create mask = 0700
#
## Windows clients look for this share name as a source of downloadable
## printer drivers
#[print$]
#   comment = Printer Drivers
#   path = /var/lib/samba/printers
#   browseable = yes
#   read only = yes
#   guest ok = no
```

### 配置自定义的共享

在配置文件最下面加入下面内容

```bash
# 分享后的目录名称，客户端看到的名字
# 如 //127.0.0.1/share
[share]

    # 说明
    comment = share
    # 共享的目录
    path = /home/bro/share
    # 除了使用者，其他人是否能浏览(需要输入 ip/share 才可以看到，直接输入IP就不会列出来了)
    browseable = yes
    # 是否可以被所有人读
    public = no
    # 是否允许客户端进行修改，no 是只读方式
    writable = no

    # 只有下面用户才可以访问共享目录
    valid users = bro,samba

```

配置完成后重启服务生效

```bash
sudo systemctl restart smbd
```

### 开放防火墙

防火墙需要开放 tcp:445 和 tcp:139 端口

## Linux 客户端使用

### 安装客户端软件

```bash
sudo apt install smbclient
```

### 命令行使用

```bash
# 如果服务器端不需要密码可以把后面的 --user=bro 省略
smbclient //127.0.0.1/FullStack02 --user=bro
```

登录后的界面如下，可以通过 help 查看支持的命令，使用方法和 ftp 类似

```bash
smb: \>
smb: \> help
?              allinfo        altname        archive        backup
blocksize      cancel         case_sensitive cd             chmod
chown          close          del            dir            du
echo           exit           get            getfacl        geteas
hardlink       help           history        iosize         lcd
link           lock           lowercase      ls             l
mask           md             mget           mkdir          more
mput           newer          notify         open           posix
posix_encrypt  posix_open     posix_mkdir    posix_rmdir    posix_unlink
posix_whoami   print          prompt         put            pwd
q              queue          quit           readlink       rd
recurse        reget          rename         reput          rm
rmdir          showacls       setea          setmode        scopy
stat           symlink        tar            tarmode        timeout
translate      unlock         volume         vuid           wdel
logon          listconnect    showconnect    tcon           tdis
tid            logoff         ..             !
smb: \>
```

### 直接挂在到目录

如果觉得上面的方式比较麻烦，可以将共享目录直接挂载到服务器上，这样就可以当普通目录来操作

```bash
# 需要先安装 cifs 工具
sudo apt install cifs-utils

sudo mount.cifs //127.0.0.1/share /mountdir

# 如果有用户名和密码的话
sudo mount.cifs //127.0.0.1/share -o username=bro%123456 /mountdir
```

## Windows 客户端使用

Windows 客户端使用就比较省事了

### 通过文件管理器查看

直接 `Win键+R` 输入`\\192.168.200.151\share`

需要注意：

- IP 换成 samba 服务器的 IP

- share 换成实际的共享名称

- 如果不需要密码直接就可以访问，如果需要的话会弹出个窗口，输入即可

### 挂载成网络文件夹

找到网络，右键选择 `映射网络驱动器` -> `选择盘符` 在文件夹处填入: `\\192.168.200.151\FullStack02`

### 清除记住的 samba 用户名密码

可能为了方便，选择了将用户名和密码记住，一般不会有问题，不过如果服务器端更改了用户名和密码，这个时候就不能登录了，就需要手动将记住的用户名和密码清除

- `Win+R` 输入 `cmd`

- `net use` 可以查看网络映射

- `net use \\192.168.200.151\FullStack02 /delete` 根据查询出来的名称，然后删除

- 也可以通过 _ 号将所有的删除 `net use _ /del`
