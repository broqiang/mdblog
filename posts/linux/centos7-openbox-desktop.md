---
title: "CentOS 7 极简桌面环境"
author: "BroQiang"
github_url: "https://broqiang.com"
created_at: 2019-06-13T11:03:00
updated_at: 2019-06-13T11:03:00
description: ""
tags: ["linux"]
---

最近因为要使用 Windows 操作系统， 在虚拟机中使用 Ubuntu 18.04 感觉有些卡， 比物理机直接使用反应慢了很多，
所以决定搞个轻量一点的， 最终决定安装 CentOS7 的 Minimal 包， 然后配个 Openbox 桌面。

## 初始配置

关闭 SELinux ， 开发环境最好关上， 避免不少麻烦。 生产环境也看需求， 个人也不推荐打开。

```bash
sudo sed -i 's/SELINUX=.*$/SELINUX=disabled /' /etc/selinux/config
```

安装 epel 源、 bash 自动完成工具、 网络工具，开发工具等常用的工具，可以根据实际需要安装。

将操作系统更新到最新， 重启。

```bash
sudo yum install -y epel-release bash-completion net-tools gcc gcc-c++ make cmake autoconf
sudo yum update -y
reboot
```

> 因为个人习惯不直接使用 root 用户， 所以配置了一个 sudo 权限的用户， 如果使用的是 root 用户， 后面命令中的 sudo 全部可以省略

配置个 PS1 ， 为了叫命令行的前缀短一点， 这个看个人喜好去配置， 或者不配置也什么都不影响

```bash
echo "export PS1='[\[\033[1;34m\]\w\[\033[0m\]]\$ '" >> ~/.bashrc
```

## 编译安装 git

CentOS7 自带的 git 有点老， 所以采用源码安装。

