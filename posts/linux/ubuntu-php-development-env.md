---
title: "Ubuntu PHP 开发环境配置"
author: "BroQiang"
github_url: "https://broqiang.com"
created_at: 2017-12-27T01:58:46
updated_at: 2017-12-27T01:58:46
description: "记录了用 Ubuntu 作为开发环境的一些简单配置，此文档是基于 Ubuntu 17.10 编写"
tags: ["linux", "php"]
---

> 版本 17.10 64 位（ubuntu-17.10-desktop-amd64.iso）

## 安装操作系统

这里不太多的赘述，基本上就是下一步下一步就可以了，每一个地方点击下一步之前注意看一下即可。

系统安装完成后，可以通过 `Ctrl+Alt+T` 打开终端，后面所有的命令都在终端中执行。

## 安装 VMTools

如果是虚拟机安装的操作系统就安装下这个，如果是物理机安装可以跳过这步。

```bash
sudo apt install open-vm-tools-desktop -y
```

安装完成后，就可以随便调整桌面的大小，并且可以和物理机之间复制文件和剪切板，如果不能生效，可以将当前登录的桌面注销或者重启操作系统。

## 更改家目录为英文

```bash
LANG=en_US
xdg-user-dirs-gtk-update
```

然后点击填出界面右下角的 `Update Names`

## 去除不必要的软件

这一步不是必须的，个人觉得这样可以减少一些磁盘空间，并且在更新操作系统的时候可以不更新那些用不到的包，可以节省一些时间。

```bash
sudo apt autoremove libreoffice-* ibus -y
```

## 安装需要的软件包

包括：主题、图标、vim、numix-icon-theme 等

```bash
# 更新源
sudo apt update

# 安装软件
sudo apt install -y gnome-tweak-tool arc-theme numix-icon-theme vim \
gnome-icon-theme git net-tools shutter kazam gnome-shell-extensions
```

## 更新操作系统到最新

```bash
sudo apt upgrade -y
# 更新完成后重启操作系统
reboot
```

重启后登录桌面的时候如果是虚拟机选择 `Ubuntu on xorg` 启动，不要选择默认，因为现在很多软件不支持 `Wayland`，如果是物理机可以尝试一下，如果显卡可以直接 xorg 也可以切换，如果不行就只能 使用 Wayland 了。

更新系统后，会出现两个内核，可以通过下面方式删除旧的内核

```bash
# 查询出当前系统中的内核
sudo dpkg -l | grep kernel
# 用眼睛找到版本老的版本，然后将包名复制删除即可
sudo dpkg -l | grep kernel | grep 4.13.0-16.19 | awk '{print $2}' | xargs sudo apt autoremove -y
```

## 安装显卡驱动

如果是虚拟机这部可以跳过了

这里只记录了 MShybrid 显卡的安装方式，并且实测的只有雷神 st pro，未来如果还会用到其他型号笔记本会将安装方法记录。

