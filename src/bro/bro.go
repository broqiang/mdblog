package bro

import "net/http"

// Engine 主引擎
type Engine struct {
	// 路由组
	RouterGroup

	// 保存 Handler 的路由树
	trees map[string]*node

	RedirectTrailingSlash bool

	RedirectFixedPath bool

	NotFound http.Handler
}

// New 初始化引擎
func New() *Engine {
	engine := Engine{
		RouterGroup: RouterGroup{
			Handlers: nil,
			path:     "/",
			root:     true,
		},
		RedirectFixedPath:     true,
		RedirectTrailingSlash: true,
	}

	engine.trees = make(map[string]*node, 7)
	engine.RouterGroup.engine = &engine

	return &engine
}

func (engine *Engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path

	if root := engine.trees[r.Method]; root != nil {
		handlers, ps, tsr := root.getValue(path)
		c := engine.createContext(w, r, ps, handlers)

		if handlers != nil {
			c.Next()
			return
		}

		if r.Method != "CONNECT" && path != "/" {
			if tsr && engine.RedirectTrailingSlash {
				redirectTrailingSlash(c)
				return
			}
			if engine.RedirectFixedPath && redirectFixedPath(c, root, engine.RedirectFixedPath) {
				return
			}
		}

	}

	engine.Handler404(w, r)

}

func redirectTrailingSlash(c *Context) {
	req := c.Request
	path := req.URL.Path
	code := PrepareRedirectCode(req.Method)

	req.URL.Path = path + "/"
	if length := len(path); length > 1 && path[length-1] == '/' {
		req.URL.Path = path[:length-1]
	}
	http.Redirect(c.Writer, req, req.URL.String(), code)
}

func redirectFixedPath(c *Context, root *node, trailingSlash bool) bool {
	req := c.Request
	path := req.URL.Path

	if fixedPath, ok := root.findCaseInsensitivePath(CleanPath(path), trailingSlash); ok {
		code := PrepareRedirectCode(req.Method)

		req.URL.Path = string(fixedPath)
		http.Redirect(c.Writer, req, req.URL.String(), code)

		return true
	}
	return false
}

// Run 启动服务
func (engine *Engine) Run(addr string) (err error) {
	defer func() {
		SystemLogError(err)
	}()

	SystemLogf("Listening and serving HTTP on %s\n", addr)

	err = http.ListenAndServe(addr, engine)

	return
}

// Handler404 用来处理 404
func (engine *Engine) Handler404(w http.ResponseWriter, r *http.Request) {
	if engine.NotFound != nil {
		engine.NotFound.ServeHTTP(w, r)
		return
	}

	http.NotFound(w, r)
}
