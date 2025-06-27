---
title: "Golang 实现和 Laravel 类似的 .env 配置"
author: "BroQiang"
github_url: "https://broqiang.com"
created_at: 2018-06-24T23:20:51
updated_at: 2018-06-24T23:20:51
description: "Golang 实现和 Laravel 类似的 .env 配置"
tags: ["go"]
---

在使用 Laravel 的时候，会觉得它是优雅的，可能 Golang 主要的领域不是在 Web 开发中，也可能是它还太年轻，没有发现一个像 Laravel 一样优雅的框架，有一些以前用 Laravel 保留下的习惯，现在没有就会觉得很别扭。

Laravel 的 .env 配置很方便，可以通过这个配置，来设置本地和服务器端使用不同的配置，Golang 可以通过 [godotenv](https://github.com/joho/godotenv) 来实现和这个差不多的功能。

在根目录下创建一个 `.env` 文件，在里面写入

```bash
PORT=8888
```

这里只是定义一个端口号，用来简单的说明下怎么使用。

然后创建一个 `main.go` 文件

```go
package main

import (
	"flag"

	"log"

	"os"

	"net/http"

	"fmt"

	"github.com/joho/godotenv"
)

var port string

func init() {
	// 可以在命令行启动服务的时候通过 -port=端口号 ，来指定 web 服务的端口号
	// 如果没有指定会使用默认的 8080
	flag.StringVar(&port, "port", "8080", "The server listening port")
	flag.Parse()

	// 正常情况下会使用上面的端口号，可以通过 .env 中的配置来对端口号进行替换
	// 初始化 .env 的配置，将 .env 中的配置加载到 Go 的 env 环境中
	if err := godotenv.Load(".env"); err != nil {
		log.Fatalln(err)
		os.Exit(1)
	}

	// 可以通过 os 包的 Getenv 获取到 .env 中配置的端口号
	envPort := os.Getenv("PORT")
	// 如果 env 里面配置了，使用 env 的，如果没有配置，仍然使用默认的
	if len(envPort) > 0 {
		port = envPort
	}
}

func Hello(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html;charset=utf8")
	w.Write([]byte("<h1>Hello World</h1>"))
}

func main() {
	http.HandleFunc("/", Hello)

	log.Printf("Starting server on localhost:%s ...", port)

	log.Fatalln(
		http.ListenAndServe(
			fmt.Sprintf("localhost:%s", port), nil))
}
```

通过上面的示例，就可以看到，原本默认的端口是命令行输入的或者默认的 8080 端口，不过读取了 .env 配置文件后就被替换成了 .env 中的 8888 端口。

是不是有点像 Laravel 了，一般项目会有一个专门处理配置文件的包或文件，在里面再将它封装一下就很方便了。不过会不会影响性能还没有测试过，今天突然发现了 [godotenv](https://github.com/joho/godotenv) 这个包，觉得可以试一下，做一个记录。
