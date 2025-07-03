---
title: "SSH 自动登录脚本"
author: "BroQiang"
created_at: 2018-03-14T01:42:26
updated_at: 2018-03-14T01:42:26
---

## 密钥登录（免密码）

### 在客户端生成公钥

```bash
ssh-keygen -t rsa -P '' -f ~/.ssh/id_rsa
```

参数说明：

- `-t` 指定密钥类型，如果不指定，默认就是 `rsa`

- `-P` 指定密码，如果不指定这个参数，将会出现交互式提示，按照提示收入密码也可以

- `-f` 指定密钥文件吗，如果不指定，处出现交互式提示，按照提示输入也可以

执行完上面的命令，会在 `~/.ssh` 下生成两个文件 `id_rsa`（私钥） 和 `id_rsa.pub`（公钥），然后将 `id_rsa.pub`上传到服务器。

上传到服务器的命令：

```bash
# 上传到远端服务器的家目录下
scp -P 22 id_rsa.pub bro@www.broqiang.com:~/
```

### 服务器端配置

```bash
# 将密钥中的内容写入到 authorized_keys 文件，如果不存在会新建一个
cat ~/id_rsa.pub >> ~/.ssh/authorized_keys

# 注意修改权限，要求 ~/.ssh/authorized_keys 文件必须是 600 权限
chmod 600 ~/.ssh/authorized_keys
```

配置完成后可以测试下，会发现不需要密码就可以登录了

```bash
ssh bro@broqiang.com -p 22
```

## 自动登录脚本

此脚本可以配置多个服务器，只需要在 servers 数组中添加即可，可以有一个 key 对应一个服务器的用户名及 IP，这样就可以用一个方便记住的名字来连接服务器，不用去记住复杂的 IP 了。

```bash
#!/bin/bash
#

declare -A servers;

servers=(
    [bro]=bro@47.93.219.74
)

if [[ -z ${1} ]]; then
    echo -e "\n\033[1;31mPlease input server\033[0m\n"
    exit
fi

server=${servers[${1}]}

if [[ -z ${server} ]]; then
    echo -e "\n\033[1;31mNo server found for ${1}\033[0m\n"

    read -p "Whether to see all servers ? [Y/N]" yn

    if [[ "${yn}" == "y" || "${yn}" == "Y" ]]; then

        echo "=================================="

        for i in ${!servers[@]}
        do
            echo -e "${i} -- ${servers[$i]}"
        done

        echo
    fi

    exit
fi


port=${2:-22}

ssh ${servers[${1}]} -p ${port}
```

因为 Ubuntu 默认 shell 改成了 dash，所以不能用 sh mssh 来执行脚本，可以用 ./mssh 或者将此脚本配置到环境变量中，如我的就是：

```bash
mssh bro
```
