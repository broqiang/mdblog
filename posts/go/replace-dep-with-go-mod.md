---
title: "将 dep 更换为 go mod"
author: "BroQiang"
created_at: 2019-08-28T11:10:07
updated_at: 2019-08-28T11:10:07
---

## 原因

之前一直在使用 dep 来管理项目， 最近碰到了两个恶心的事情， 就是安装 ali-oss-sdk 和
go-ethereum 客户端的时候超级慢，lantern 也是半费状态， 时断时连，go-ethereum
等了半个小时都没装上（开灯和关灯都试了）， 最后 dep ensure 状态回家， 第二天上班才好。

正好在[learnku](https://learnku.com/) 的帖子中 appleboy 给我留言说
[可以改用 go module 了] ， 然后就试了下， 泪奔， 只能说太好用了。

结果就是我就将所有的项目全部用 module 替换了。

## 替换

这里记录下过程， 用我的博客来做示例。

这个 `go mod` 是 go 自带的， 替换起来非常简单， 正好我的 go 环境使用的是最 go1.12 ，
配置起来很容易。

### 将原本的 dep 的配置删除

```bash
# 我已经将代码从 GOPATH 中改到 /www 目录了， 现在可以脱离 GOPATH 了,
# 代码也可以放在任意目录了。
# 如果代码还在 GOPATH 下， 就需要手动添加下面配置， 不叫 go 去自动识别
# export GO111MODULE=on
cd /www/mdblog

# 删除原本的 dep 的配置， 这个我也不纠结了， 原本 dep 对版本的控制就不太好
# 如果对使用的包版本有要求， 稍后可以手动修改 go.mod 将版本改为指定的
rm -rf Gopkg.* vendor
```

### 配置 go mod

这个也非常简单，它初始化后会在 $GOPATH/pkg/mod 目录中将所有下载过的依赖包保存，
并且可以保存多个版本， 下次再使用已经下载过的版本依赖时， 不会再去网上下载，
应该算是个本地仓库的感觉吧。

```bash
# 初始化
go mod init

# 如果是一个新的项目， 初始化的时候需要指定项目名称， 如下面两个例子
# go mod init demoname
# go mod init github.com/broqiang/mdblog
```

现在就会在项目目录下生成一个 go.mod 文件， 并且当在当前项目目录下执行 `go get` ,
`go build`, `go run` 等命令的时候， 自动将依赖的包下载, 同时将版本信息写入到
go.mod 中，还会生成一个新的 go.sum 详细的记录

就这么简单的将项目从 dep 转成 module 管理。

## 代理

换成 module 方式之后可以很方便的来使用代理了，例如可以使用
[https://goproxy.io/](https://goproxy.io/) 代理， 里面有配置的说明。

Linux 下添加下面变量就可以生效

```bash
# Enable the go modules feature
export GO111MODULE=on
# Set the GOPROXY environment variable
export GOPROXY=https://goproxy.io
```

这里有个坑， 如果使用了私有仓库， 这个代理是找不到的， 正确的方法不清楚， 我的处理办法是：
当项目使用了私有仓库（公司项目全是私有的）时， 临时手动把代理去掉
`export GOPROXY=""` 。
