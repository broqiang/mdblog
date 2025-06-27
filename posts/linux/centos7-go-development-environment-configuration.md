---
title: "CentOS7 Go 开发环境配置"
author: "BroQiang"
github_url: "https://broqiang.com"
created_at: 2018-06-05T19:03:14
updated_at: 2018-06-05T19:03:14
description: "CentOS7.5 G语言开发环境配置"
tags: ["linux", "go"]
---

因为最近将开发环境从 Ubuntu 换成了 CentOS ，所以写了这个文档记录一下。

## 操作系统

操作系统安装的是 mini 版本，最小安装。 安装后的显卡等硬件驱动自行解决，不在这里说明。

### 安装桌面环境

安装前先将操作系统更新到最新。

```bash
sudo yum update -y
```

我使用的是 GNOME 桌面，所以这里安装的是 GNOME 。

```bash
sudo yum groupinstall 'GNOME Desktop' -y
```

安装完修改启动级别，将默认的字符界面修改成图形界面启动。

```bash
sudo systemctl set-default graphical
```

配置完成后重启系统生效，进入图形界面前选择 GNOME ，在输入界面有个 ⚙ ，点击就可以选择，个人不喜欢默认的 GNOME 经典模式。

### 配置默认快捷键

打开设置

- 右上角菜单点进设置界面

- 或快捷键 `Alt+F2` 输入 `gnome-control-center`

左侧菜单找到设备->键盘 ，在右侧设置快捷键，也可以在最下面点击 + 号 ，自定义快捷键。

#### 设置的快捷键

> 这个按照个人习惯去配置

- 启动器->主目录 `Super+E`

- 导航->窗口上移一个工作区 `Shift+Ctrl+Alt+K`

- 导航->窗口下移一个工作区 `Shift+Ctrl+Alt+J`

- 导航->移动到上层工作区 `Ctrl+Alt+K`

- 导航->移动到下层工作区 `Ctrl+Alt+J`

- 导航->隐藏所有正常窗口 `Super+D`

- 窗口->切换全屏模式 `Alt+回车`

- 添加自定义快捷键-> `gnome-terminal` -> `Ctrl+Alt+T`

### 更改家目录为英文

这个看个人喜好，个人不喜欢中文目录，因为命令行的时候要经常切换输入法，比较麻烦。

```bash
LANG=en_US
xdg-user-dirs-gtk-update
```

然后点击填出界面右下角的 `Update Names`

### 安装需要的软件

epel 源、主题、配置工具等

```bash
# 安装 epel 源
sudo yum install epel-release -y

# 安装软件
sudo yum install gnome-tweak-tool vim net-tools arc-theme numix-icon-theme wqy-*
```

### 安装 git

