# Bro Qiang 博客

示例： [broqiang.com](https://broqiang.com) （个人博客，正在使用）

源码： [github.com/broqiang/mdblog](https://github.com/broqiang/mdblog)

## 环境

- 开发环境 Ubuntu 18.04

- 服务器环境 CentOS 7.6

其他环境没有测试过，不确定是否兼容（如 Windows、MAC）， 后面的所有步骤都假定是在 Linux（上面两种系统） 下完成。

## 快速使用

如果不需要自己编译，只是想要查看下效果，可以直接下载[mdblog.tar.gz](https://github.com/BroQiang/mdblog/releases/download/v1.1.0/mdblog.tar.gz)
这个我已经编译好了的版本到本地，然后执行：

```bash
tar xzvf mdblog.tar.gz
cd mdblog
./blog
```

然后浏览器访问 [http://localhost:8091](http://localhost:8091)

即可查看效果，只需要在 `resource/blog-docs` 目录中放入 markdown 文件即可使用。

> 需要注意， markdown 文件必须放在 blog-docs 的二级目录下才会被解析，如： resources/blog-docs/linux

如果想要修改，或者查看源码，请继续往下看文档。

## 需求

- go 语言为主的博客

- 博客放在我自己服务器上

- 博客展示的是 markdown 文档

- 博客的原始文件（ *.md 文件）放在 github 上，直接 github 直接更新文件，博客同步更新

按照上面的需求，最终就是一个静态动态结合的博客，完成后就是本博客了。
大体的思路是这样的，在 github 上创建一个专门用来保存 .md 文件的仓库 [blog-docs](https://github.com/BroQiang/blog-docs)，
将这个仓库 clone 到本博客项目的一个目录下（可以通过配置放在任意目录）。然后启动本博客系统，来解析 markdown 文本文件，用来展示。
当 [blog-docs](https://github.com/BroQiang/blog-docs) 文件仓库更新时，博客系统运行的内容自动更新。

## 依赖

### Go 语言的依赖

- go 版本是 `go1.12 linux/amd64` ， 其他版本未测试。

- [golang/dep](https://github.com/golang/dep) 用来管理 go 中的依赖

- [gin-gonic/gin](https://github.com/gin-gonic/gin) 引擎及路由

- [BurntSushi/toml](https://github.com/BurntSushi/toml) 用来处理 toml 格式的配置文件

- [russross/blackfriday](https://github.com/russross/blackfriday/releases/tag/v1.5.2) 用来将 markdown 文本转换成 HTML 使用的是 1.5.2 版本。

- [microcosm-cc/bluemonday](https://github.com/microcosm-cc/bluemonday) HTML清理程序， 配合 blackfriday 来处理 markdown 的转换。

### 前端的依赖

- nodejs 因为前端 scss 和 js 的编译是依赖于这个

- [laravel-mix](https://github.com/JeffreyWay/laravel-mix) 相当于一个 webpack 的包装器，用来管理前端静态资源，以前用 Laravel 框架的时候觉得这个用来管理前端的内容很顺手，试了下，原来可以单独使用， 就用了。

- [jquery](https://github.com/jquery/jquery) 因为 Bootstrap 依赖它。

- [bootstrap](https://github.com/twbs/bootstrap) 前端框架，用于页面展示。

- [sindresorhus/github-markdown-css](https://github.com/sindresorhus/github-markdown-css) markdown 文档的样式使用的这个 css

- [highlightjs/highlight.js](https://github.com/highlightjs/highlight.js) markdown 代码块高亮使用的此插件

- [阿里矢量图标库](https://www.iconfont.cn) 这个不错，按照需求定制，暂时我只用了 3 个图标

## 编译

编译需要编译两部分的内容，前端（webpack）和 Go 。

### Go 编译

编译前要保证已经正确安装了 Go 环境，并配置好环境变量， GOROOT 和 GOPATH

#### 安装 dep

这个是一个官方的依赖管理工作，见 [https://github.com/golang/dep](https://github.com/golang/dep) ， 有详细的介绍和安装步骤。

安装

```bash
go get -u github.com/golang/dep/cmd/dep
```

然后执行下面命令，查看结果，如果有下面类似的结果证明安装成功

```bash
$ dep version
dep:
 version     : devel
 build date  :
 git hash    :
 go version  : go1.12
 go compiler : gc
 platform    : linux/amd64
 features    : ImportDuringSolve=false
```

#### 安装博客

这是要指定目录，因为 GOPATH 的引入着这样做的，如果想要换成其他目录，需要同时修改包的导入路径。

```bash
# 创建项目保存的目录
mkdir $GOPATH/src/github.com/broqiang

# 进入到目录
cd $GOPATH/src/github.com/broqiang

# 下载项目
git clone https://github.com/BroQiang/mdblog.git

# 初始化项目（恢复依赖）
dep ensure
```

#### 修改配置文件

需要将项目 mdblog 目录下的 `config/config.example.toml` 复制一份

```bash
cd yourpath/config
cp config.example.toml config.toml
```

然后根据需要修改下配置，在配置文件中有详细的注释，说明每一个配置是用来做什么的。

#### 启动服务

__go run__ 方式：

这种方式一般就是开发时候使用，编译后的文件会生成在临时目录，所以一些静态文件会找不到，
所以这里要传入一个参数，告诉程序项目根目录在哪里。

```bash
# 可以指定目录（要绝对路径），$(pwd) 就是当前所在目录的绝对路径
go build -root=$(pwd)
```

此时就可以访问了，如配置文件默认的 8091 端口，在浏览器访问 [http://localhost:8091](http://localhost:8091)

__go build__ 编译方式:

```bash
# 编译
go build blog.go

# 运行
./blog

# 或通过绝对路径运行
/yourpath/blog
```

使用编译后的二进制文件启动服务的话，默认的项目根目录就是可执行文件所在的目录，如果目录结构没有改变的话可以不用传入
`-root` 参数，如果改变了，静态文件和可执行文件不在相同的目录，仍然需要在启动的时候传入 `-root=` 参数。

### 前端资源编译

如果只是使用本博客，不需要自定义 css 和 js 等，可以略过此步骤。

编译前要保证已经安装了 npm 环境，因为 npm 有些慢，我使用的 yarn 来初始化前端依赖（npm 仍然需要安装）

这里使用了个 laravel-mix 来统一管理静态资源（个人觉得比较方便，不喜欢的可以随意替换）。

配置文件：

```bash
config/webpack.mix.js
```

这里也没做什么特殊的事情，就是将所有的 scss 和 js 分别打包成一个单一的文件

打包前的文件保存在 `resources/assets` ， 分别存放了 scss 文件和 js 文件

编译后的文件存放在 `resources/static` 中，这个是项目真实使用的静态文件目录，对应的 http 方式的访问路径是 /public

## markdown 文件配置

markdown 的一些初始配置保存在 `config/categories.toml` 中，这个主要用于配置顶部的菜单导航，
可以有链接和需要解析的 markdown 文件仓库的二级目录，详细规则查看配置文件中的注释，有详细的备注。

## 配置自动启动

在 Linux 下，可以将启动程序加入到 Systemd 中，就可以通过 `systemctl` 命令来管理服务，可以方便的启动或停止

编写一个脚本，将下面内容保存到 /lib/systemd/system 目录中， 例如

用 vim 打开　`sudo vim /lib/systemd/system/mdblog.service` ，此时默认是一个空文件（如果没有使用过），　写入：

```bash
[Unit]
Description=blog - This is BroQiang blog
Documentation=https://broqiang.com
After=network.target remote-fs.target nss-lookup.target

[Service]
Type=simple
ExecStart=/www/web/mdblog/blog
ExecReload=/bin/kill $MAINPID && /www/web/mdblog/blog
ExecStop=/bin/kill -s QUIT $MAINPID
Restart=always

[Install]
WantedBy=multi-user.target
```

> 注意里面的项目目录 /www/web/mdblog/blog 要替换成实际的项目目录

此版本还没有实现热启动，启动和停止都是强制的，后面会陆续实现（不过这是个静态的博客，暂时没有数据交互，冷启动貌似也没大问题）。

添加完成上面的脚本就可以通过 Systemd 服务来管理了。

```bash
# 启动服务
sudo systemctl start mdblog

# 停止服务
sudo systemctl stop mdblog

# 配置开机自动启动
sudo systemctl enable mdblog

# 取消开机自动启动
sudo systemctl disable mdblog
```

## 配置　Nginx 转发

如果服务器只有一个　Go 服务，完全不需要　Nginx 了， 直接在配置文件中将端口改成 80 就可以了。
（不支持 ssl ， 也不计划支持，因为我自己实际使用是通过 Nginx 代理的，在 Nginx 中已经配置了证书）

在 Nginx 中添加一个虚拟主机（server 部分）

```nginx
server {
    listen       443; # 如果没有使用 https 可以将此端口改成 80

    # 如果没有使用 https 可以将下面三行关于证书的配置删除
    ssl on;
    ssl_certificate /etc/letsencrypt/live/broqiang.com/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/broqiang.com/privkey.pem;

    server_name broqiang.com www.broqiang.com;

    # 下面的配置是将 www.broqiang.com 的请求 301 到 broqiang.com ，如果只有一个域名可以去掉下面 3 行
    if ($host = 'www.broqiang.com'){
            rewrite ^/(.*)$ https://broqiang.com/$1 permanent;
    }

    # 访问日志， 这里和 mdblog 启动一个即可，要想关闭的话 access_log off;
    access_log /www/webLogs/go-blog_access.log;

    location / {
        # 将请求转发到 mdblog ，根据实际的监听地址修改
        proxy_pass http://127.0.0.1:8091;
    }
}
```

## 配置 github 钩子

在 `https://github.com/BroQiang/blog-docs` 项目中添加一个 Webhooks， 配置下面内容：

- Payload URL: `https://broqiang.com/webhook`

- Content type: 选择 `application/json`

- Secret: 自定义一个密钥，要和 `config.toml` 中的 secret 的值保持一直

钩子生效后， blog-docs 再 push 的时候 mdblog 就可以自动更新文档并显示了。

## 更新日志

### 2019-04-28 添加 github 钩子，自动同步 blog-docs 的文档
