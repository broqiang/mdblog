---
title: "httprouter 源码分析"
author: "BroQiang"
created_at: 2019-04-17T01:05:48
updated_at: 2019-04-17T01:05:48
---

关于 [httprouter](https://github.com/julienschmidt/httprouter) 本身就不过多说了，可以直接去查看源码及 README 。

这个包相对还是比较简单了，只有几个文件，并且除了标准库没有外部的依赖。难理解的就是基数树，需要算法基础。

抛砖引玉，有不对的地方望指出，我及时修改。

## 入口

使用的是代码追踪的方式，可以从官方给的 demo 来入手：

```go
package main

import (
    "fmt"
    "github.com/julienschmidt/httprouter"
    "net/http"
    "log"
)

func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
    fmt.Fprint(w, "Welcome!\n")
}

func Hello(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
    fmt.Fprintf(w, "hello, %s!\n", ps.ByName("name"))
}

func main() {
    router := httprouter.New()
    router.GET("/", Index)
    router.GET("/hello/:name", Hello)

    log.Fatal(http.ListenAndServe(":8080", router))
}
```

可以看到，这个 demo 还是比较简单的， 主要就做了几件事：

- 定义了两个业务处理的函数（虽然只做了简单的字符串输出）。

- 初始化路由， httprouter.New()

- 将上面的两个函数注册到路由中

- 使用 httprouter 的 handler 启动服务

下面就从 main 函数开始看，通过 `New()` 函数初始化了一下 httprouter ，然后调用了两次 GET 方法，分别传入了 Index 和 Hello 的业务处理函数。

下面就追踪到 New() 中，看看它到底做了什么，这个函数也很简单， 只是初始化了一下 Router 结构，将几个参数的默认值设置为 true ，并返回了 \*Router （指针）。

```go
func New() *Router {
    return &Router{
        RedirectTrailingSlash:  true,
        RedirectFixedPath:      true,
        HandleMethodNotAllowed: true,
        HandleOPTIONS:          true,
    }
}
```

这个 Route 是什么？ 里面的参数又代表了什么，追踪到 Router 去看一下。

## Router struct

先看看这个结构体里面包含了什么字段，顺便将注释翻译一下,这里就没有保留注释的原文，如果想要了解原文，可以直接去 [源码](https://github.com/julienschmidt/httprouter/blob/master/router.go#L110) 中查看即可。

```go
// Router 是一个 http.Handler 可以通过定义的路由将请求分发给不同的函数
type Router struct {
    trees map[string]*node

    // 这个参数是否自动处理当访问路径最后带的 /，一般为 true 就行。
    // 例如： 当访问 /foo/ 时， 此时没有定义 /foo/ 这个路由，但是定义了
    // /foo 这个路由，就对自动将 /foo/ 重定向到 /foo (GET 请求
    // 是 http 301 重定向，其他方式的请求是 http 307 重定向）。
    RedirectTrailingSlash bool

    // 是否自动修正路径， 如果路由没有找到时，Router 会自动尝试修复。
    // 首先删除多余的路径，像 ../ 或者 // 会被删除。
    // 然后将清理过的路径再不区分大小写查找，如果能够找到对应的路由， 将请求重定向到
    // 这个路由上 ( GET 是 301， 其他是 307 ) 。
    RedirectFixedPath bool

    // 用来配合下面的 MethodNotAllowed 参数。
    HandleMethodNotAllowed bool

    // 如果为 true ，会自动回复 OPTIONS 方式的请求。
    // 如果自定义了 OPTIONS 路由，会使用自定义的路由，优先级高于这个自动回复。
    HandleOPTIONS bool


    // 路由没有匹配上时调用这个 handler 。
    // 如果没有定义这个 handler ，就会返回标准库中的 http.NotFound 。
    NotFound http.Handler

    // 当一个请求是不被允许的，并且上面的 HandleMethodNotAllowed 设置为 ture 的时候，
    // 如果这个参数没有设置，将使用状态为 with http.StatusMethodNotAllowed 的 http.Error
    // 在 handler 被调用以前，为允许请求的方法设置 "Allow" header 。
    MethodNotAllowed http.Handler

    // 当出现 panic 的时候，通过这个函数来恢复。会返回一个错误码为 500 的 http error
    // (Internal Server Error) ，这个函数是用来保证出现 painc 服务器不会崩溃。
    PanicHandler func(http.ResponseWriter, *http.Request, interface{})
}
```

现在可以看到， Router 这个结构就是 httprouter 的一个核心的部分，这里定义了路由的一些初始配置，基本通过注释就可以知道它们是做什么用的，上面还有个 trees 是这个结构最核心的内容，竟然没注释（原文就没有注释），这里面保存的就是注册的路由，按照 method (POST,GET ...) 方法分开，每一个方法对应一个基数树（radix tree）。这个树比较复杂，暂时先不去理会，简单的理解成将路由保存进来即可。

在 Router 结构和 New() 函数之间还有一行比较有意思的写法：

```go
var _ http.Handler = New()
```

这个其实就是通过 New() 函数初始化一个 Router ，并且指定 Router 实现了 http.Handler 接口 ，如果 New() 没有实现 http.Handler 接口，在编译的时候就会报错了。这里只是为了验证一下， New() 函数的返回值并不需要，所以就把它赋值给 \_ ，相当于是给丢弃了。

到这里，我们就可以看到了， Router 也是基于 http.Handler 做的实现，如果要实现 http.Handler 接口，就必须实现 `ServeHTTP(w http.ResponseWriter, req *http.Request)` 这个方法，下面就可以去追踪下 ServerHTTP 都做了什么。

## ServerHTTP

这个代码比较长，就将分析的步骤写在注释中了，关于基数树放在最后说明。

```go
// ServeHTTP makes the router implement the http.Handler interface.
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
    // 在最开始，就做了 panic 的处理，这样可以保证业务代码出现问题的时候，不会导致服务器崩溃
    // 这里就是简单的在浏览器返回一个 500 错误，如果想在 panic 的时候自己处理
    // 只需要将 r.PanicHandler = func(http.ResponseWriter, *http.Request, interface{})
    // 重写， 添加上自己的处理逻辑即可。
    if r.PanicHandler != nil {
        defer r.recv(w, req)
    }

    // 通过 Request 获取当前请求的路径
    path := req.URL.Path

    // 到基数树中去查找匹配的路由
    // 首先看是否有请求的方法（比如 GET）是否存在路由，如果存在就继续寻找
    if root := r.trees[req.Method]; root != nil {

        // 如果路由匹配上了，从基数树中将路由取出
        if handle, ps, tsr := root.getValue(path); handle != nil {
            // 获取的是一个函数，签名 func(http.ResponseWriter, *http.Request, Params)
            // 这个就是我们向 httprouter 注册的函数
            handle(w, req, ps)
            // 处理完成后 return ，此次生命周期就结束了
            return

            // 当没有找到的时候，并且请求的方法不是 CONNECT 并且 路径不是 / 的时候
        } else if req.Method != "CONNECT" && path != "/" {
            // 这里就要做重定向处理， 默认是 301
            code := 301 // Permanent redirect, request with GET method

            // 如果请求的方式不是 GET 就将 http 的响应码设置成 307
            if req.Method != "GET" {
                // Temporary redirect, request with same method
                // As of Go 1.3, Go does not support status code 308.
                // 看上面的注释的意思貌似是作者想使用 308（永久重定向），但是 Go 1.3
                // 不支持 308 ，所以使用了一个临时重定向。 为什么不能用 301 呢？
                // 因为 308 和 307 不允许请求方法从 POST 更改为 GET
                code = 307
            }

            // tsr 返回值是一个 bool 值，用来判断是否需要重定向, getValue 返回来的
            // RedirectTrailingSlash 这个就是初始化时候定义的，只有为 true 才会处理
            if tsr && r.RedirectTrailingSlash {
                // 如果 path 的长度大于 1，只有大于 1 才会出现这种情况，如 p/，path/
                // 并且路径的最后是 / 的时候将最后的 / 去除
                if len(path) > 1 && path[len(path)-1] == '/' {
                    req.URL.Path = path[:len(path)-1]
                } else {
                    // 如果不是 / 结尾的，给结尾添加一个 /
                    // 假设定义了一个路由 '/foo/' ，并且没有定义路由 /foo ,
                    // 实际访问的是 '/foo'， 将 /foo 重定向到 /foo/
                    req.URL.Path = path + "/"
                }

                // 将处理过的路由重定向， 这个是一个 http 标准包里面的方法
                http.Redirect(w, req, req.URL.String(), code)
                return
            }

            // Try to fix the request path
            // 路由没有找到，重定向规则也不符合，这里会尝试修复路径
            // 需要在初始化的时候定义 RedirectFixedPath 为 true，允许修复
            if r.RedirectFixedPath {
                // 这里就是在处理 Router 里面说的，将路径通过 CleanPath 方法去除多余的部分
                // 并且 RedirectTrailingSlash 为 ture 的时候，去匹配路由
                // 比如： 定义了一个路由 /foo , 但实际访问的是 ////FOO ，就会被重定向到 /foo
                fixedPath, found := root.findCaseInsensitivePath(
                    CleanPath(path),
                    r.RedirectTrailingSlash,
                )
                if found {
                    req.URL.Path = string(fixedPath)
                    http.Redirect(w, req, req.URL.String(), code)
                    return
                }
            }
        }
    }

    // 路由也没有匹配，重定向也没有找到，修复后仍然没有匹配，就到了这里

    // 如果没有任何路由匹配，请求方式又是 OPTIONS， 并且允许响应 OPTIONS
    // 就会给 Header 设置一个 Allow 头，返回去
    if req.Method == "OPTIONS" && r.HandleOPTIONS {
        // Handle OPTIONS requests
        if allow := r.allowed(path, req.Method); len(allow) > 0 {
          w.Header().Set("Allow", allow)
            return
        }
    } else {
        // 如果不是 OPTIONS 或者 不允许 OPTIONS 时
        // Handle 405
        // 如果初始化的时候 HandleMethodNotAllowed 为 ture
        if r.HandleMethodNotAllowed {
            // 返回 405 响应，通过 allowed() 方法来处理 405 时 allow的值。
            // 大概意思是这样的，比如定义了一个 POST 方法的路由 POST("/foo",...)
            // 但是调用却是通过 GET 方式，这是就会给调用者返回一个包含 POST 的 405
            if allow := r.allowed(path, req.Method); len(allow) > 0 {
                w.Header().Set("Allow", allow)
                if r.MethodNotAllowed != nil {
                    // 这里默认就会返回一个字符串 Method Not Allowed
                    // 可以自定义 r.MethodNotAllowed = 自定义的 http.Handler 实现
                    // 来按照需求响应， allowd 方法也不太难，就是各种判断和查找，就不分析了。
                    r.MethodNotAllowed.ServeHTTP(w, req)
                } else {
                    http.Error(w,
                        http.StatusText(http.StatusMethodNotAllowed),
                        http.StatusMethodNotAllowed,
                    )
                }
                return
            }
        }
    }

    // Handle 404
    // 如果什么的没有找到，就只会返回一个 404 了
    // 如果定义了 NotFound ，就会调用自定义的
    if r.NotFound != nil {
        // 如果需要自定义，在初始化之后给 NotFound 赋一个值就可以了。
        // 可以简单的通过 http.HandlerFunc 包装一个 handler ，例如：
        // router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        //         w.Write([]byte("什么都没有找到"))
        // })
        //
        r.NotFound.ServeHTTP(w, req)
    } else {
        // 如果没有定义，就返回一个 http 标准库的  NotFound 函数，返回一个 404 page not found
        http.NotFound(w, req)
    }
```

到这里基本也就清除了它的基本调用的流程， 就是定义一个 Router 结构，并且使它实现 http.Handler 接口， 也就是添加一个 ServeHTTP 方法。 然后在 ServeHTTP 中去做路由的处理。

## 路由设置

现在为止，看到的都是路由的调用，那些路由又是哪里来的呢？ 其实如果看了 router.go 这个文件，就可以在里面发现，GET， POST， PUT 等方法，就是这些方法来设置的。

还记得开始的时候，引用了一个官方的 demo，里面就有两个设置路由的例子：

```go
...


func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
    fmt.Fprint(w, "Welcome!\n")
}

...

router.GET("/", Index)

```

下面就追踪这个 GET 方法，看看它到底做了什么：

```go
func (r *Router) GET(path string, handle Handle) {
    r.Handle("GET", path, handle)
}

...

// DELETE is a shortcut for router.Handle("DELETE", path, handle)
func (r *Router) DELETE(path string, handle Handle) {
    r.Handle("DELETE", path, handle)
}

...

```

可以看到，从 GET 方法一直到 DELETE 方法，几乎长得都一样，也没做其它事情，就是调用了一下 Handle 方法，继续追踪 Handle 方法。

```go
// 通过给定的路径和方法，注册一个新的 handle 请求，对于 GET, POST, PUT, PATCH 和 DELETE
// 请求，相对应的方法，直接调用这个方法也是可以的，只需要多传入第一个参数即可。
// 官方给的建议是： 在批量加载或使用非标准的自定义方法时候使用。
//
// 它有三个参数，第一个是请求方式（GET，POST……）， 第二个是路由的路径， 前两个参数都是字符串类型
// 第三个参数就是我们要注册的函数，比如上面的 Index 函数，可以追踪一下，在这个文件的最上面
// 有 Handle 的定义，它是一个函数类型，只要和这个 Handle 的签名一直的都可以作为参数
func (r *Router) Handle(method, path string, handle Handle) {

    // 验证路径是否是 / 开头的，如： /foo 就是可以的， foo 就会 panic，编译的时候机会出错
    if path[0] != '/' {
        panic("path must begin with '/' in path '" + path + "'")
    }

    // 因为 map 是一个指针类型，必须初始化才可以使用，这里做一个判断，
    // 如果从来没有注册过路由，要先初始化 tress 属性的 map
    if r.trees == nil {
        r.trees = make(map[string]*node)
    }

    // 因为路由是一个基数树，全部是从根节点开始，如果第一次调用注册方法的时候跟是不存在的，
    // 就注册一个根节点， 这里是每一种请求方法是一个根节点，会存在多个树。
    root := r.trees[method]

    // 根节点存在就直接调用，不存在就初始化一个
    if root == nil {
        // 需要注意的是，这里使用的 new 来初始化，所以 root 是一个指针。
        root = new(node)
        r.trees[method] = root
    }

    // 向 root 中添加路由，树的具体操作在后面单独去分析。
    root.addRoute(path, handle)
}
```

现在，路由就已经注册，可以使用了。 看一下上面函数的第三个参数 Handle ，追踪一下

```go
type Handle func(http.ResponseWriter, *http.Request, Params)
```

可以看到，这个签名和 `http.HandlerFunc` 的签名比较像，用途其实也是一样的，只不过多了一个 Params 参数。这个参数就是用来支持路径中的参数绑定的， 比如： `/hello/:name` ,可以通过 `Params.ByName("name")` 来获取路径绑定参数的值。其实这个 Params 就是一个 Param 的切片，
也可以自己将所有的值遍历出来，追踪过去看一下就很容易明白了：

```go
// Parram 是一个 URL 参数，包含了一个 Key 和 一个 Value 。
// Key 就是注册路由时候的使用的字符串，如： "/hello/:name" ，这个 Key 就是 name
// Value 就是访问路径中获取到的值， 如： localhost:8080/hello/BroQiang ，
// 此时 Value 就是 BroQiang
type Param struct {
    Key   string
    Value string
}

// Params 就是一个 Param 的切片，这样就可以看出来， URL 参数可以设置多个了。
// 它是在 tree 的 GetValue() 方法调用的时候设置的，一会分析树的时候可以看到。
// 这个切片是有顺序的，第一个设置的参数就是切片的第一个值，所以通过索引获取值是安全的
type Params []Param

// ByName 返回传入的 name 的值，如果没有找到就会返回一个空字符串
// 这里就是做了个循环，一担找到就将值返回。
func (ps Params) ByName(name string) string {
    for i := range ps {
        if ps[i].Key == name {
            return ps[i].Value
        }
    }
    return ""
}
```

## 定义静态文件目录

到此时，我们会发现整个 `router.go` 文件貌似都已经追踪完了，还整下两个函数没有提到： 和 `Lookup` 和 `ServeFiles` 。

`Lookup` 暂时没有发现哪里用到了，一会查看其他文件的时候看看有没有用到的地方。先留个坑，最后再填。

下面来看看 `ServeFiles` 这个方法。

```go
// 这个是用来定义静态文件的，比如 js, css, 图片等
// path 是定义的 URL 路径，必须是 /*filepath 这种格式
// root 对应的是系统中的目录或文件系统，因为这个函数其实是使用的 http 标准库中的 FileServer
// 来实现的文件服务，所以 root 参数要传入一个 http.FileSystem 类型的参数，可以通过
// http.Dir("/etc/...") 这种方式转换一下，因为 Dir 实现了 http.FileSystem 接口，
// 直接使用它来转换就行了。
//
// 示例： 比如定义了路由 router.ServeFiles("/public/*filepath", http.Dir("/etc"))
// 此时访问 localhost:8080/public/passwd ， 就可以将 /etc/passwd 文件的内容显示了
func (r *Router) ServeFiles(path string, root http.FileSystem) {
    // 就是因为这三行代码的处理，是验证路径是否 /xxx/*filepath 这种结构定义
    // 因为标准库要求必须这样设置路径，这里也就这样了,
    // 个人还真没想到为什么要这样处理，应该有什么意义吧
    if len(path) < 10 || path[len(path)-10:] != "/*filepath" {
        panic("path must end with /*filepath in path '" + path + "'")
    }

    fileServer := http.FileServer(root)

    r.GET(path, func(w http.ResponseWriter, req *http.Request, ps Params) {
        req.URL.Path = ps.ByName("filepath")
        fileServer.ServeHTTP(w, req)
    })
}
```

这个方法其实最终也是调用的标准库，并且不太灵活，看需求，不用也是可以的，比如需要给静态文件设置缓存，调用这个方法就没法实现了，不过可以自己包装一下 http.FileServer 就可以了，详细参考这个 [Issues#40](https://github.com/julienschmidt/httprouter/issues/40)

到这里 httprouter 的基本路由处理已经读完，在不考虑性能的情况下，完全可以模仿这个来实现一个自己的路由了。如果是一个小的服务，没有几个路由，其实完全不用考虑，没必要使用树，直接一个 map 就可以了。标准库的 ServeMux 中的 muxEntry 也只有简单的两个字段（所以支持的功能比较少）。

## Radix tree

这个就是这个路由号称最快的关键部分（这个没去核实过，也不用去纠结）。开始读源代码前，先了解下原理（这个就是之前一直没去读源码的原因），详细见官方给的说明： [How does it work?](https://github.com/julienschmidt/httprouter#how-does-it-work)

这里简单的解读下，路由使用了一个有共同前缀的一个树结构，这个树就是一个压缩前缀树（ compact [prefix tree](https://zh.wikipedia.org/wiki/Trie) ） 或者就叫基数树（ [Radix tree](https://zh.wikipedia.org/wiki/%E5%9F%BA%E6%95%B0%E6%A0%91) ）。也就是具有共同前缀的节点拥有相同的父节点，官方给出的一个 GET 请求的例子，通过示例就比较好理解了：

图示，直接拔过来的。

```go
Priority   Path             Handle
9          \                *<1>
3          ├s               nil
2          |├earch\         *<2>
1          |└upport\        *<3>
2          ├blog\           *<4>
1          |    └:post      nil
1          |         └\     *<5>
2          ├about-us\       *<6>
1          |        └team\  *<7>
1          └contact\        *<8>


// 这个图相当于注册了下面这几个路由
GET("/search/", func1)
GET("/support/", func2)
GET("/blog/:post/", func3)
GET("/about-us/", func4)
GET("/about-us/team/", func5)
GET("/contact/", func6)

```

通过上面的示例可以看出：

- \*<数字> 代表一个 handler 函数的内存地址（指针）

- search 和 support 拥有共同的父节点 s ，并且 s 是没有对应的 handle 的， 只有叶子节点（就是最后一个节点，下面没有子节点的节点）才会注册 handler 。

- 从根开始，一直到叶子节点，才是路由的实际路径。

- 路由搜索的顺序是从上向下，从左到右的顺序，为了快速找到尽可能多的路由，包含子节点越多的节点，优先级越高。

### node

大概了解下，开始读代码，详细去看看。一样是先找一个入口，上面在分析的时候已经留了好几个坑， 可以慢慢填了。

首先在 router.go 中找到

```go
...
type Router struct {
    trees map[string]*node

    ...
}
```

这里是一个 map ， 每一种请求方式（GET，POST ……) 单独管理一颗树，官方说这样比每个节点中去保存方法节省空间，并且查找的时候速度会更快。

接下来看看这个 node 是什么：

```go
type node struct {
    // 当前节点的 URL 路径
    // 如上面图中的例子的首先这里是一个 /
    // 然后 children 中会有 path 为 [s, blog ...] 等的节点
    // 然后 s 还有 children node [earch,upport] 等，就不再说明了
    path      string

    // 判断当前节点路径是不是含有参数的节点, 上图中的 :post 的上级 blog 就是wildChild节点
    wildChild bool

    // 节点类型: static, root, param, catchAll
    // static: 静态节点, 如上图中的父节点 s （不包含 handler 的)
    // root: 如果插入的节点是第一个, 那么是root节点
    // catchAll: 有*匹配的节点
    // param: 参数节点，比如上图中的 :post 节点
    nType     nodeType

    // path 中的参数最大数量，最大只能保存 255 个（超过这个的情况貌似太难见到了）
    // 这里是一个非负的 8 进制数字，最大也只能是 255 了
    maxParams uint8

    // 和下面的 children 对应，保留的子节点的第一个字符
    // 如上图中的 s 节点，这里保存的就是 eu （earch 和 upport）的首字母
    indices   string

    // 当前节点的所有直接子节点
    children  []*node

    // 当前节点对应的 handler
    handle    Handle

    // 优先级，查找的时候会用到,表示当前节点加上所有子节点的数目
    priority  uint32
}
```

现在已经知道了 node 是用来做什么的了，就是用来保存树结构的，初始化的时候是一个空的树，没有任何数据， 接着看看怎么向这颗树中添加数据。

直接读注册路由的时候，Handle 方法中的 root.addRoute(path, handle) 没有去分析，现在就开始啃这个了， 追踪过去看看：

代码很长，真不喜欢读这种结构的代码，看起来比较累。

```go
// addRoute 将传入的 handle 添加到路径中
// 需要注意，这个操作不是并发安全的！！！！
func (n *node) addRoute(path string, handle Handle) {
    fullPath := path
    // 请求到这个方法，就给当前节点的权重 + 1
    n.priority++

    // 计算传入的路径参数的数量
    numParams := countParams(path)

    // 如果树不是空的
    // 判断的条件是当前节点的 path 的字符串长度和子节点的数量全部都大于 0
    // 就是说如果当前节点是一个空的节点，或者当前的节点是一个叶子节点，就直接
    // 进入 else 在当前节点下面添加子节点
    if len(n.path) > 0 || len(n.children) > 0 {
    // 定义一个 lable ，循环里面可以直接 break 到这里，适合这种嵌套的比较深的
    walk:
        for {
            // 如果传入的节点的最大参数的数量大于当前节点记录的数量，替换
            if numParams > n.maxParams {
                n.maxParams = numParams
            }

            // 查找最长的共同的前缀
            // 公共前缀不包含 ":" 或 "*"
            i := 0

            // 将最大值设置成长度较小的路径的长度
            max := min(len(path), len(n.path))

            // 循环，计算出当前节点和添加的节点共同前缀的长度
            for i < max && path[i] == n.path[i] {
                i++
            }

            // 如果相同前缀的长度比当前节点保存的 path 短
            // 比如当前节点现在的 path 是 sup ，添加的节点的 path 是 search
            // 它们相同的前缀就变成了 s ， s 比 sup 要短，符合 if 的条件，要做处理
            if i < len(n.path) {
                // 将当前节点的属性定义到一个子节点中，没有注释的属性不变，保持原样
                child := node{
                    // path 是当前节点的 path 去除公共前缀长度的部分
                    path:      n.path[i:],
                    wildChild: n.wildChild,
                    // 将类型更改为 static ，默认的，没有 handler 的节点
                    nType:     static,
                    indices:   n.indices,
                    children:  n.children,
                    handle:    n.handle,
                    // 权重 -1
                    priority:  n.priority - 1,
                }

                // 遍历当前节点的所有子节点（当前节点变成子节点之后的节点），
                // 如果最大参数数量大于当前节点的数量，更新
                for i := range child.children {
                    if child.children[i].maxParams > child.maxParams {
                        child.maxParams = child.children[i].maxParams
                    }
                }

                // 在当前节点的子节点定义为当前节点转换后的子节点
                n.children = []*node{&child}

                // 获取子节点的首字母,因为上面分割的时候是从 i 的位置开始分割
                // 所以 n.path[i] 可以去除子节点的首字母，理论上去 child.path[0] 也是可以的
                // 这里的 n.path[i] 取出来的是一个 uint8 类型的数字（代表字符），
                // 先用 []byte 包装一下数字再转换成字符串格式
                n.indices = string([]byte{n.path[i]})

                // 更新当前节点的 path 为新的公共前缀
                n.path = path[:i]

                // 将 handle 设置为 nil
                n.handle = nil

                // 肯定没有参数了，已经变成了一个没有 handle 的节点了
                n.wildChild = false
            }

            // 将新的节点添加到此节点的子节点， 这里是新添加节点的子节点
            if i < len(path) {
                // 截取掉公共部分，剩余的是子节点
                path = path[i:]

                // 如果当前路径有参数
                // 如果进入了上面 if i < len(n.path) 这个条件，这里就不会成立了
                // 因为上一个 if 中将 n.wildChild 重新定义成了 false
                // 什么情况会进入到这里呢 ?
                // 1. 上面的 if 不生效，也就是说不会有新的公共前缀， n.path = i 的时候
                // 2. 当前节点的 path 是一个参数节点就是像这种的 :post
                // 就是定义路由时候是这种形式的： blog/:post/update
                //
                if n.wildChild {
                    // 如果进入到了这里，证明这是一个参数节点，类似 :post 这种
                    // 不会这个节点进行处理，直接将它的子节点赋值给当前节点
                    // 比如： :post/ ，只要是参数节点，必有子节点，哪怕是
                    // blog/:post 这种，也有一个 / 的子节点
                    n = n.children[0]

                    // 又插入了一个节点，权限再次 + 1
                    n.priority++

                    // 更新当前节点的最大参数个数
                    if numParams > n.maxParams {
                        n.maxParams = numParams
                    }

                    // 更改方法中记录的参数个数，个人猜测，后面的逻辑还会用到
                    numParams--

                    // 检查通配符是否匹配
                    // 这里的 path 已经变成了去除了公共前缀的后面部分，比如
                    // :post/update ， 就是 /update
                    // 这里的 n 也已经是 :post 这种的下一级的节点，比如 / 或者 /u 等等
                    // 如果添加的节点的 path >= 当前节点的 path &&
                    // 当前节点的 path 长度和添加节点的前面相同数量的字符是相等的， &&
                    if len(path) >= len(n.path) && n.path == path[:len(n.path)] &&
                        // 简单更长的通配符，
                        // 当前节点的 path >= 添加节点的 path ，其实有第一个条件限制，
                        // 这里也只有 len(n.path) == len(path) 才会成立，
                        // 就是当前节点的 path 和 添加节点的 path 相等 ||
                        // 添加节点的 path 减去当前节点的 path 之后是 /
                        // 例如： n.path = name, path = name 或
                        // n.path = name, path = name/ 这两种情况
                        (len(n.path) >= len(path) || path[len(n.path)] == '/') {

                        // 跳出当前循环，进入下一次循环
                        // 再次循环的时候
                        // 1. if i < len(n.path) 这里就不会再进入了，现在 i == len(n.path)
                        // 2. if n.wildChild 也不会进入了，当前节点已经在上次循环的时候改为 children[0]
                        continue walk
                    } else {
                        // 当不是 n.path = name, path = name/ 这两种情况的时候，
                        // 代表通配符冲突了，什么意思呢？
                        // 简单的说就是通配符部分只允许定义相同的或者 / 结尾的
                        // 例如：blog/:post/update，再定义一个路由 blog/:postabc/add，
                        // 这个时候就会冲突了，是不被允许的，blog 后面只可以定义
                        // :post 或 :post/ 这种，同一个位置不允许使用多种通配符
                        // 这里的处理是直接 panic 了，如果想要支持，可以尝试重写下面部分代码


                        // 下面做的事情就是组合 panic 用到的提示信息
                        var pathSeg string

                        // 如果当前节点的类型是有*匹配的节点
                        if n.nType == catchAll {
                            pathSeg = path
                        } else {
                            // 如果不是，将 path 做字符串分割
                            // 这个是通过 / 分割，最多分成两个部分,然后取第一部分的值
                            // 例如： path = "name/hello/world"
                            // 分割两部分就是 name 和 hello/world , pathSeg = name
                            pathSeg = strings.SplitN(path, "/", 2)[0]
                        }

                        // 通过传入的原始路径来处理前缀, 可以到上面看下，方法进入就定义了这个变量
                        // 在原始路径中提取出 pathSeg 前面的部分在拼接上 n.path
                        // 例如： n.path = ":post" , fullPath="/blog/:postnew/add"
                        // 这时的 prefix = "/blog/:post"
                        prefix := fullPath[:strings.Index(fullPath, pathSeg)] + n.path

                        // 最终的提示信息就会生成类似这种：
                        // panic: ':postnew' in new path '/blog/:postnew/update/' \
                        // conflicts with existing wildcard ':post' in existing \
                        // prefix '/blog/:post'
                        // 就是说已经定义了 /blog/:post 这种规则的路由，
                        // 再定义 /blog/:postnew 这种就不被允许了
                        panic("'" + pathSeg +
                            "' in new path '" + fullPath +
                            "' conflicts with existing wildcard '" + n.path +
                            "' in existing prefix '" + prefix +
                            "'")
                    }
                }


                // 如果没有进入到上面的参数节点，当前节点不是一个参数节点 :post 这种
                c := path[0]

                // slash after param
                // 如果当前节点是一个参数节点，如 /:post && 是 / 开头的 && 只有一个子节点
                if n.nType == param && c == '/' && len(n.children) == 1 {
                    // /:post 这种节点不做处理，直接拿这个节点的子节点去匹配
                    n = n.children[0]
                    // 权重 + 1 ， 因为新的节点会变成这个节点的子节点
                    n.priority++

                    // 结束当前循环，再次进行匹配
                    // 比如： /:post/add 当前节点是 /:post 有个子节点 /add
                    // 新添加的节点是 /:post/update ，再次循环的时候将会用
                    // /add 和 /update 进行匹配， 就相当于是两次新的匹配，
                    // 由于 node 本身是一个指针，所以它们一直都会在 /:post 下面搞事情
                    continue walk
                }

                // 检查添加的 path 的首字母是否保存在在当前节点的 indices 中
                for i := 0; i < len(n.indices); i++ {
                    // 如果存在
                    if c == n.indices[i] {
                        // 这里处理优先级和排序的问题，把这个方法看完再去查看这个方法干了什么
                        i = n.incrementChildPrio(i)

                        // 将当前的节点替换成它对应的子节点
                        n = n.children[i]

                        // 结束当前循环，继续下一次
                        // 此时的当前 node n，已经变成了它的子字节，下次循环就从这个子节点开始处理了
                        // 比如：
                        // 当前节点是 s , 包含两个节点 earch 和 upport
                        // 这时 indices 就会有字母 e 和 u 并且是和子节点 earch 和 uppor 相对应
                        // 新添加的节点如果叫 subject ， 与当前节点匹配去除公共前缀 s 后， 就变成了
                        // ubject ，这时 n = upport 这个节点了，path = ubject
                        // 下一次循环就拿 upport 这个 n 和 ubject 这个 path 去开始下次的匹配
                        continue walk
                    }
                }

                // 如果上面 for 中也没有匹配上，就将新添加的节点插入

                // 如果第一个字符不是通配符 : 并且不是 *
                if c != ':' && c != '*' {
                    // []byte for proper unicode char conversion, see #65
                    // 将 path 的首字母 c 拼接到 indices 中
                    n.indices += string([]byte{c})

                    // 初始化一个新的节点
                    child := &node{
                        maxParams: numParams,
                    }

                    // 将这个新初始化的节点添加到当前节点的子节点中
                    n.children = append(n.children, child)

                    // 这个应该还是在处理优先级，一会再去看这个方法
                    n.incrementChildPrio(len(n.indices) - 1)

                    // 将当前节点替换成新生成的节点
                    n = child
                }

                // 用当前节点发起插入子节点的动作
                // 注意这个 n 已经替换成了上面新初始化的 child 了，只初始化了 maxParams
                // 属性，相当于是一个空的节点。insertChild 这个方法一会再去具体查看，
                // 现在只要知道它在做具体的插入子节点的动作就行，先把这个方法追踪完。
                n.insertChild(numParams, path, fullPath, handle)
                return

            // 如果公共前缀和当前添加的路径长度相等
            } else if i == len(path) { // Make node a (in-path) leaf
                // 如果当前节点不是 static 类型的， 就是已经注册了 handle 的节点
                // 就证明已经注册过这个路由，直接 panic 了
                if n.handle != nil {
                    panic("a handle is already registered for path '" + fullPath + "'")
                }

                // 如果是 nil 的，就是没有注册，证明这个节点是之前添加别的节点时候拆分出来的共同前缀
                // 例如： 添加过两个截点 /submit 和 /subject ，处理后就会有一个 static 的前缀 sub
                // 这是再添加一个 /sub 的路由，就是这里的情况了。

                n.handle = handle
            }

            // 这个新的节点被添加了， 出现了 return ， 只有出现这个才会正常退出循环，一次添加完成。
            return
        }
    } else { // Empty tree
        // 如果 n 是一个空格的节点，就直接调用插入子节点方法
        n.insertChild(numParams, path, fullPath, handle)

        // 并且它只有第一次插入的时候才会是空的，所以将 nType 定义成 root
        n.nType = root
    }
}
```

### incrementChildPrio

上面分析 addRoute 的时候调用了两次这个方法，大概知道它是处理优先级的，下面就追踪到这个方法，看看它到底干了什么：

```go
// increments priority of the given child and reorders if necessary
// 通过之前两次的调用，我们知道，这个 pos 都是 n.indices 中指定字符的索引，也就是位置
func (n *node) incrementChildPrio(pos int) int {
    // 因为 children 和 indices 是同时添加的，所以索引是相同的
    // 可以通过 pos 代表的位置找到， 将对应的子节点的优先级 + 1
    n.children[pos].priority++
    prio := n.children[pos].priority

    // 调整位置（向前移动）
    newPos := pos

    // 这里的操作就是将 pos 位置的子节点移动到最前面，增加优先级，比如：
    // 原本的节点是 [a, b, c, d], 传入的是 2 ,就变成 [c, a, b, d]
    // 这只是个示例，实际情况还要考虑 n.children[newPos-1].priority < prio ，不一定是移动全部
    for newPos > 0 && n.children[newPos-1].priority < prio {
        // swap node positions
        n.children[newPos-1], n.children[newPos] = n.children[newPos], n.children[newPos-1]

        newPos--
    }

    // build new index char string
    // 重新组织 n.indices 这个索引字符串，和上面做的是相同的事，只不过是用切割 slice 方式进行的
    // 从代码本身来说上面的部分也可以通过切割 slice 完成，不过移动应该是比切割组合 slice 效率高吧
    // 感觉是这样，没测试过，有兴趣的可以测试下。
    if newPos != pos {
        n.indices = n.indices[:newPos] + // unchanged prefix, might be empty
            n.indices[pos:pos+1] + // the index char we move
            n.indices[newPos:pos] + n.indices[pos+1:] // rest without char at 'pos'
    }

    // 最后返回的是调整后的位置，保证外面调用这个方法之后通过索引找到的对应的子节点是正确的
    return newPos
}
```

代码看完，大概总结一下，这个方法主要是在处理节点之前的关系，添加或者修改已经存在的，拼接出树的结构，真正写入插入节点的数据是在这个方法处理完关系后，调用 insertChild 方法来完成的。

### insertChild

现在就再追踪下，看看它是怎么插入数据的：

> 打开之后看一下，又是很长…… 有没有大脑翁的一下，仔细一行一行看吧。

```go
// 它传入的几个参数，我们可以回到 addRoute 看看都传给它什么了
// numParams 参数个数
// path 插入的子节点的路径
// fullPath 完整路径，就是注册路由时候的路径，没有被处理过的
// 注册路由对应的 handle 函数
func (n *node) insertChild(numParams uint8, path, fullPath string, handle Handle) {
    var offset int // 已经处理过的路径的所有字节数

    // find prefix until first wildcard (beginning with ':'' or '*'')
    // 查找前缀，知道第一个通配符（ 以 ':' 或 '*' 开头
    // 就是要将 path 遍历，提取出参数
    // 只要不是通配符开头的就不做处理，证明这个路由是没有参数的路由
    for i, max := 0, len(path); numParams > 0; i++ {
        c := path[i]

        // 如果不是 : 或 * 跳过本次循环，不做任何处理
        if c != ':' && c != '*' {
            continue
        }

        // 查询通配符后面的字符，直到查到 '/' 或者结束
        end := i + 1

        for end < max && path[end] != '/' {
            switch path[end] {
            // 通配符后面的名称不能包含 : 或 * ， 如 ::name 或 :*name 不允许定义
            case ':', '*':
                panic("only one wildcard per path segment is allowed, has: '" +
                    path[i:] + "' in path '" + fullPath + "'")
            default:
                end++
            }
        }

        // 检查通配符所在的位置，是否已经有子节点，如果有，就不能再插入
        // 例如： 已经定义了 /hello/name ， 就不能再定义 /hello/:param
        if len(n.children) > 0 {
            panic("wildcard route '" + path[i:end] +
                "' conflicts with existing children in path '" + fullPath + "'")
        }

        // 检查通配符是否有一个名字
        // 上面定义 end = i+1 ， 后面的 for 又执行了 ++ 操作，所以通配符 : 或 * 后面最少
        // 要有一个字符, 如： :a 或 :name ， :a 的时候 end 就是 i+2
        // 所以如果 end - i < 2 ，就是通配符后面没有对应的名称， 就会 panic
        if end-i < 2 {
            panic("wildcards must be named with a non-empty name in path '" + fullPath + "'")
        }

        // 如果 c 是 : 通配符的时候
        if c == ':' { // param
            // split path at the beginning of the wildcard
            // 从 offset 的位置，到查询到通配符的位置分割 path
            // 并把分割出来的路径定义到节点的 path 属性
            if i > 0 {
                n.path = path[offset:i]
                // 开始的位置变成了通配符所在的位置
                offset = i
            }

            // 将参数部分定义成一个子节点
            child := &node{
                nType:     param,
                maxParams: numParams,
            }

            // 用新定义的子节点初始化一个 children 属性
            n.children = []*node{child}
            // 标记上当前这个节点是一个包含参数的节点的节点
            n.wildChild = true

            // 将新创建的节点定义为当前节点，这个要想一下，到这里这种操作已经有不少了
            // 因为一直都是指针操作，修改都是指针的引用，所以定义好的层级关系不会被改变
            n = child
            // 新的节点权重 +1
            n.priority++
            // 最大参数个数 - 1
            numParams--

            // if the path doesn't end with the wildcard, then there
            // will be another non-wildcard subpath starting with '/'
            // 这个 end 有可能是结束或者下一个 / 的位置
            // 如果小于路径的最大长度，代表还包含子路径（也就是说后面还有子节点）
            if end < max {
                // 将通配符提取出来，赋值给 n.path， 现在 n.path 是 :name 这种格式的字符串
                n.path = path[offset:end]
                // 将起始位置移动到 end 的位置
                offset = end

                // 定义一个子节点，无论后面还有没有子节点 :name 这种格式的路由后面至少还有一个 /
                // 因为参数类型的节点不会保存 handler
                child := &node{
                    maxParams: numParams,
                    priority:  1,
                }

                // 将初始化的子节点赋值给当前节点
                n.children = []*node{child}

                // 当前节点又变成了新的子节点（到此时的 n 的身份已经转变了几次了，看这段代表的时候
                // 脑中要有一颗树，实在想不出来的话可以按照开始的图的结构，将节点的变化记录到一张纸上，
                // 然后将每一次的转变标记出来，就完全能明白了）
                n = child

                // 如果进入倒了这个 if ，执行到这里，就还会进入到下一次循环，将可能生成出来的参数再次去匹配
            }

            // 如果走到了这里，并且是循环的最后一次，就是已经将当前节点 n 定义成了叶子节点
            // 就可以进入到最下面代码部分，进行插入了
            // （上面的注释说明的位置是 else 上面，不是为下面的 else 做的注释）


        // 进入到 else ，就是包含 * 号的路由了
        } else { // catchAll
            // 这里的意思是， * 匹配的路径只允许定义在路由的最后一部分
            // 比如 : /hello/*world 是允许的， /hello/*world/more 这种就会 painc
            // 这种路径就是会将 hello/ 后面的所有内容变成 world 的变量
            // 比如地址栏输入： /hello/one/two/more ，获取到的参数 world = one/twq/more
            // 不会再将后面的 / 作为路径处理了
            if end != max || numParams > 1 {
                panic("catch-all routes are only allowed at the end of the path in path '" + fullPath + "'")
            }

            // 这种情况是，新定义的 * 通配符路由和其他已经定义的路由冲突了
            // 例如已经定义了一个 /hello/bro ， 又定义了一个 /hello/*world ，此时就会 panic 了
            if len(n.path) > 0 && n.path[len(n.path)-1] == '/' {
                panic("catch-all conflicts with existing handle for the path segment root in path '" + fullPath + "'")
            }

            // currently fixed width 1 for '/'
            i--
            // 这里是查询通配符前面是否有 / 没有 / 是不行的，panic
            if path[i] != '/' {
                panic("no / before catch-all in path '" + fullPath + "'")
            }

            n.path = path[offset:i]

            // first node: catchAll node with empty path
            // 后面的套路基本和之前看到的类似，就是定义一个子节点，保存通配符前面的路径，
            // 有变化的就是将 nType 定义为 catchAll，就是说代表这是一个  * 号匹配的路由
            child := &node{
                wildChild: true,
                nType:     catchAll,
                maxParams: 1,
            }
            n.children = []*node{child}
            n.indices = string(path[i])
            n = child
            n.priority++

            // second node: node holding the variable
            // 将下面的节点再添加到上面，不过 * 号路由不会再有下一级的节点了，因为它会将后面的
            // 的所有内容当做变量，即使它是个 / 符号
            child = &node{
                path:      path[i:],
                nType:     catchAll,
                maxParams: 1,
                handle:    handle,
                priority:  1,
            }
            n.children = []*node{child}

            // 这里 return 了，看下上面，因为已经将 handle 保存，查到了叶子节点，所以就直接结束了当前方法
            return
        }
    }

    // insert remaining path part and handle to the leaf
    // 这里给所有的处理的完成后的节点（不包含 * 通配符方式）的 handle 和 path 赋值
    n.path = path[offset:]
    n.handle = handle
}
```

代码读完，简单分析下，就是将没有了上级管理的路由（关系被 addRoute 方法处理了）， 交给这个方法处理，处理的时候会根据实际情况搞出新的关系，并最终在叶子节点写入注册的数据（主要是 path 和 handler)

### getValue

前面已经分析完成怎么向树中插入数据，现在不知道是否还记得，在 ServeHTTP 方法中还用到了个从树中取数据的方法 getValue， 其实插入如果搞明白了，这个取的方法理解起来就很容易了，按照插入时定义的规则，取出来即可，直接看代码：

```go
// 通过传入的参数的路径，寻找指定的路由
// 参数： 字符串类型的路径
// 返回值：
// Handle 通过路径找到的 handler ， 如果没有找到，会返回一个 nil
// Params 如果路由注册的是参数路由，这里会将参数及值返回，是一个 Param 结构的 slice
// tsr 是一个 boolean 类型的标志，告诉调用者这个路由是否可以被重定向，配合 Router
// 里面的 RedirectTrailingSlash 属性
func (n *node) getValue(path string) (handle Handle, p Params, tsr bool) {
walk: // outer loop for walking the tree
    for {
        // 如果要找的路径长度大于节点的长度
        // 入了根目录 / 一般情况都会满足这个条件，进入到这里走一圈
        if len(path) > len(n.path) {
            // 如果寻找的路径和节点保存的路径有相同的前缀
            if path[:len(n.path)] == n.path {
                // 将寻找路径的前缀部分去除
                path = path[len(n.path):]
                // 如果当前的节点不包含通配符（子节点有 ： 或者 * 这种的），就可以直接去子节点继续搜索。
                if !n.wildChild {
                    // 通过查找路径的首字符，快速找到包含这个字符的子节点（性能提升的地方就体现出来了）
                    c := path[0]
                    for i := 0; i < len(n.indices); i++ {
                        // 如果找到了，就继续去 indices 中对应的子节点去找
                        // 并且进入下一次循环
                        if c == n.indices[i] {
                            n = n.children[i]
                            continue walk
                        }
                    }

                    // 如果进入到了这里，代表没有找到对应的子节点

                    // 如果寻找路径正好是 / ，因为去除了公共路径
                    // 假设寻找的是 /hello/ ， n.path = hello， 这时 path 就是 /
                    // 当前的节点如果注册了 handle ，就证明这个是一个路由，tsr 就会变成 true
                    // 调用者就知道可以将路由重定向到 /hello 了，下次再来就可以找到这个路由了
                    // 可以去 Router.ServeHTTP 方法看下，是不是会发起一个重定向
                    tsr = (path == "/" && n.handle != nil)

                    // 返回，此时会返回 (nil,[],true|false) 这 3 个值
                    return
                }

                // 处理是通配符节点（包含参数或*）的情况
                // 如果当前节点是通配符节点，子节点肯定只有一个，因为 :name 这种格式的只允许注册一个
                // 所以可以把当前节点替换成 :name 或 *name 这种子节点
                n = n.children[0]
                switch n.nType {
                // 当子节点是参数类型的时候，:name 这种格式的
                case param:
                    // find param end (either '/' or path end)
                    // 查找出寻找路径中到结束或者 / 以前的部分所在的位置
                    // 比如是 :name 或者 :name/more ，end 就会等于 4
                    end := 0
                    for end < len(path) && path[end] != '/' {
                        end++
                    }

                    // save param value
                    // 保存参数的值， 因为参数是一个 slice ，使用前要初始化一下
                    if p == nil {
                        // lazy allocation
                        p = make(Params, 0, n.maxParams)
                        // 现在 p 的值是 []Param
                    }
                    // 下面两行是将 slice 的长度扩展
                    // 如果第一次来， p 是一个空格， i = 0
                    // p = p[:0+1] p 取的是一个位置 0 到 1 的切片，因为切片前开后闭，
                    // 所以它的长度就变成了 1， 下面就可以向切片添加一个值了
                    // 第二次来就是在原有的长度再扩展 1 的长度，然后再添加一个值，
                    // 这样这个参数就会永远是一个有效的最小长度，应该可以提高性能 ?
                    // 不清楚作者为什么这样做，直接 append 不可以吗？ 有兴趣的可以去确定下。
                    i := len(p)
                    p = p[:i+1] // expand slice within preallocated capacity

                    // 这个就是给 Parm 赋值， Key 就是去掉 : 后的值，比如 :name 就是 name
                    // Value 就是路径到结束位置 end 的值
                    p[i].Key = n.path[1:]
                    p[i].Value = path[:end]

                    // we need to go deeper!
                    // 上面已经获取了第一个值，但是有可能还会有两个或者更多，所以还要继续处理
                    // 如果之前找到第一个值的位置比寻找路径的长度小，就是还有 :name 后面还有内容需要去匹配
                    if end < len(path) {
                        // 如果还存在子节点， 就将 path 和 node 替换成子节点的，跳出当前循环再找一遍
                        if len(n.children) > 0 {
                            path = path[end:]
                            n = n.children[0]
                            continue walk
                        }

                        // ... but we can't
                        tsr = (len(path) == end+1)
                        // 这里返回的是 (nil, [Parm{name, 路径上的值}], true|false)
                        return
                    }

                    // 如果当前节点保存了 handle ，就将 handle 返回给调用者，去发起业务逻辑
                    if handle = n.handle; handle != nil {
                        return
                    } else if len(n.children) == 1 {
                        // No handle found. Check if a handle for this path + a
                        // trailing slash exists for TSR recommendation
                        // 确认它是否有一个 / 的子节点，如果有 tsr = true
                        n = n.children[0]
                        tsr = (n.path == "/" && n.handle != nil)
                    }

                    // 这里返回的是 [nil, [找到的参数...], true|false]
                    return

                // 如果是 * 这种方式的路由, 就是直接处理参数，并返回，因为它肯定是最后一个节点
                // 所以也不会出现子节点，及时路径上后面根再多的 / 也都只是参数的值
                case catchAll:
                    // save param value
                    if p == nil {
                        // lazy allocation
                        p = make(Params, 0, n.maxParams)
                    }
                    i := len(p)
                    p = p[:i+1] // expand slice within preallocated capacity
                    p[i].Key = n.path[2:]
                    p[i].Value = path

                    handle = n.handle
                    return

                default:
                    // 没想到什么情况会出现无效的类型， 如果出现了就直接 panic
                    // 现在这个 panic 不会导致服务器崩溃，因为在 ServeHTTP 中做了 Recovery
                    panic("invalid node type")
                }
            }
        // 如果查询的路径和节点保存的路径相同
        } else if path == n.path {
            // We should have reached the node containing the handle.
            // Check if this node has a handle registered.
            // 检查下这个节点是否包含 handle ，如果包含就返回
            // 会有不包含的情况，比如恰巧这个就是一个相同前缀的节点
            if handle = n.handle; handle != nil {
                return
            }

            // 如果查找的路径是 / 又不是根节点，还是个参数节点，就允许重定向
            if path == "/" && n.wildChild && n.nType != root {
                tsr = true
                return
            }

            // No handle found. Check if a handle for this path + a
            // trailing slash exists for trailing slash recommendation
            // 如果没有找到，但是有一个 / 的子节点，也允许它重定向
            // 就是这个意思，比如定义了一个 /hello/ 路由，但是访问的是 /hello
            // 也允许它重定向到 /hello/ 上
            for i := 0; i < len(n.indices); i++ {
                if n.indices[i] == '/' {
                    n = n.children[i]
                    tsr = (len(n.path) == 1 && n.handle != nil) ||
                        (n.nType == catchAll && n.children[0].handle != nil)
                    return
                }
            }

            return
        }

        // Nothing found. We can recommend to redirect to the same URL with an
        // extra trailing slash if a leaf exists for that path
        // 是否允许重定向的一个验证
        tsr = (path == "/") ||
            (len(n.path) == len(path)+1 && n.path[len(path)] == '/' &&
                path == n.path[:len(n.path)-1] && n.handle != nil)
        return
    }
}
```

到这里源码就分析完了，整理的代码量不算多， `router.go` 文件还好理解，基本看一遍就清除，这个 `tree.go` ，看了几个小时才基本上看明白，逻辑有点复杂，而且也不是很习惯这个写法。 搞明白了之后再用的时候遇到问题也是很清楚哪里出现的，并且定义路由的时候也知道要怎么去定义了，那种可以支持，那种不可以支持。