可以直接 `sudo yum install git` 安装，不过这个版本有点低，个人建议编译安装。这里只简单记录了安装步骤，详细的说明请参见： [Git 安装](https://broqiang.com/posts/8)

```bash
# 安装依赖关系
sudo yum install curl-devel expat-devel gettext-devel openssl-devel \
zlib-devel  perl-ExtUtils-MakeMaker autoconf gcc gcc-c++ asciidoc xmlto docbook2X

# CentOS 有个 bug，docbook2X 安装完成后的包名叫 db2x_docbook2texi ，所以需要创建一个 软链接
sudo ln -s /usr/bin/db2x_docbook2texi /usr/bin/docbook2x-texi

# 下载当前最新版本的 git
wget https://www.kernel.org/pub/software/scm/git/git-2.17.1.tar.gz

# 解压进入源码目录
sudo tar xzvf git-2.17.1.tar.gz -C /usr/local/src/
cd /usr/local/src/git-2.17.1

# 如果是 github 上下载的版本，需要执行下面命令，官方下载的可以忽略
sudo make configure

# 配置 MakeFile 文件
sudo ./configure --prefix=/usr/local/git

# 编译、安装
sudo make all doc info
sudo make install install-doc install-html install-info

# 配置环境变量
echo -e "export GIT_HOME=/usr/local/git\nexport PATH=\$GIT_HOME/bin:\$PATH" \
| sudo tee /etc/profile.d/git.sh

# 生效环境变量
. /etc/profile.d/git.sh

# 删除被 yum 安装的 git
sudo yum autoremove git
```

### 安装 sogou 拼音

官方没有 CentOS 版本的安装包，可以到 [我的码云](https://gitee.com/BroQiang/Software_CentOS_SogouPinYin) 下载，里面有详细的配置文档。

```bash
# 删除 ibus
sudo yum autoremove ibus

# 获取软件
wget https://gitee.com/BroQiang/Software_CentOS_SogouPinYin/raw/master/sogou-pinyin-1.1.0.0037-1.el7.centos.x86_64.rpm

# 安装
sudo yum install sogou-pinyin-1.1.0.0037-1.el7.centos.x86_64.rpm -y

# 配置环境变量
echo -e "export GTK_IM_MODULE=fcitx\nexport XMODIFIERS=\"@im=fcitx\"\nexport QT_IM_MODULE=fcitx" | sudo tee /etc/profile.d/fcitx.sh

# 切换默认输入法
imsettings-switch fcitx

# 重启系统，配置生效
reboot
```

如果出现一个面板加载失败，请重新启动的错误，应该是缺少了一个依赖： opencc 。

```bash
sudo yum install opencc -y
```

### 安装 lantern

到 [Github](https://github.com/getlantern/lantern/releases) 下载

需要安装一个依赖关系：

```bash
sudo yum install -y libappindicator-gtk3
```

不过在 Linux 下使用有 bug ，可以参考： [Linux 下 Lantern 关闭后不清除代理](https://broqiang.com/posts/27)

### 安装 google-chrome 浏览器

这个需要科学上课，要打开上面安装的 lantern

到 [官网](https://www.google.com/chrome/) 去下载

下载后通过 yum 去安装，可以自动解决依赖。

```bash
sudo yum install google-chrome-stable_current_x86_64.rpm -y
```

## Go 环境

如果网络顺畅的话可以直接去 [官网](https://golang.org) ，如果官网打不开也可以选择 [Go 语言中文网](https://studygolang.com) 下载及查看文档。

这里使用的是当前官方最新版本： `1.10.2`

### 下载

#### 官网

下载地址： [https://dl.google.com/go/go1.10.2.linux-amd64.tar.gz](https://dl.google.com/go/go1.10.2.linux-amd64.tar.gz)

```bash
wget https://dl.google.com/go/go1.10.2.linux-amd64.tar.gz
```

#### Go 语言中文网

下载地址： [https://studygolang.com/dl/golang/go1.10.2.linux-amd64.tar.gz](https://studygolang.com/dl/golang/go1.10.2.linux-amd64.tar.gz)

```bash
wget https://studygolang.com/dl/golang/go1.10.2.linux-amd64.tar.gz
```

### 解压及配置

如果是下载的二进制版，直接解压就可以使用（Ubuntu 和 CentOS 一般都可以直接使用二进制版本，其他发行版本没有测试过）。

```bash
# 解压
sudo tar xzf go1.10.2.linux-amd64.tar.gz -C /usr/local/

# 配置环境变量
echo -e "export GOROOT=/usr/local/go\nexport PATH=\$PATH:\$GOROOT/bin" | sudo tee /etc/profile.d/go.sh

# 生效及测试
source /etc/profile.d/go.sh
go version
```

看到下面的内容，证明配置成功

```bash
go version go1.10.2 linux/amd64
```

### 配置 GOPATH 环境变量

工作空间就是 GOPATH 环境变量， 如果不配置的话，默认就是 `$HOME/go` ，可以自定义一个目录，如 `$HOME/code/go`。

工作空间下要有 3 个目录，`bin`、`pkg`、`src`，具体的介绍见 [如何使用 Go 编程](http://docscn.studygolang.com/doc/code.html)

#### 配置

```bash
# 创建目录
mkdir -p ~/code/go/{bin,pkg,src}

# 配置环境变量， 这个一般配置在当前用户的 .bashrc 中即可
echo -e "\nexport GOPATH=\$HOME/code/go\nexport PATH=\$PATH:\$GOPATH/bin" >> ~/.bash_profile
```

需要重启系统或从新登录桌面环境生效，如果是 xorg 桌面环境也可以直接 `Alt +F2`，输入 `r`

## 安装配置 IntelliJ IDEA

这个看个人喜好了，Go 支持几乎所有主流的 IDE 和 文本编辑器

### 下载 IDEA

直接从 [官网](https://www.jetbrains.com/idea/download/#section=linux) 下载即可

```bash
wget https://download.jetbrains.8686c.com/idea/ideaIU-2018.1.4.tar.gz
```

### 安装

直接解压就可以使用

```bash
sudo tar xzf ideaIU-2018.1.4.tar.gz -C /opt/
```

解压首次启动可以通过绝对路径

```bash
/opt/idea-IU-181.5087.20/bin/idea.sh
```

启动后会提示是否导入配置文件，如果是第一次，就选择默认的不导入即可，如果已经有了配置，就将备份的配置导入。

将他的协议下拉到最下（也可以看完），点击 Accept ，也就是传说中的同意协议。

License 授权，如果不差钱就去买一个吧，也可以使用 [http://idea.youbbs.org](http://idea.youbbs.org) ，来自： [https://www.youbbs.org/t/2115](https://www.youbbs.org/t/2115)。

选择一个主题

创建一个应用菜单

在 /usr/local/bin/idea 创建一个脚本，下载启动就不需要绝对路径了，只需要 `Alt+F2` 输入 idea 即可。

选择 idea 的任务（选择启动的插件），我只留了一个 web 、 版本控制、 数据库，其他的全部都禁用了。

然后看到一个可以创建任务的界面，软件安装完成。这时就可以将软件关闭，可以通过菜单中的图标或快捷方式启动即可。

### 安装插件

现再 idea 还没法开发 Go ，还需要安装和 Go 相关的插件。

点击下侧的 configure -> Plugins -> 浏览 -> Go ，将 Go 插件安装 idea 就可以支持 Go 了，然后就可以创建 Go 项目了。 如果下载不成功可以换个时间再下或者单独将插件下载下来选择本地安装。

## VSCode 安装

如果觉得 idea 或者 goland 太重，可以考虑使用 vscode ，用起来也不错。

### 获取软件

[官网下载地址](https://code.visualstudio.com) ，选择对应的包下载即可

### 安装 VSCode 插件

`Ctrl+Shift+X` 输入 `go` 选中安装即可，安装完成后会依赖一些插件，都是 go 的包，第一次打开 go 文件的时候会提示安装，从 vscode 安装可能会有些慢，也可以手动通过终端( cli ) 来安装,输入下面命令即可：

> > 注意这里是一条命令，不要有换行

```bash
go get -v {github.com/mdempsky/gocode,\
github.com/ramya-rao-a/go-outline,github.com/uudashr/gopkgs/cmd/gopkgs,\
github.com/acroca/go-symbols,golang.org/x/tools/cmd/guru,github.com/acroca/go-symbols,\
golang.org/x/tools/cmd/gorename,github.com/go-delve/delve/cmd/dlv,\
github.com/rogpeppe/godef,github.com/lukehoban/go-find-references,\
github.com/sqs/goreturns}
```
