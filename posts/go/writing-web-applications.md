---
title: "Go 编写 Web 应用"
author: "BroQiang"
github_url: "https://broqiang.com"
created_at: 2019-04-02T22:54:51
updated_at: 2019-04-02T22:54:51
description: "这是一篇官方的 https://golang.org/doc/articles/wiki 的翻译， 描述了怎样使用 go 的 http 包快速创建一个 web 应用"
tags: ["go", "翻译"]
---

这是一篇官方的 [Writing Web Applications](https://golang.org/doc/articles/wiki) 的翻译

如果 `golang.org` 打不开的话，可以把所有链接中的 `golang.org` 更换成 `golang.google.cn` ，这是一个官方的国内镜像，和 `golang.org` 的内容是一致的。

## 简介

### 本教程包含下面内容

- 通过 load 与 save 方法创建数据结构

- 使用 `net/http` 去构建 web 应用

- 使用 `html/template` 去处理 HTML 模板

- 使用 `regexp` 包去验证用户的输入

- 使用闭包

### 需要已经掌握下面知识

- 编程经验

- 了解基本的 web 技术(HTTP,HTML)

- 会一些简单 UNIX/DOS 命令行

## 入门

现在，你可以在 FreeBSD, Linux, OS X, 或 Windows 上运行 Go ， 我们使用 `$` 代表命令行的提示符， `$` 开始的内容是输入的命令。

安装 Go (查看 [安装文档](https://golang.org/doc/install) )

在 GOPATH 中为这个项目创建一个新的目录，然后切换到这个目录(cd)

```bash
cd $GOPATH/src
mkdir gowiki
cd go wiki
```

创建一个文件 `wiki.go` ， 用你喜欢的编辑器打开它，用添加下面内容

```go
package main

import (
    "fmt"
    "io/ioutil"
)
```

我们从 Go 的标准库中导入了 `fmt` 和 `io/ioutil` 包, 稍后，当我们实现其他功能时，将会在 `import` 声明中添加更多的包。

## 数据结构

我们从定义数据结构开始，一个 `wiki` 由许多相关的页面组成，每一个页面都包含一个 title 和 body （页面的内容）。我们定义一个 Page 结构，包含两个代表 title 和 body 的字段。

```go
type Page struct {
    Title string
    Body []byte
}
```

类型 `[]byte` 就是一个 byte 切片（关于切片的更多信息，查看： [Slices: usage and internals](https://blog.golang.org/go-slices-usage-and-internals)）。 Body 是元素是一个 `[]byte` 而不是 `string`， 是因为将要使用的 `io` 库需要这个类型，可以在下面看到。

Page 结构中体现了如何将页面数据保存到内存中，但是如果是持久存储怎么办？我们可以在 Page 上创建一个 save 方法来解决：

```go
func (p *Page) save() error {
    filename := p.Title + ".txt"

    return ioutil.WriteFile(filename, p.Body, 0600)
}
```

这个方法的签名解读："这是一个名字叫 save 的方法，它的接收者 p 指向一个 Page 的指针。它不需要参数，并返回一个类型为 error 的值"

这个方法将保存 Page 的 Body 到一个文本文件。为了简单，我们使用 Title 来做它的文件名。

save 方法返回了一个 error 类型的值，是因为它是 WriteFile 的返回类型（一个用于将字节切片写入到文件的标准库函数）。 save 方法返回的是一个 error 类型的值，所以在写入文件出现错误的时候应该去处理它。如果没有出现错误，应该返回一个 nil （指针，接口和一些其他类型的零值）

八进制整数 `0600`， 是传给 WriteFile 的第三个参数，表示只有当前用于拥有文件的读写权限。（通过 Unix man page open（2）查看详细说明 "译者注： Linux 上可以通过这个命令查看： man 2 open "）

除了保存页面，我们也需要加载页面：

```go
func loadPage(title string) *Page {
    filename := title + ".txt"
    body, _ := ioutil.ReadFile(filename)

    return &Page{Title: title, Body: body}
}
```

loadPage 通过参数 title 拼接了文件名，将文件内容读取到一个新的变量 body 中，并且返回了一个通过 title 和 body 构造的指向 Page 字面的指针。

函数可以返回多个值。标准库函数 io.ReadFile 返回了一个 []byte 和一个 error 。在 loadPage 还没有处理错误，通过使用空白符 `_` 将返回的错误丢弃（就是将值赋给一个空）。

但是，如果 ReadFile 遇到一个错误会发生什么？ 例如，文件不存在。 我们不应该忽略这个错误, 让我们修改函数，来返回 \*Page 和 error 。

```go
func loadPage(title string) (*Page, error) {
    filename := title + ".txt"
    body, err := ioutil.ReadFile(filename)

    if err != nil {
        return nil, err
    }

    return &Page{Title: title, Body: body}, nil
}
```

此函数的调用者现在可以通过检查第二个参数； 如果它是 nil ，表示成功加载了一个 Page 。如果不是，它将是一个 error ，并且可以由调用者处理（详细查看 [语言规范](https://golang.org/ref/spec#Errors) ）。

现在，我们有了一个简单的数据结构，可以保存并加载一个文件。让我们写一个 main 函数来测试我们写的东西。

```go
func main() {
    p1 := &Page{Title: "TestPage", Body: []byte("This is a simple Page.")}
    p1.save()
    p2, _ := loadPage("TestPage")
    fmt.Println(string(p2.Body))
}
```

编译并执行这个代码，一个包含了 p1 的内容的名字是 TestPage.txt 的文件将被创建。这个文件将被读进结构 p2 中，并将它的 Body 元素打印到屏幕上。

你可以像这样编译并运行这个程序：

```bash
$ go build wiki.go
$ ./wiki
This is a simple Page.
```

（如果你使用的是 Windows ，你必须输入 `wiki` ，去掉 `./` 去执行这个程序）

[点击这里去查看我们到现在写的代码](https://golang.org/doc/articles/wiki/part1.go)

## net http 包简介

这是一个简单的 web 服务器完整的工作示例：

```go
package main

import (
    "fmt"
    "log"
    "net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
}

func main() {
    http.HandleFunc("/", handler)
    log.Fatal(http.ListenAndServe(":8080", nil))
}
```

main 函数从调用 http.HandleFunc 开始，它告诉 http 包使用 handler 去处理所有访问 web 根目录（"/"）的请求。

然后它调用 http.ListenAndServe ，指定它在任意接口上监听 8080 端口（"8080"）。 （现在不用去管它的第二个参数， nil ） 这个函数将被阻塞，知道程序终止。

ListenAndServe 始终返回一个 error ，并且它只有发生意外错误时才会返回。为了记录这个错误，我们在它的外面包裹一个叫 log.Fatal 的函数。

handler 是一个 http.HandlerFunc 类型的函数，它需要两个参数，http.ResponseWriter 和 http.Request 。

http.ResponseWriter 的值集合了 HTTP 服务器的响应。通过写入它，我们将数据发送到 HTTP 客户端。

http.Request 是表示客户端 HTTP 请求的数据结构。 r.URL.Path 是一个请求地址组成的路径。它的后面跟随 `[1:]` 的意思是： "创建一个从第一个字符到结尾的子切片" 。从路径中删除开头的 `/` 。

如果你运行程序，并访问这个地址：

`http://localhost:8080/monkeys`

程序将会显示包含下面内容的页面：

`Hi there, I love monkeys!`

## 使用 net http 去服务 wiki 页面

要想使用 net/http 包，它必须被导入：

```go
import (
    "fmt"
    "io/ioutil"
    "net/http"
)
```

我们来创建一个叫 viewHandler 的 handler， 用户可以通过它去查看 wiki 页面。它将去处理包含前缀 `/view/` 的 URL 。

```go
func viewHandler(w http.ResponseWriter, r *http.Request) {
    title := r.URL.Path[len("/view/"):]
    p, _ := loadPage(title)
    fmt.Fprintf(w, "<h1>%s</h1><div>%s</div>", p.Title, p.Body)
}
```

再次注意，使用 `_` 来忽略 loadPage 返回的错误。这里是为了简单，但这是一个坏的习惯，稍后我们会处理它。

首先，这个函数从请求 URL 的 path 组件 `r.URL.Path` 中提取页面的标题。再通过切片 `[len("/view/"):]` 去掉前面路径前面的 `/view/` ，这是因为路径总是以 `/view/` 开始，它不是页面的一部分。

然后函数去加载页面数据，格式化成一个简单的 HTML 格式的字符串并写入到 http.ResponseWriter `w` 。

要使用这个 handler ，我们重写我们的 main 函数，去初始化 http 通过 viewHandler 去处理每一个 `/view/` 下面的请求。

```go
func main() {
    http.HandleFunc("/view/", viewHandler)
    log.Fatal(http.ListenAndServe(":8080", nil))
}
```

[点击这里去查看我们到现在写的代码](https://golang.org/doc/articles/wiki/part2.go)

我们创建一些测试页面数据（如： test.txt），编译我们的代码，并尝试服务器 wiki 页面。

使用编辑器打开 test.txt 文件，写入 `Hello world` 。

```bash
go build wiki.go
./wiki
```

（如果你使用的是 Windows ，你必须输入 `wiki` ，去掉 `./` 去执行这个程序）

服务启动后，访问 [http://localhost:8080/view/test](http://localhost:8080/view/test) 将会显示页面，标题是 test ，内容是 Hello world 。

## 编辑页面

wiki 不是一个不能编辑页面的 wiki 。我们来创建两个新的 handler， 一个叫 editHandler ，用来显示编辑页面的 form 表单，另一个叫 saveHandler ，用来保存 form 表单提交的数据。

首先，我们添加它们到 main() 中：

```go
func main() {
    http.HandleFunc("/view/", viewHandler)
    http.HandleFunc("/edit/", editHandler)
    http.HandleFunc("/save/", saveHandler)
    log.Fatal(http.ListenAndServe(":8080", nil))
}
```

editHandler 函数加载页面（或者，它不存在时创建一个空的 Page 结构）， 并且显示一个 HTML form 表单。

```go
func editHandler(w http.ResponseWriter, r *http.Request) {
    title := r.URL.Path[len("/edit/"):]
    p, err := loadPage(title)
    if err != nil {
        p = &Page{Title: title}
    }
    fmt.Fprintf(w, "<h1>Editing %s</h1>"+
        "<form action=\"/save/%s\" method=\"POST\">"+
        "<textarea name=\"body\">%s</textarea><br>"+
        "<input type=\"submit\" value=\"Save\">"+
        "</form>",
        p.Title, p.Title, p.Body)
}
```

这个函数可以很好的工作，但是所有硬编码的 HTML 都是非常丑的，所以还有更好的办法。

## html template 包

`html/template` 是 Go 标准库的一部分。我们使用 html/template ，可以将 HTML 保存到一个单独的文件中，允许我们在不改动底层 Go 代码的情况下改变我们的编辑页面的布局。

首先，我们必须将 html/template 添加到 import 的列表中。我们也不会再使用 fmt 包了，所以将它删除。

```go
import (
    "html/template"
    "io/ioutil"
    "net/http"
)
```

我们来创建一个包含 HTML 表单的模板文件。 打开一个文件，命名为 edit.html 并且添加下面内容：

```html
<h1>Editing {{.Title}}</h1>

<form action="/save/{{.Title}}" method="POST">
  <div>
    <textarea name="body" rows="20" cols="80">{{printf "%s" .Body}}</textarea>
  </div>
  <div><input type="submit" value="Save" /></div>
</form>
```

修改 editHandler 去使用模板，替换硬编码的 HTML：

```go
func editHandler(w http.ResponseWriter, r *http.Request) {
    title := r.URL.Path[len("/edit/"):]
    p, err := loadPage(title)
    if err != nil {
        p = &Page{Title: title}
    }
    t, _ := template.ParseFiles("edit.html")
    t.Execute(w, p)
}
```

`template.ParseFiles` 函数将读取 `edit.html` 的内容，并返回一个 `*template.Template`

执行模板的 `t.Execute` 方法，将生成的 HTML 写入到 http.ResponseWriter 。`.Title` 和 `.Body` 的点符号参考 `p.Title` 和 `p.Body` 。

模板指定用双大括号括起来，`printf "%s" .Body` 是一个函数调用，它将输出的 `.Body` 的字节流替换成字符串，与调用 `fmt.Printf` 相同。`html/template` 可以保证模板操作只有安全的并且正确的 HTML 被生成，它会自动转译大于号符号 `>`，替换成 `&gt;`，确保用户数据不会破坏表单数据。

我们现在已经使用模板了，就再创建一个用于 viewHandler 函数调用的 view.html 模板：

```html
<h1>{{.Title}}</h1>

<p>[<a href="/edit/{{.Title}}">edit</a>]</p>

<div>{{printf "%s" .Body}}</div>
```

修改对应的 viewHandler 函数：

```go
func viewHandler(w http.ResponseWriter, r *http.Request) {
    title := r.URL.Path[len("/view/"):]
    p, _ := loadPage(title)
    t, _ := template.ParseFiles("view.html")
    t.Execute(w, p)
}
```

注意，我们在两个 `handler` 中使用了几乎完全一样的模板代码。我们将模板代码放到一个单独的函数中，来去除这个重复：

```go
func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
    t, _ := template.ParseFiles(tmpl + ".html")
    t.Execute(w, p)
}
```

修改连个 handler 来使用这个函数：

```go
func viewHandler(w http.ResponseWriter, r *http.Request) {
    title := r.URL.Path[len("/view/"):]
    p, _ := loadPage(title)
    renderTemplate(w, "view", p)
}
```

```go
func editHandler(w http.ResponseWriter, r *http.Request) {
    title := r.URL.Path[len("/edit/"):]
    p, err := loadPage(title)
    if err != nil {
        p = &Page{Title: title}
    }
    renderTemplate(w, "edit", p)
}
```

我们可以先注释掉在 main 函数中注册的未实现的 `save` `handler`，然后就可以重新编辑并测试我们的程序了。

[点击这里去查看我们到现在写的代码](https://golang.org/doc/articles/wiki/part3.go)

## 处理不存在的页面

如果你输入 [/view/APageThatDoesntExist](http://localhost:8080//view/APageThatDoesntExist) ，你将看到一个包含 HTML 的页面，这是因为它忽略了 `loadPage` `error` 返回值，并继续去尝试填充了一个没有数据的模板去替换，如果请求的页面不存在，它应该重定向到编辑页面，这样就可以创建内容。

```go
func viewHandler(w http.ResponseWriter, r *http.Request) {
    title := r.URL.Path[len("/view/"):]
    p, err := loadPage(title)
    if err != nil {
        http.Redirect(w, r, "/edit/"+title, http.StatusFound)
        return
    }
    renderTemplate(w, "view", p)
}
```

`http.Redirect` 添加一个 http 状态码 `http.StatusFound` （302） 和 一个地址到 http 响应。

## 保存页面

`saveHandler` 函数将处理编辑页面 form 表单提交的数据，在 mian 函数中取消和此相关的注释，然后实现这个 handler 。

```go
func saveHandler(w http.ResponseWriter, r *http.Request) {
    title := r.URL.Path[len("/save/"):]
    body := r.FormValue("body")
    p := &Page{Title: title, Body: []byte(body)}
    p.save()
    http.Redirect(w, r, "/view/"+title, http.StatusFound)
}
```

这个页面的 `title` （由 URL 提供）和表单中的唯一字段 `Body` 被保存到一个新的 `Page` ，然后调用 `save()` 方法写入数据到一个文件，并且将客户端重定向到 `/view/` 页面。

`FormValue` 的返回值是一个字符串类型，我们在将它填充到 Page 前必选转换成 []byte，使用 `[]byte(body)` 可以执行转换。

## 错误处理

在我们的程序中，有几个地方忽略了错误。这种做法非常不好，尤其是当错误发生时，我们的程序会出现意外的行为。一个好的方案是去处理错误，并将错误消息返回给用户。这样，当出现错误的时候，服务器将按照我们想要的方式运行，并且通知给用户。

首先，我们在 `renderTemplate` 中处理错误：

```go
func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
    t, err := template.ParseFiles(tmpl + ".html")
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    err = t.Execute(w, p)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}
```

`http.Error` 函数发送一个指定的 HTTP 状态码（示例中是 `Internal Server Error`）和错误消息。这个抽离出一个单独的函数是一个多么明智的做法，否则要改多个地方了。

现在，我们来修改 `saveHandler` ：

```go
func saveHandler(w http.ResponseWriter, r *http.Request) {
    title := r.URL.Path[len("/save/"):]
    body := r.FormValue("body")
    p := &Page{Title: title, Body: []byte(body)}
    err := p.save()
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    http.Redirect(w, r, "/view/"+title, http.StatusFound)
}
```

在 `p.save()` 时发生任何错误，将会报告给用户。

## 模板缓存

这是一段低效的代码： 页面每一次展示的时候，`renderTemplate` 都会调用 `ParseFiles` 。好的方式是在程序初始化的时候只调用一个 `ParseFiles` ，解析所有的模板到一个单一的 `*Template` ，然后我们可以使用 `ExecuteTemplate` 方法去渲染指定的模板。

首先，创建一个全局变量 `templates` ，并且通过 `ParseFiles` 初始化它。

```go
var templates = template.Must(template.ParseFiles("edit.html", "view.html"))
```

`template.Must` 函数是一个方便的包装器，当传入一个不是 `nil` 的 `error` 时会产生一个 `panic`，除此之外，返回的 `*Template` 不会改变。一个 `panic` 在这里是合适，如果模板不能被加载只有退出程序才是明智的。

`ParseFiles` 函数接受任意数量的用来识别我们的模板文件的字符串参数，然后将这些文件解析到 `templates`，并且通过基础的文件名命名。如果我们要添加更多的模板到我们的程序，我们需要将它们的名字添加到 `ParseFiles` 的参数中。

然后我们修改 `renderTemplate` 函数，通过 `templates.ExecuteTemplate` 去调用与名称相对应的模板。

```go
func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
    err := templates.ExecuteTemplate(w, tmpl+".html", p)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}
```

注意，模板名称是模板的文件名，所以我们必须在 teml 参数后面添加 `.html` 。

## 验证

可能你已经注意到，这个程序有一个严重的安全漏洞：用户可以通过任意路径在服务器上读/写。为了解决这个，我们可以编写一个函数，通过正则表达式来验证 title 。

首先，在 `import` 列表添加 `regexp` ,然后我们可以创建一个全局变量去保存我们的正则表达式。

```go
var validPath = regexp.MustCompile("^/(edit|save|view)/([a-zA-Z0-9]+)$")
```

函数 `regexp.MustCompile` 将解析并编译正则表达式，返回一个 `regexp.Regexp` 。`MustCompile` 区别于 `Compile` 的是，当正则表达式编译失败的时候，它将产生一个 `panic`， 而 `Compile` 会通过第二个参数返回一个 `error` 。

现在，我们编写一个函数来使用 `validPath` 来验证路径并提取页面的 title 。

```go
func getTitle(w http.ResponseWriter, r *http.Request) (string, error) {
    m := validPath.FindStringSubmatch(r.URL.Path)
    if m == nil {
        http.NotFound(w, r)
        return "", errors.New("Invalid Page Title")
    }
    return m[2], nil // title 是第二个子表达式
}
```

如果 title 是有效的，它将和一个值为 `nil` 的 `error` 一起返回。如果 title 是无效的，函数将写入一个 `404 Not Found` 错误到 HTTP 链接，并且返回一个错误到 handler 。要创建一个新的 `error` ，我们还必须要导入 `errors` 包。

我们将 `getTitle` 调用放到每一个 handler 中：

```go
func viewHandler(w http.ResponseWriter, r *http.Request) {
    title, err := getTitle(w, r)
    if err != nil {
        return
    }
    p, err := loadPage(title)
    if err != nil {
        http.Redirect(w, r, "/edit/"+title, http.StatusFound)
        return
    }
    renderTemplate(w, "view", p)
}
```

```go
func editHandler(w http.ResponseWriter, r *http.Request) {
    title, err := getTitle(w, r)
    if err != nil {
        return
    }
    p, err := loadPage(title)
    if err != nil {
        p = &Page{Title: title}
    }
    renderTemplate(w, "edit", p)
}
```

```go
func saveHandler(w http.ResponseWriter, r *http.Request) {
    title, err := getTitle(w, r)
    if err != nil {
        return
    }
    body := r.FormValue("body")
    p := &Page{Title: title, Body: []byte(body)}
    err = p.save()
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    http.Redirect(w, r, "/view/"+title, http.StatusFound)
}
```

## 函数文字和闭包简介

捕获每一个函数中的错误处理条件会产生大量的冗余代码。如果我们将每一个 handler 包装进一个函数去处理验证和错误检查呢？ Go 的函数字面提供了一个强大的抽象功能可以帮助我们实现它。

首先，我们重写每一个 handler 函数的定义，接收一个字符串类型的 title 参数：

```go
func viewHandler(w http.ResponseWriter, r *http.Request, title string)
func editHandler(w http.ResponseWriter, r *http.Request, title string)
func saveHandler(w http.ResponseWriter, r *http.Request, title string)
```

现在我们定义一个参数是上面函数类型的包装函数，并且返回一个类型为 `http.HandlerFunc` 的函数（ 这个返回值是为了满足 http.HandleFunc 的参数 ）。

```go
func makeHandler(fn func (http.ResponseWriter, *http.Request, string)) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
        // 这里将从 Request 中提取 title
		// 并且调用提供的 handler 'fn'
	}
}
```

这个返回的函数叫闭包，因为它包含了在它外部定义的值。在这种情况，变量 fn (传给 makeHandler 的单个参数) 包裹在闭包中， 变量 fn 将会是我们的 `save`,`edit` 或 `view` handler 中的一个。

现在我们可以从 `getTitle` 中提取代码，并在此处使用它（稍作修改）：

```go
func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        m := validPath.FindStringSubmatch(r.URL.Path)
        if m == nil {
            http.NotFound(w, r)
            return
        }
        fn(w, r, m[2])
    }
}
```

`makeHandler` 返回的闭包是一个包含 `http.ResponseWriter` 和 `http.Request` 参数的函数（ 换句话说，一个 `http.HandlerFunc` ）。这个闭包从请求路径中提取 title ，并且通过 title 验证器 `regexp` 验证它。 如果标题是无效的，一个 `error` 将通过 `http.NotFound` 函数写入到 `ResponseWriter` 。如果 title 是有效的，它包裹的 handler 函数 fn 将传入 ResponseWriter, Request, 和 title 参数被调用。

现在，在 main 函数中，我们可以在 `handler` 被注册到 `http` 包之前，通过 `makeHandler` 来包装它们：

```go
func main() {
    http.HandleFunc("/view/", makeHandler(viewHandler))
    http.HandleFunc("/edit/", makeHandler(editHandler))
    http.HandleFunc("/save/", makeHandler(saveHandler))

    log.Fatal(http.ListenAndServe(":8080", nil))
}
```

最后，我们从 handler 中删除对 `getTitle` 的调用，使它们更简单：

```go
func viewHandler(w http.ResponseWriter, r *http.Request, title string) {
    p, err := loadPage(title)
    if err != nil {
        http.Redirect(w, r, "/edit/"+title, http.StatusFound)
        return
    }
    renderTemplate(w, "view", p)
}
```

```go
func editHandler(w http.ResponseWriter, r *http.Request, title string) {
    p, err := loadPage(title)
    if err != nil {
        p = &Page{Title: title}
    }
    renderTemplate(w, "edit", p)
}
```

```go
func saveHandler(w http.ResponseWriter, r *http.Request, title string) {
    body := r.FormValue("body")
    p := &Page{Title: title, Body: []byte(body)}
    err := p.save()
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    http.Redirect(w, r, "/view/"+title, http.StatusFound)
}
```

## 试试看

[点击这里可以查看最终的代码](https://golang.org/doc/articles/wiki/final.go)

重新编译代码并运行它：

```bash
go build wiki.go
./wiki
```

访问 [http://localhost:8080/view/ANewPage](http://localhost:8080/view/ANewPage) 将显示页面编辑表单, 你应该可以输入一些文本，点击 `save`，然后会重定向到新创建的页面。

## 其他任务

下面是希望你可能希望自己解决的任务清单：

- 将模板保存到 `tmpl/` 目录并且将页面数据保存到 `data/` 。

- 添加一个 handler 将 web 根目录重定向到 `/view/FrontPage` 。

- 完善页面模板，使它们成为一个有效的 HTML 并添加一些 CSS 规则。

- 通过 [PageName] 实例实现一些页面间的链接，`<a href="/view/PageName">PageName</a>` （提示： 你可以通过 `regexp.ReplaceAllFunc` 来实现这个）。
