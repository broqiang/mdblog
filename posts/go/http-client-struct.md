---
title: "HTTP 客户端 - 使用 Client 类型"
author: "BroQiang"
github_url: "https://broqiang.com"
created_at: 2018-06-17T01:35:15
updated_at: 2018-06-17T01:35:15
description: "主要介绍了 Client 类型以及 Do 和 Head 的使用。"
tags: ["go", "http"]
---

这章主要介绍了 Client 类型以及 Do 和 Head 的使用。

## client 类型

[示例代码](https://github.com/BroQiang/go-packages-study/tree/master/packages/net/http/client)

Client 类型代表 HTTP 客户端。它的零值（ DefaultClient ）是一个可用的使用 DefaultTransport 的客户端。

Client 的 Transport 字段一般会含有内部状态（缓存 TCP 连接），因此 Client 类型值应尽量被重用而不是每次需要都创建新的。 Client 类型值可以安全的被多个 go 程同时使用。

Client 类型的层次比 RoundTripper 接口（如 Transport ）高，还会管理 HTTP 的 cookie 和重定向等细节。

`Client` 类型的结构体

```go
type Client struct {
    // Transport 指定执行独立、单次 HTTP 请求的机制。
    // 如果 Transport 为 nil，则使用 DefaultTransport 。
    Transport RoundTripper

    // CheckRedirect 指定处理重定向的策略。
    // 如果 CheckRedirect 不为 nil，客户端会在执行重定向之前调用本函数字段。
    // 参数 req 和 via 是将要执行的请求和已经执行的请求（切片，越新的请求越靠后）。
    // 如果 CheckRedirect 返回一个错误，本类型的 Get 方法不会发送请求 req，
    // 而是返回之前得到的最后一个回复和该错误。（包装进 url.Error 类型里）
    //
    // 如果CheckRedirect为nil，会采用默认策略：连续10此请求后停止。
    CheckRedirect func(req *Request, via []*Request) error

    // Jar 指定 cookie 管理器。
    // 如果Jar为nil，请求中不会发送 cookie ，回复中的 cookie 会被忽略。
    Jar CookieJar

    // Timeout 指定本类型的值执行请求的时间限制。
    // 该超时限制包括连接时间、重定向和读取回复主体的时间。
    // 计时器会在 Head 、 Get 、 Post 或 Do 方法返回后继续运作并在超时后中断回复主体的读取。
    //
    // Timeout 为零值表示不设置超时。
    //
    // Client 实例的 Transport 字段必须支持 CancelRequest 方法，
    // 否则 Client 会在试图用 Head 、 Get 、 Post 或 Do 方法执行请求时返回错误。
    // 本类型的 Transport 字段默认值（ DefaultTransport ）支持 CancelRequest 方法。
    Timeout time.Duration
}
```

在此目录下初始化了一个 Client ，用来配合后面测试使用。

```go
package myclient

import (
	"net/http"
)

var Client = &http.Client{}
```

## Do 方法

[示例代码](https://github.com/BroQiang/go-packages-study/tree/master/packages/net/http/client/Do)

Do 方法发送请求，返回 HTTP 回复。它会遵守客户端 c 设置的策略（如重定向、cookie、认证）。

如果客户端的策略（如重定向）返回错误或存在 HTTP 协议错误时，本方法将返回该错误；如果回应的状态码不是 2xx，本方法并不会返回错误。

如果返回值 err 为 nil，resp.Body 总是非 nil 的，调用者应该在读取完 resp.Body 后关闭它。如果返回值 resp 的主体未关闭，c 下层的 RoundTripper 接口（一般为 Transport 类型）可能无法重用 resp 主体下层保持的 TCP 连接去执行之后的请求。

请求的主体，如果非 nil，会在执行后被 c.Transport 关闭，即使出现错误。

一般应使用 Get 、 Post 或 PostForm 方法就可以代替 Do 方法，其实它们最终执行的也是 Do ，只不过做了一些包装。

当有一些比较特殊的，上面的三种方式不能满足时，就要自己初始化 Request ，然后调用 Do 方法。

### 语法：

```go
func (c *Client) Do(req *Request) (resp *Response, err error)
```

### 参数：

- \*Request 可以通过这个参数来自定义 Request

### 返回值：

- \*Response 如果获取到了数据，会将数据保存在 Response 中

- error 如果请求数据的时候出现错误，会返回一个 error ，并将具体的错误记录到 error 中

### 示例：

详细的使用请看示例代码，已经在里面写了详细的注释。

`myclient.Client` 是在上一级目录的 client.go 文件中初始化的。

### 服务器端

`server.go`

```go
package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

func main() {
	http.HandleFunc("/", TestDo)

	log.Fatalf("%v", http.ListenAndServe("localhost:8080", nil))
}

func TestDo(w http.ResponseWriter, req *http.Request) {
	// 休眠 2 秒再处理，一会来测试超时
	time.Sleep(2e9)

	// 验证 csrf
	if req.Header.Get("_csrf") != "123456" {
		http.Error(w, "无效的 csrf token", 400)
		return
	}

	// 获取 cookie
	cookie := req.Header.Get("Cookie")
	if cookie == "" {
		http.Error(w, "请登陆后再操作", 401)
		return
	}

	// 获取用户名，简单的利用 strings 去截取字符串，只是个简单示例，没有考虑那么多可能性。
	fmt.Printf("%v\n", cookie)
	index := strings.Index(cookie, "=")
	name := cookie[index+1:]
	fmt.Printf("%v\n", name)
	if name != "BroQiang" {
		http.Error(w, "当前用户没有权限操作", 401)
		return
	}

	io.WriteString(w, "Hello "+name)
}
```

### 客户端

`client.go`

```go
package main

import (
	"log"
	"net/http"
	"os"

	"io/ioutil"

	"fmt"
	"io"

	"github.com/broqiang/go-packages-study/packages/net/http/client"
)

func main() {
	// 测试自定义的 Get 函数
	GetData()

	// 当然，我们不是为了自定义 Get 、 Post 方法，才去实现 Do 方法的调用
	// 是因为有一些特殊的需求，比如要携带 cooki ， 连接超时时间设置等。
	MyConnection()
}

func MyConnection() {
	log.Println("--------------- 开始执行自定义的 Do 方法 ----------------")
	// 正常执行自定义的请求
	myDo()

	log.Println("--------------- 开始执行超时后的自定义的 Do 方法 ----------------")
	// 给 Client 设置一个超时，然后测试效果
	// 这里只是简单的验证超时，没有管 Transport 的设置
	myclient.Client.Timeout = 1e9
	myDo()
	// 可以看到会有一个链接超时的错误，如下面这样
	// Get http://localhost:8080: net/http: request canceled (Client.Timeout exceeded while awaiting headers)
}

func myDo() {
	// 先自定义一个 Request
	req, err := http.NewRequest("GET", "http://localhost:8080", nil)
	ErrPrint(err)

	// 设置一个 csrf token, 服务器端会去验证这个
	req.Header.Set("_csrf", "123456")
	// 设置一个 Cookie, 服务器端会去验证这个
	req.Header.Set("Cookie", "name=BroQiang")
	// 一会写完服务器端可以尝试分别将上面两行注释去测试

	resp, err := myclient.Client.Do(req)
	ErrPrint(err)
	defer resp.Body.Close()

	DataPrint(resp.Body)
	// 打印下状态码，看下效果
	fmt.Printf("返回的状态码是： %v\n", resp.StatusCode)
	fmt.Printf("返回的信息是： %v\n", resp.StatusCode)
}

func GetData() {
	log.Println("------------ 自定义的 Get 方法获取数据 ----------")
	url := "https://broqiang.com"
	resp, err := MyGet(url)
	ErrPrint(err)
	defer resp.Body.Close()

	DataPrint(resp.Body)

}

// 实现和 Get 函数一样的功能，其实这个就是源码改的
// Post 和 PostForm 就演示了，实现原理差不多，最终夜都是要调用 client.Do
func MyGet(url string) (resp *http.Response, err error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	return myclient.Client.Do(req)
}

func DataPrint(body io.ReadCloser) {
	// 拿到数据
	bytes, err := ioutil.ReadAll(body)
	ErrPrint(err)

	// 这里要格式化再输出，因为 ReadAll 返回的是字节切片
	fmt.Printf("%s\n", bytes)
}

func ErrPrint(err error) {
	if err != nil {
		log.Fatalln(err)
		os.Exit(1)
	}
}
```

## Head

[示例代码](https://github.com/BroQiang/go-packages-study/tree/master/packages/net/http/client/Head)

Head 向指定的 URL 发出一个 HEAD 请求，如果回应的状态码如下，Head 会在调用 client.CheckRedirect 后执行重定向。

就是 Head 只会请求 Head 部分，可以快速的返回。比如你写了一个爬虫，可以先通过 Head 测试下，是否是 200 的，就根据返回的状态码来做处理，Head 请求没有 Body 部分，速度会快很多。

除了没有返回的 Body ，基本上用起来和 Get 差不多。

### Head 函数

#### 语法：\*\*

```go
Head(url string) (resp *Response, err error)
```

#### 参数：

- 字符串类型的 url 地址，需要注意的是这里要是完整地址，要加上 `http://` 或 `https://` 的地址

#### 返回值：

- \*Response 如果获取到了数据，会将数据保存在 Response 中

- error 如果请求数据的时候出现错误，会返回一个 error ，并将具体的错误记录到 error 中

### 另外一种方式 Head 方法

可以通过 client 结构体的 Head() 方法获取数据，其实两种方式是一样的，Head() 函数也是调用的结构体中的 Head() 方法。详细的使用可以见示例中的用法

#### 示例：

`client.go`

```go
package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func main() {
	fmt.Println("----------访问一个正常的 URL ----------")
	MyHead("https://broqiang.com")
	// 结果：返回的状态码是： 200

	fmt.Println("----------访问一个正常的 URL ----------")
	// 可以将之前测试 Do 用的服务器端启动，因为这个服务器端有限制，
	// 可以通过返回的 code 来直到访问的状况。
	MyHead("http://localhost:8080")
	// 结果： 返回的状态码是： 400
}

func MyHead(url string) {
	resp, err := http.Head(url)
	ErrPrint(err)

	defer resp.Body.Close()

	bytes, err := ioutil.ReadAll(resp.Body)
	ErrPrint(err)

	fmt.Printf("%q\n", bytes)
	fmt.Printf("返回的状态码是： %d\n", resp.StatusCode)
}

func ErrPrint(err error) {
	if err != nil {
		log.Fatalln(err)
		os.Exit(1)
	}
}
```
