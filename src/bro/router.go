package bro

import (
	"fmt"
	"net/http"
	"path"
	"strings"
)

var allowHTTPMethod = []string{"GET", "POST", "PUT", "PATH", "HEAD", "OPTIONS", "DELETE", "CONNECT", "TRACE"}

func checkAllowHTTPMethod(method string) {
	for _, m := range allowHTTPMethod {
		if method == m {
			return
		}
	}

	panic(fmt.Sprintf("http method %s is not valid. \n", method))
}

type (
	// Handler 自定义 Handler
	Handler func(*Context)

	// Handlers 是一个 Handler 的切片
	Handlers []Handler

	// RouterGroup 是一个路由分组
	RouterGroup struct {
		Handlers Handlers
		path     string
		engine   *Engine
		root     bool
	}
)

// GET 注册 GET 方式的路由
func (group *RouterGroup) GET(relativePath string, handlers ...Handler) {
	group.handle("GET", relativePath, handlers)
}

// POST 注册 POST 方式的路由
func (group *RouterGroup) POST(relativePath string, handlers ...Handler) {
	group.handle("POST", relativePath, handlers)
}

// PUT 注册 PUT 方式的路由
func (group *RouterGroup) PUT(relativePath string, handlers ...Handler) {
	group.handle("PUT", relativePath, handlers)
}

// PATCH 注册 PATCH 方式的路由
func (group *RouterGroup) PATCH(relativePath string, handlers ...Handler) {
	group.handle("PATCH", relativePath, handlers)
}

// HEAD 注册 HEAD 方式的路由
func (group *RouterGroup) HEAD(relativePath string, handlers ...Handler) {
	group.handle("HEAD", relativePath, handlers)
}

// OPTIONS 注册 OPTIONS 方式的路由
func (group *RouterGroup) OPTIONS(relativePath string, handlers ...Handler) {
	group.handle("OPTIONS", relativePath, handlers)
}

// DELETE 注册 DELETE 方式的路由
func (group *RouterGroup) DELETE(relativePath string, handlers ...Handler) {
	group.handle("DELETE", relativePath, handlers)
}

// CONNECT 注册 CONNECT 方式的路由
func (group *RouterGroup) CONNECT(relativePath string, handlers ...Handler) {
	group.handle("CONNECT", relativePath, handlers)
}

// TRACE 注册 TRACE 方式的路由
func (group *RouterGroup) TRACE(relativePath string, handlers ...Handler) {
	group.handle("TRACE", relativePath, handlers)
}

// Any 会同时注册 "GET", "POST", "PUT", "PATH", "HEAD", "OPTIONS", "DELETE", "CONNECT", "TRACE"
func (group *RouterGroup) Any(relativePath string, handlers ...Handler) {

}

// Handle 注册路由
func (group *RouterGroup) Handle(method, relativePath string, handlers ...Handler) {
	checkAllowHTTPMethod(method)

	group.handle(method, relativePath, handlers)
}

func (group *RouterGroup) handle(method, relativePath string, handlers Handlers) {
	path := JoinPath(group.path, relativePath)
	handlers = group.combineHandlers(handlers)

	root := group.engine.trees[method]
	if root == nil {
		root = new(node)
		group.engine.trees[method] = root
	}

	root.addRoute(path, handlers)
}

func (group *RouterGroup) combineHandlers(handlers Handlers) Handlers {
	gl := len(group.Handlers)

	mergedHandlers := make(Handlers, gl+len(handlers))

	copy(mergedHandlers, group.Handlers)
	copy(mergedHandlers[gl:], handlers)

	return mergedHandlers
}

// Last 获取 slice 中的最后一个 Handler
func (handlers Handlers) Last() Handler {
	if length := len(handlers); length > 0 {
		return handlers[length-1]
	}

	return nil
}

// 初始化一个 Context
func (group *RouterGroup) createContext(w http.ResponseWriter, r *http.Request, ps Params, handlers Handlers) *Context {
	return &Context{
		Request:  r,
		Writer:   w,
		Params:   ps,
		handlers: handlers,
		engine:   group.engine,
		index:    -1,
	}
}

// StaticFile 注册一个静态文件
func (group *RouterGroup) StaticFile(relativePath, filepath string) {
	if strings.Contains(relativePath, ":") || strings.Contains(relativePath, "*") {
		panic("URL parameters can not be used when serving a static file")
	}

	handler := func(c *Context) {
		http.ServeFile(c.Writer, c.Request, filepath)
	}

	group.GET(relativePath, handler)
	group.HEAD(relativePath, handler)
}

// StaticDir 允许传入字符串类型的目录
func (group *RouterGroup) StaticDir(relativePath, root string) {
	group.StaticFS(relativePath, http.Dir(root))
}

// StaticFS 静态文件服务，需要传入 FileSystem
func (group *RouterGroup) StaticFS(relativePath string, fs http.FileSystem) {
	if strings.Contains(relativePath, ":") || strings.Contains(relativePath, "*") {
		panic("URL parameters can not be used when serving a static folder")
	}
	handler := group.createStaticHandler(relativePath, fs)
	urlPattern := path.Join(relativePath, "/*filepath")

	// Register GET and HEAD handlers
	group.GET(urlPattern, handler)
	group.HEAD(urlPattern, handler)
}

func (group *RouterGroup) createStaticHandler(relativePath string, fs http.FileSystem) Handler {
	absolutePath := JoinPath(group.path, relativePath)
	fileServer := http.StripPrefix(absolutePath, http.FileServer(fs))

	return func(c *Context) {
		rPath := c.Request.URL.Path
		// 处理不允许访问目录
		if rPath == absolutePath+"/" || rPath == absolutePath {
			group.engine.Handler404(c.Writer, c.Request)
			return
		}
		fileServer.ServeHTTP(c.Writer, c.Request)
	}
}