详细方法见：[Ubuntu 17.10 安装 MShybrid 显卡](https://broqiang.com/2018/01/13/Ubuntu17.10ForLeiShen/)

## 配置主题

这个看个人爱好，区别只是美观，我使用的是 arc-theme 和 Numix 的图标，上面安装软件的时候已经安装。

快捷键 `Alt+F2` 在窗口中输入 `gnome-tweak-tool`，然后就按照喜好去配置就可以了。

## 设置快捷键

### 系统快捷键

打开设置

- 右上角菜单点进设置界面

- 快捷键 `Alt+F2` 输入 `gnome-control-center`

左侧菜单找到`设备`->`键盘` ，在右侧设置快捷键，也可以在最下面点击 `+ 号` ，自定义快捷键。

### 设置 gnome-terminal （终端）快捷键

因为这个使用频率非常的高，以前使用 xshell 习惯了全屏 Alt+Enter，所以这个也改成这个。

菜单上的`编辑`->`首选项`->`快捷键` ，然后根据需要去配置就可以了。

## 配置搜狗输入法

前面已经安装了，如果没装用下面命令安装即可，不过需要注意，要配置 ubuntukylin 源。

可以先从[官网下载](https://pinyin.sogou.com/linux/) deb 包

```bash
# 下载搜狗拼音
sudo dpkg -i sogoupinyin_2.2.0.0102_amd64.deb
sudo apt install -f
```

打开`设置`界面->`区域和语言`->`Manage Installed Languages`

在弹出界面选择 fcitx，然后通过 `Ctrl+空格` 就可以切换输入法。

## 配置 git

软件已经在前面 apt install 的时候安装，下面进行一些简单的配置，如果想要手动安装，参考：[Git 安装](https://broqiang.com/posts/8)

```bash
# 提交数据的时候的用户邮箱和用户名
git config --global user.email "broqiang@qq.com"
git config --global user.name "Bro Qiang"
# 保存密码
git config --global credential.helper store

# 设置使用 vim 作为默认文本编辑器，nano 实在是用不习惯……
git config --global core.editor vim

# 配置 git 别名
echo -e "\nalias gs='git status'\nalias gaa='git add .'\nalias ga='git add '\nalias gp='git push'\nalias gc='git commit -m '\nalias gl='git log'\nalias grao='git remote add origin '\nalias gpo='git push origin '\nalias gb='git branch'" >> ~/.bash_aliases
```

## 配置 vim

```bash
echo -e ":set nu\n:set tabstop=4\n:set expandtab\n:set autoindent\n:set shiftwidth=4\n:set smartindent" > ~/.vimrc
```

## 安装 lantern

```bash
# 下载
wget https://raw.githubusercontent.com/getlantern/lantern-binaries/master/lantern-installer-64-bit.deb

# 安装
sudo dpkg -i lantern-installer-64-bit.deb

# 如果提示需要依赖关系，执行下面命令
sudo apt install -f
```

## 安装 Sublime

下面的是一个已经按照我的习惯配置好了的 Sublime，更多配置参考本博客的 Sublime 相关文档

- [Github 仓库](https://github.com/BroQiang/Software_Sublime.git)

- [码云仓库](https://gitee.com/BroQiang/Software_Sublime.git)

```bash
git clone https://github.com/BroQiang/Software_Sublime.git
```

## 安装 Atom

具体步骤：

- [Github 仓库](https://github.com/BroQiang/Software_Atom.git)

- [码云仓库](https://gitee.com/BroQiang/Software_Atom.git)

## 安装 chrome 浏览器

直接去官网安装下载安装即可，如果打不开确定下前面的 alntern 是否安装

[chrome 浏览器官方下载地址](https://www.google.cn/intl/zh-CN/chrome/browser/desktop)

```bash
# 安装软件包
sudo dpkg -i google-chrome-stable_current_amd64.deb
# 解决依赖关系
sudo apt install -f
```

使用，直接 `Alt+F2` 输入 `google-chrome`

## 安装 Navicat

下载地址： [https://gitee.com/BroQiang/Software_Navicat.git](https://gitee.com/BroQiang/Software_Navicat.git)

里面有详细的安装步骤

安装完成可以配置一个启动脚本，方便快速启动

```bash
echo -e "/opt/navicat112_Premium_cs_x64/start_navicat" > ~/bin/navicat
chmod +x ~/bin/navicat
```

使用，可以直接在终端或者 `Alt+F2` 输入 `navicat` 就可以运行了。

## 配置截图工具

### shutter

在 Ubuntu 下没有什么太好用的截图工具，尝试了几个之后，确定为 `shutter`，这个工具很强大，说它不好用是因为它太重了，有点不够灵巧，功能有点太多了。

使用方法，可以直接通过 `shutter`命令来截图，并可以通过系统设置，给截图命令配置快捷键，方便快速使用截图。

- `shutter -s` 区域截图

- `shutter -f` 全屏截图

- `shutter -m` 截取菜单（一般工具没法截取右键菜单，这个很强大）

### flameshot `（推荐）`

如果使用的不是 xorg， 而是 Wayland，这个是个不错的选择。

#### 安装

```bash
# 下载代码
git clone https://github.com/lupoDharkael/flameshot.git

# 安装依赖关系
sudo apt install -y g++ build-essential qt5-qmake qt5-default qttools5-dev-tools

# 编译
qmake && make

# 安装
make install
```

#### 使用

可以直接通过命令使用，为了方便，可以将下面命令配置成自己习惯的快捷键。

```bash
# 开始截图，直接就是截图界面，有些类似 QQ 截图和 深度截图，不过比 QQ 截图更加强大
flameshot gui

# 开始截图，并配置默认的保存位置，和上面的界面相同，只是有了一个默认保存位置
# 创建默认保存目录，只要是真实存在并且有写入权限的目录即可，可以根据跟人爱好去配置
mkdir -p ~/Pictures/Screenshot
# 截图后图片会默认保存到 Screenshot 目录
lameshot gui -p ~/Pictures/Screenshot

# 延时启动，单位事毫秒，如 2000 就是 2秒，这种方式适用于截取右键菜单，也是比较方便的
flameshot gui -d 2000
```

成功截取图片后，可以通过下面两个快捷键来保存图片：

- `Ctrl+S` - 保存到磁盘

- `Ctrl+C` - 保存到剪切板

## 配置录屏工具

这里使用的是 `kazam` ，在前面已经安装了，如果没有安装再安装一下

### 快捷键

- 开始录 - `Supper+Ctrl+R`

- 暂停 - `Supper+Ctrl+P`

- 结束保存 - `Supper+Ctrl+F`

## 配置 WPS

详细步骤见 [github](https://github.com/BroQiang/Software_WPS)

## 禁用 IPv6

这个看个人看好吧，不禁用也不影响使用，只是查询来的内容会多不少，看起来有点乱，看需求去决定是否禁用。

```bash
sudo vim /etc/default/grub
# 找到下面内容
GRUB_CMDLINE_LINUX_DEFAULT="quiet splash"
# 讲上面内容替换成
GRUB_CMDLINE_LINUX_DEFAULT="ipv6.disable=1 quiet splash"

# 更新
sudo update-grub

# 重启后生效
reboot
```

## 安装配置 MySQL PHP Nginx 环境

可以使用 [自动安装脚本](https://github.com/BroQiang/lnmp)，自动安装环境，使用的时候注意下 config 文件中的配置，如：操作系统版本、用户等。

如果想要手动配置，请参考(注意下安装顺序，建议按照下面的顺序安装)：

- [Linux 编译安装 Mysql5.7](https://broqiang.com/posts/11)

- [Linux 编译安装 PHP 7](https://broqiang.com/posts/12)

- [Linux 编译安装 Nginx](https://broqiang.com/posts/13)

## 安装 Composer

没有什么特别的，下载下来就可以使用，不过要保证系统中有 PHP CLI 环境。

### 开始安装

```bash
sudo wget https://getcomposer.org/download/1.5.6/composer.phar -O /usr/local/bin/composer
sudo chmod +x /usr/local/bin/composer
```

### 配置国内镜像

[Composer 中文网](https://www.phpcomposer.com/) 可以在这里看到配置方式，和中文的各种文档。

```bash
composer config -g repo.packagist composer https://packagist.phpcomposer.com
```

## 安装 node 环境

### 安装 node 和 npm

```bash
sudo apt install nodejs npm
```

### 安装 yarn

因为 npm 的效率等问题，不是非常好用，所以一般都将 npm install 的工作用 yarn 来替换

```bash
# 如果没有安装 curl 就先安装，如果已经安装可以忽略
sudo apt install -y curl

# 导入仓库的秘钥
curl -sS https://dl.yarnpkg.com/debian/pubkey.gpg | sudo apt-key add -
# 配置 yarn 的仓库，一定注意，要导入仓库在安装，Ubuntu 默认带的 yarn 是个莫名其妙的东西，反正不是这个 yarn 就对了
echo "deb https://dl.yarnpkg.com/debian/ stable main" | sudo tee /etc/apt/sources.list.d/yarn.list

# 更新软件源，这个也是必须要做的，否则安装的时候还是系统自带的 yarn
sudo apt-get update

# 安装
sudo apt-get install yarn
```

### 配置国内源

这里采用的都是淘宝的镜像，其他的没有尝试过

```bash
# npm 直接将下面内容写去 到 ~/.bash_alias 中即可
alias cnpm="npm --registry=https://registry.npm.taobao.org \
--cache=$HOME/.npm/.cache/cnpm \
--disturl=https://npm.taobao.org/dist \
--userconfig=$HOME/.cnpmrc"

# yarn 国内镜像
yarn config set registry https://registry.npm.taobao.org
```

## 安装屏幕画笔

这个看个人需求吧，因为经常会有做演示的需求，个人还是非常有需求的。

```bash
sudo apt-get install gromit-mpx
```

使用非常简单，可以通过菜单找到，直接点击就可以运行，为了方便也可以放在 dock 中，或者直接快捷键 `Alt+F2` 输入 `gromit-mpx`，这个时候就将软件启动了，可以见右上角。

### 屏幕画快捷键

- F9 ---- 开始画图，再次点击就停止画图

- Shift+F9 清空绘图，这个如果不手动清空，画的图案会一直在屏幕上显示

- Ctrl+F9 隐藏桌面绘图，但是不清空，可以再次点击显示

- 画图时默认是红色的线，按住 Shift 是蓝线，安装 Ctrl 是黄线

## 画图工具

推荐使用 dia，可以直接普通图案、UML 图、流程图等。

### 画图工具安装

```bash
sudo apt install dia
```

### 解决不能输入中文

默认是不能输入中文的，只能用一个折中的办法，在启动的时候执行 `dia --classic` ，不过这样导航栏就拆离出来了，暂时没有更好的解决方案，并且好几年了都没有人解决……
