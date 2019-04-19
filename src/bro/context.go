package bro

import (
	"net/http"
)

// Context HTTP 上下文
type Context struct {
	Request  *http.Request
	Writer   http.ResponseWriter
	Params   Params
	handlers Handlers
	engine   *Engine
	index    int8
	Keys     map[string]interface{}
}

// Next 遍历调用 handlers
func (c *Context) Next() {
	c.index++
	for s := int8(len(c.handlers)); c.index < s; c.index++ {
		c.handlers[c.index](c)
	}
}

// HandlerName 获取注册的 Handler 名称
func (c *Context) HandlerName() string {
	return nameOfFunction(c.handlers.Last())
}

// Handler 获取注册的 Handler， 最后一个是最终的 Handler ，前面的都是中间件
func (c *Context) Handler() Handler {
	return c.handlers.Last()
}

// Set 用于保存 Context 中可以使用的一个值
func (c *Context) Set(key string, value interface{}) {
	if c.Keys == nil {
		c.Keys = make(map[string]interface{})
	}
	c.Keys[key] = value
}

// Get 可以获取 Set 的值，如果不存在会返回一个 nil, false
func (c *Context) Get(key string) (value interface{}, exists bool) {
	value, exists = c.Keys[key]
	return
}

// Param 可以通过 key 获取路由参数， 如果不存在返回一个空字符串
func (c *Context) Param(key string) string {
	return c.Params.ByName(key)
}

// WriterString 可以直接向 Writer 中写入字符串
func (c *Context) WriterString(str string) {
	c.Writer.Write([]byte(str))
}
