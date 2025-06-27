---
title: "CentOS 下 Nginx 配置 Let's Encrypt 证书"
author: "BroQiang"
github_url: "https://broqiang.com"
created_at: 2018-03-15T01:38:12
updated_at: 2018-03-15T01:38:12
description: "记录了 CentOS 7 下 Let's Encrypt 单域名及通配符域名的申请已经 Nginx 配置 https 服务"
tags: ["linux", "nginx"]
---

体验了下免费的 ssl 证书，看起来很好用。

效果：

![示例图](https://image.broqiang.com/qJ8d953bvu0fwOsnYTnK6OJ2WE4RrP.png)

[官网](https://letsencrypt.org/)

安装要求在 shell 环境下，用官方提供的 ACEI 工具来安装，官方只提供了 Linux 下的安装环境，正好服务器就是 CentOS 7，就直接在服务器上安装 证书。

[官方安装文档](https://certbot.eff.org/lets-encrypt/centosrhel7-nginx)

官方提供的 RPM 包是有 bug 的，使用的时候会报错，所以这里采用 pip 安装

pip 是 python 的一个在线安装工具，因为 certbot 是 Python 开发的，所以也可以用 pip 在线安装

## 安装

```bash
# 安装 Python 的的开发环境
yum install openssl-devel python-devel python2-pip

# 升级 pip
pip install --upgrade pip setuptools

# 安装 证书配置工具 certbot 及依赖关系
# pm 包就是下面的依赖关系有问题，等未来 rpm 包更新之后还是可以继续按照官方的方式安装
pip install certbot
pip install requests urllib3 pyOpenSSL --force --upgrade
```

## 申请证书

### 准备

以前（2017-3-14 之前），只支持单独域名证书，14 日之后发布了 V2 版本，可以支持通配符域名证书了。

#### 查询版本

```bash
certbot --version
```

结果：

```bash
certbot 0.22.0
```

可以看到，这个版本就可以支持通配符域名了

#### 开始申请证书

```bash
sudo certbot certonly  -d *.broqiang.com --manual --preferred-challenges dns --server https://acme-v02.api.letsencrypt.org/directory -m broqiang@qq.com
```

- `-d` 参数是用来指定域名的

这里有个坑，_.broqiang.com 和 broqiang.com 不是一样的，我第一次申请完证书发现不信任，提示一直是 broqiang.com 的证书是由 _.broqiang.com 发布，不被信任， 重复申请了好多次才发现。

- `--preferred-challenges dns` 验证域名的方式，需要配置一个 txt 记录，才确保是域名的拥有者

- `--server https://acme-v02.api.letsencrypt.org/directory` 指定服务器

因为 V2 版本 和 V1 版本不是相同的服务器，因此要是申请通配符域名需要在这里指定服务器。

`-m broqiang@qq.com` 指定申请者邮箱，要填写有效的邮箱，因为需要邮箱点击一个确定链接，并且未来域名快到期，也会发这个邮箱来通知。

执行完上面的命令会有一些交互式的操作，按照提示一步一步完成即可。

```bash
Saving debug log to /var/log/letsencrypt/letsencrypt.log
Plugins selected: Authenticator manual, Installer None
Obtaining a new certificate
Performing the following challenges:
dns-01 challenge for broqiang.com

-------------------------------------------------------------------------------
NOTE: The IP of this machine will be publicly logged as having requested this
certificate. If you're running certbot in manual mode on a machine that is not
your server, please ensure you're okay with that.

Are you OK with your IP being logged?
-------------------------------------------------------------------------------
(Y)es/(N)o: yes
```

这里提示是否将域名 IP 绑定，点击 y 或者 yes

```bash
-------------------------------------------------------------------------------
Please deploy a DNS TXT record under the name
_acme-challenge.broqiang.com with the following value:

ruLiZ00NxCovX3vRou5I3qdMQ0QGvOPJ5x3IWJvKI_I

Before continuing, verify the record is deployed.
-------------------------------------------------------------------------------
Press Enter to Continue
```

这里提示你的域名接卸一个 txt 记录，然后验证。

如上，就要配置一个 key 为 `_acme-challeng` ，value 为：`ruLiZ00NxCovX3vRou5I3qdMQ0QGvOPJ5x3IWJvKI_I` 的 txt 记录。

配置好了回车确定，这里不要急着回车，需要确定解析已经生效。

```bash
# 确认生效的办法
dig _acme-challenge.broqiang.com txt

# 如果提示没有这个命令可以用下面命令安装
sudo yum install -y bind-utils
```

执行完上面命令将返回的值和配置的值对比。

```bash
Waiting for verification...
Cleaning up challenges

IMPORTANT NOTES:
 - Congratulations! Your certificate and chain have been saved at:
   /etc/letsencrypt/live/broqiang.com/fullchain.pem
   Your key file has been saved at:
   /etc/letsencrypt/live/broqiang.com/privkey.pem
   Your cert will expire on 2018-06-13. To obtain a new or tweaked
   version of this certificate in the future, simply run certbot
   again. To non-interactively renew *all* of your certificates, run
   "certbot renew"
 - If you like Certbot, please consider supporting our work by:

   Donating to ISRG / Let's Encrypt:   https://letsencrypt.org/donate
   Donating to EFF:                    https://eff.org/donate-le
```

验证通过后会生成证书。

证书保存在 `/etc/letsencrypt/archive/broqiang.com` 但是他会在 `/etc/letsencrypt/live/broqiang.com` 中创建证书的软链接，推荐使用 live 中的，因为它会自动将最新的版本链接到 live 中。

通配符的域名就已经申请完成，不过需要注意，此时只可以用 \*.broqiang.com 不包括主域名 broqiang.com ，如果主域名需要的话，可以再申请一份证书，只需要将前面的命令的 -d 参数修改即可。

```bash
sudo certbot certonly  -d broqiang.com --manual --preferred-challenges dns --server https://acme-v02.api.letsencrypt.org/directory -m broqiang@qq.com
```

## Nginx 中使用证书

```nginx
# 将 http 的访问自动 301 重定向到 https 上
server {
    listen       80;
    server_name  broqiang.com www.broqiang.com;

    return 301 https://broqiang.com$request_uri;
}


# 项目使用的虚拟主机
server {
    listen       443;
    ssl on;
    ssl_certificate /etc/letsencrypt/live/broqiang.com-0001/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/broqiang.com-0001/privkey.pem;
    ssl_trusted_certificate /etc/letsencrypt/live/broqiang.com-0001/chain.pem;

    server_name  broqiang.com www.broqiang.com;
    index index.php index.html;
    root  /www/web/www.broqiang.com/public;

   # 301 重定向
   if ($host = 'www.broqiang.com'){
        rewrite ^/(.*)$ https://broqiang.com/$1 permanent;
   }

    location / {
        try_files $uri $uri/ /index.php?$query_string;
    }

    location ~ \.php {
        fastcgi_pass    127.0.0.1:9000;
        fastcgi_split_path_info ^(.+\.php)(.*)$;
        fastcgi_param PATH_INFO $fastcgi_path_info;
        fastcgi_param PATH_TRANSLATED $document_root$fastcgi_path_info;
        fastcgi_param SCRIPT_FILENAME $document_root$fastcgi_script_name;
        include         fastcgi_params;
    }

    fastcgi_intercept_errors on;
    access_log /www/webLogs/www.broqiang.com_access.log;
    #access_log off;
}

# 禁止 IP 直接访问(http)，也就是说没有绑定虚拟主机的域名也不能访问
server
{
    listen 80 default;

    #server_name _;

    return 500;

}

# 禁止 IP 直接访问(https)
# 需要注意，这里也要配置证书，否则会先报个证书不可信，点击确定好再跳转到 500，非常不友好。
server
{
    listen 443 ssl default_server;
    ssl_certificate /etc/letsencrypt/live/broqiang.com/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/broqiang.com/privkey.pem;

    return 500;

}
```

配置完成后重启 Nginx 服务即可，具体的 域名、路径等根据实际的项目去配置。

## 配置自动更新证书

因为证书的有效期是 3 个月，又懒得去计算有效期了，所以写个计划任务，每天凌晨自动去更新一次证书。

```bash
# 要用 root 用户去执行更新证书的动作，因为之前获取证书的时候就是 root 账号做的。
sudo crontab -e
```

写入下面内容

```bash
02 00 * * * root certbot renew
```

确定是否配置成功

```bash
sudo crontab -l
```