可以参考 [官方文档](https://git-scm.com/book/zh/v2/%E8%B5%B7%E6%AD%A5-%E5%AE%89%E8%A3%85-Git) 安装，
源码可以从 [官方仓库](https://mirrors.edge.kernel.org/pub/software/scm/git/) 下载,
这里采用的是写此文档时的最新版本。

```bash
# 安装依赖关系
sudo yum install -y curl-devel expat-devel gettext-devel \
    openssl-devel zlib-devel

# 下载源码
curl -fO https://mirrors.edge.kernel.org/pub/software/scm/git/git-2.22.0.tar.xz

# 解压安装包
tar xf git-2.22.0.tar.xz && cd git-2.22.0/

# 编译安装
make configure
./configure --prefix=/usr
make all
make install
```

安装完成后可以查看下版本， 显示下面信息及安装成功

```bash
$ git --version
git version 2.22.0
```

安装完成后可以按照需要做一些简单的配置。

```bash
# 提交数据的时候的用户邮箱和用户名
git config --global user.email "broqiang@qq.com"
git config --global user.name "Bro Qiang"
# 保存密码
git config --global credential.helper store

# 设置使用 vim 作为默认文本编辑器，nano 实在是用不习惯……
git config --global core.editor vim

# 配置方便使用的别名， 根据个人习惯， 也可以不配置
echo -e "\n\nalias gs='git status' \
\nalias gaa='git add . \
'\nalias ga='git add ' \
\nalias gp='git push' \
\nalias gc='git commit -m ' \
\nalias gl='git log' \
\nalias grao='git remote add origin ' \
\nalias gpo='git push origin ' \
\nalias gb='git branch'" >> ~/.bashrc
```

## 安装 vim

CentOS 7 自带的 vim 版本太低了， 很多 vim 插件都不支持或支持的不好， 比如 YouCompleteMe , 所以就直接编译安装最新版本了。

可以从 [github.com/vim/vim](https://github.com/vim/vim) 找到最新的版本， 我这里使用的是 release 中的最新版， 没有使用 master

### 编译安装 vim

```bash
curl -fLO https://github.com/vim/vim/archive/v8.1.1523.tar.gz
tar xf v8.1.1523 && cd vim-8.1.1523/

# 安装依赖关系
sudo yum install -y gcc-c++ ncurses-devel python-devel

# 编译安装
./configure \
  --disable-nls \
  --enable-cscope \
  --enable-gui=no \
  --enable-multibyte  \
  --enable-pythoninterp \
  --enable-rubyinterp \
  --with-features=huge  \
  --with-python-config-dir=/usr/lib64/python2.7/config \
  --with-tlib=ncurses \
  --without-x
make
make install
```

编译完成后输入下面命令查看是否安装完成

```bash
$ vim --version
VIM - Vi IMproved 8.1 (2018 May 18, compiled Jun 12 2019 22:29:49)
...... 省略更多信息
```

### 配置 vim

```bash
# 安装插件管理工具
curl -fLo ~/.vim/autoload/plug.vim --create-dirs \
    https://raw.githubusercontent.com/junegunn/vim-plug/master/plug.vim

# 下载 .vimrc 配置文件
curl -fLo ~/.vimrc \
    https://raw.githubusercontent.com/BroQiang/vim-go-ide/master/vimrc-centos-base
```

上面配置完成后， 直接打开 vim ， 输入 `:PlugInstall` ， 慢慢等待插件安装， 看网络状况， 有可能会要等很久。

插件下载安装完成后，需要编译下 YouCompleteMe

```bash
cd ~/.vim/plugged/YouCompleteMe
python install.py
```

编译完成后就可以比较优雅的使用 vim 了，这里没有写具体步骤，直接使用的我自己的配置，更多的配置可以参考
[vim-go-ide](https://github.com/BroQiang/vim-go-ide)

## 安装桌面环境

为了方便， 安装一个极简的 openbox 桌面 ， 可以保证虚拟机中桌面环境的性能，
同时又能做一些桌面下的操作，如果物理机安装的话个人还是更喜欢 Ubuntu + GNOME 桌面。

```bash
# 安装 X Window 环境
sudo yum groupinstall -y "X Window System"

# 安装 openbox 桌面环境和一些相关软件， 加--exclude 选项， 过滤开发相关的包
sudo yum install -y xfce4-terminal wqy-* \
    lightdm openbox gmrun tint2 obconf-qt thunar

# 安装登录管理器LightDM
sudo systemctl enable lightdm

# 配置开机默认启动图形界面
sudo systemctl set-default graphical.target

sudo reboot
```

> 如果桌面不能最大化， 可以安装 VMtools

### 配置 tint2 托盘自动启动

tint2 就是一个简易的系统托盘， 可以显示一些基本的内容， 如果不需要的话也可以不配置。

编辑 autostart

```bash
sudo vim /etc/xdg/openbox/autostart
```

在下面写入下面内容

```bash
tint2 &
```

> autostart 是一个配置文件， 想要什么程序在启动 openbox 的时候自己启动， 写入这里就可以了。

### 配置快捷键

配置 gmrun (这是一个运行可执行程序的一个工具， 和 GNOME 桌面下的 `Alt+F2` 的功能类似，
所以这里也给它配置一个 `Alt + F2` 的快捷键） 、 xfce4-terminal (终端）和
thunar （轻量级的文件管理器）。

编辑 rc.xml （一般都会有， 如果不存在的话就自己从 /etc/xdg/openbox/rc.xml 复制一个即可）

```bash
vim ~/.config/openbox/rc.xml
```

找到 `<keyboard>` 标签， 在 `<keyboard>` 和 `</keyboard>` 之间插入一组新的标签

```bash
<keybind key="A-F2">
  <action name="execute">
    <execute>gmrun</execute>
  </action>
</keybind>
<keybind key="C-A-T">
  <action name="execute">
    <execute>xfce4-terminal</execute>
  </action>
</keybind>
<keybind key="W-E">
  <action name="execute">
    <execute>thunar</execute>
  </action>
</keybind>
<keybind key="A-Return">
  <action name="ToggleDecorations"/>
  <action name="ToggleMaximizeFull"/>
</keybind>
<keybind key="W-Up">
  <action name="GrowToEdge">
    <direction>north</direction>
  </action>
</keybind>
<keybind key="W-Down">
  <action name="GrowToEdge">
    <direction>south</direction>
  </action>
</keybind>
<keybind key="W-Left">
  <action name="GrowToEdge">
    <direction>west</direction>
  </action>
</keybind>
<keybind key="W-Right">
  <action name="GrowToEdge">
    <direction>east</direction>
  </action>
</keybind>
```

删除另外一个已经存在的 W-E ， 这个默认的貌似是个 kde 的文件管理器， 包有点多， 就直接把他删了，
找到下面内容删除即可。

```bash
<!-- Keybindings for running applications -->
<keybind key="W-e">
  <action name="Execute">
    <startupnotify>
      <enabled>true</enabled>
      <name>Konqueror</name>
    </startupnotify>
    <command>kfmclient openProfile filemanagement</command>
  </action>
</keybind>
```

> 这是一个 xml 格式的配置文件， 注意下标签的闭合， 只要不放在其他子标签中即可。

## 安装 google chrome 浏览器

```bash
curl -fO https://dl.google.com/linux/direct/google-chrome-stable_current_x86_64.rpm

sudo yum localinstall -y google-chrome-stable_current_x86_64.rpm
```

## 安装 sogou 拼音

虽然也有一些其他输入法， 不过都用不太习惯， 还是习惯 sogou 输入法， 不过官方只有 dep 包， 只能自己处理下了。

### 安装 fcitx

```bash
sudo yum -y  install fcitx fcitx-pinyin fcitx-configtool qtwebkit
```

### 处理 dep 包

先从 [搜狗官方](https://pinyin.sogou.com/linux/?r=pinyin) 将 dep 包下载，
这里下载的是 64 位的包， 32 位没测试过。

```bash
# 下载 dep 包
curl -O http://cdn2.ime.sogou.com/dl/index/1524572264/sogoupinyin_2.2.0.0108_amd64.deb

# 将 dep 包解包
mkdir sogou
mv sogoupinyin_2.2.0.0108_amd64.deb sogou
cd sogou
ar vx sogoupinyin_2.2.0.0108_amd64.deb

# 将上面打开的 data 包解压, 手动复制 sogoupinyin 的文件
tar xf data.tar.xz

# 复制搜狗拼音的文件, 需要注意， 是相对路径的 user ， 不是从 / 开始的，
# 文件是从 data.tar.gz 中解压出来的
sudo cp usr/lib/x86_64-linux-gnu/fcitx/* /usr/lib64/fcitx/
sudo cp -r etc/* /etc/
sudo cp -r usr/bin/* /usr/bin/
sudo cp -r usr/share/* /usr/share/
```

### 配置 fcitx 自动启动

编辑配置文件 `autostart`

```bash
sudo vim /etc/xdg/openbox/autostart
```

在末尾添加下面内容

```bash
fcitx &
```

注销， 下次重新登录 openbox 桌面的时候就会自动启动 fcitx 输出法

### 配置 fcitx 加载 sogoupinyin

打开 fcitx 配置管理工具， 可以通过也右下角的键盘图标，或直接输入 `fcitx-config-gtk3` 命令，
点击左下角的 + 号， 找到 Sogou Pinyin ， 添加即可。 一般保留一个英文输入法和一个 Sougou
输入法， 使用 Ctrl + 空格就可以切换了。

### DPI 设置

如果是高分屏， 默认显示会非常的小， 设置下 dpi （类似 Windows 下的缩放功能）即可。

参考： [https://wiki.archlinux.org/index.php/HiDPI#X_Resources](https://wiki.archlinux.org/index.php/HiDPI#X_Resources)

`vim ~/.Xresources` (~/.Xresources 如果不存在就新建一个)

写入下面内容， 然后重新登录桌面即可生效

```bash
Xft.dpi: 130
Xft.autohint: 0
Xft.lcdfilter:  lcddefault
Xft.hintstyle:  hintfull
Xft.hinting: 1
Xft.antialias: 1
Xft.rgba: rgb
```

上面的 `Xft.dpi: 130` 根据自己的分辨率来设置

## 完成

到这里基本就完成了， 为了虚拟机的性能， 没有配置过多花哨的东西， 可以满足基本的开发就可以了。
