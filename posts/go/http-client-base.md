---
title: "HTTP 客户端 - 基础"
author: "BroQiang"
created_at: 2018-06-16T21:50:46
updated_at: 2018-06-16T21:50:46
---

Go 封装了非常方便的客户端方法，不再需要借助第三方包，就可以实现客户端的请求。本篇文章主要将的是 Get 、 Post 、PostForm 三个基础的方法。

文档已经在 [github](https://github.com/BroQiang/go-packages-study) 创建仓库，欢迎围观，欢迎修正。

## Get

[github](<(https://github.com/BroQiang/go-packages-study/tree/master/packages/net/http/Get)>)

这个方法非常的简单，就是从指定的地址获取到内容，可以简单的理解成我们用浏览器打开一个页面。

### Get 函数

#### 语法：

```go
Get(url string) (resp *Response, err error)
```

#### 参数：

- 字符串类型的 url 地址，需要注意的是这里要是完整地址，要加上 `http://` 或 `https://` 的地址

#### 返回值：

- \*Response 如果获取到了数据，会将数据保存在 Response 中

- error 如果请求数据的时候出现错误，会返回一个 error ，并将具体的错误记录到 error 中

### 另外一种方式 Get 方法

可以通过 client 结构体的 Get() 方法获取数据，其实两种方式是一样的，Get() 函数也是调用的结构体中的 Get() 方法。详细的使用可以见示例中的用法

### 示例：

`client.go`

```go
package main

import (
	"net/http"
	"log"
	"io/ioutil"
	"os"
	"fmt"
)

const url  = "https://broqiang.com"

func main() {
	// 方式一，直接通过 Get 函数
	resp, err := http.Get(url)
	ErrPrint(err)
	defer resp.Body.Close()

	// 拿到数据
	bytes, err := ioutil.ReadAll(resp.Body)
	ErrPrint(err)

	// 这里要格式化再输出，因为 ReadAll 返回的是字节切片
	fmt.Println("------------- 方法一 ---------------")
	fmt.Printf("%s",bytes)

	// 方式二，通过 client 结构体的 Get 方法
	client := new(http.Client)
	// 或
	client = &http.Client{}

	resp, err = client.Get(url)
	ErrPrint(err)
	defer resp.Body.Close()

	res, err := ioutil.ReadAll(resp.Body)
	ErrPrint(err)
	fmt.Println("\n\n\n------------- 方法二 ---------------")
	fmt.Printf("%s",res)
}

func ErrPrint(err error)  {
	if err != nil {
		log.Fatalln(err)
		os.Exit(1)
	}
}
```

## Post

[github](https://github.com/BroQiang/go-packages-study/tree/master/packages/net/http/Post)

这个方法也比较简单，就是通过 Post 方式和服务器端交互。

## POST 函数

### 语法：

```go
(url string, contentType string, body io.Reader) (resp *Response, err error)
```

### 参数：

- 字符串类型的 url 地址，需要注意的是这里要是完整地址，要加上 `http://` 或 `https://` 的地址

- 字符串类型的 contentType ，Post 提交数据的类型，常见的有下面 4 种：

  - `application/x-www-form-urlencoded` 不设置 enctype 属性的原生 form 表单提交方式。

  - `multipart/form-data` 上传文件时的数据提交方式，相当于 form 表单的 enctype 等于 multipart/form-data 。

  - `application/json` 用来告诉服务端消息主体是序列化后的 JSON 字符串。

  - `text/xml` 它是一种使用 HTTP 作为传输协议，XML 作为编码方式的远程调用规范，和 json 作用类型。

- 实现了 io.Reader 接口的数据。 如： 可以通过 `strings.NewReader()` 方法将普通字符串实现 io.Reader 接口。

### 返回值：

- \*Response 如果获取到了数据，会将数据保存在 Response 中

- error 如果请求数据的时候出现错误，会返回一个 error ，并将具体的错误记录到 error 中

### 另外一种方式 Post 方法

可以通过 client 结构体的 Post() 方法获取数据，其实两种方式是一样的，Post() 函数也是调用的结构体中的 Post() 方法。详细的使用可以见示例中的用法

### 示例：

这个示例要和服务器端交互数据，所以要有一个服务器端和一个客户端，服务器端见目录中的服务器端内容（如果你很快就发现了这个文档，可能还没来得及写服务器端）。

### 服务器端代码

`server.go`

可以通过 [http://localhost:8080/hello](http://localhost:8080/hello) 来访问服务器端，但是要是 Post 方式才可以。测试的时候注意要先启动服务器端，再启动客户端。

```go
package main

import (
	"net/http"
	"fmt"
	"log"
)

func main() {
	http.HandleFunc("/hello", func(w http.ResponseWriter, req *http.Request) {
		if req.Method == "POST" {
			fmt.Fprintf(w,"Hello %s\n",req.FormValue("name"))
			return
		}

		http.NotFound(w,req)
	})

	log.Fatalf("%v",http.ListenAndServe("localhost:8080",nil))
}
```

### 客户端代码

`client.go`

```go
package main

import (
	"log"
	"os"
	"net/http"
	"strings"
	"fmt"
	"io/ioutil"
	"io"
)

const url  = "http://localhost:8080/hello"

func main() {
	// 方式一，直接通过 Post 函数
	fmt.Println("------------- 方法一 ---------------")
	resp, err := http.Post(url,"application/x-www-form-urlencoded",
		strings.NewReader("name=Bro Qiang"))
	ErrPrint(err)
	defer resp.Body.Close()

	DataPrint(resp.Body)

	// 方式二，通过 client 结构体中的 Post 方法
	fmt.Println("------------- 方法二 ---------------")
	client := &http.Client{}
	resp, err = client.Post(url,"application/x-www-form-urlencoded",
		strings.NewReader("name=New Bro Qiang"))
	ErrPrint(err)
	defer resp.Body.Close()

	DataPrint(resp.Body)
}

func DataPrint(body io.ReadCloser) {
	// 拿到数据
	bytes, err := ioutil.ReadAll(body)
	ErrPrint(err)

	// 这里要格式化再输出，因为 ReadAll 返回的是字节切片
	fmt.Printf("%s",bytes)
}

func ErrPrint(err error)  {
	if err != nil {
		log.Fatalln(err)
		os.Exit(1)
	}
}
```

## PostForm

[github](https://github.com/BroQiang/go-packages-study/tree/master/packages/net/http/PostForm)

这个方法是通过 Post 方式向服务器端提交一个 form 表单，和 html 中的 form 表单提交类型。

### PostForm 函数

#### 语法：

```go
PostForm(url string, data url.Values) (resp *Response, err error)
```

#### 参数：

- 字符串类型的 url 地址，需要注意的是这里要是完整地址，要加上 `http://` 或 `https://` 的地址

- url.Values 类型的数据, 它的类型实际是一个 `map[string][]string`

#### 返回值：

- \*Response 如果获取到了数据，会将数据保存在 Response 中

- error 如果请求数据的时候出现错误，会返回一个 error ，并将具体的错误记录到 error 中

### 另外一种方式 PostForm 方法

可以通过 client 结构体的 PostForm() 方法获取数据，其实两种方式是一样的，PostForm() 函数也是调用的结构体中的 PostForm() 方法。详细的使用可以见示例中的用法

#### 示例：

创建一个服务器端，用来接收 from 表单，并将数据在控制台打印，并给客户端返回是否成功。

创建一个客户端，用来向服务器端发送 form 表单。

> 注意： 客户端的 resp.Body 不要忘记关闭。

#### 服务器端

`server.go`

```go
package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/form", MyForm)

	log.Fatalf("%v",
		http.ListenAndServe("localhost:8080", nil))
}

func MyForm(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
		return
	}

	formData := r.Form
	log.Printf("收到的数据： %v", formData)

	fmt.Fprintf(w, "提交成功")
}
```

#### 客户端代码

`client.go`

```go
package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

const url = "http://localhost:8080/form"

func main() {
	data := map[string][]string{"name": {"Bro Qiang"}, "gender": {"male"}}

	// 方法一：PostForm 函数
	resp, err := http.PostForm(url, data)
	ErrPrint(err)
	defer resp.Body.Close()

	DataPrint(resp.Body)

	// 方法二：client 结构体的 PostForm 方法
	client := &http.Client{}
	resp, err = client.PostForm(url, data)
	ErrPrint(err)
	defer resp.Body.Close()

	DataPrint(resp.Body)
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
