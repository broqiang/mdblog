package midderware

import (
	"net/http/httputil"

	"github.com/broqiang/mdblog/app/mylog"

	"github.com/gin-gonic/gin"
)

// Recovery 用于页面出现 panic 时候的恢复
func Recovery(c *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
			httprequest, _ := httputil.DumpRequest(c.Request, false)
			mylog.LogErr.Printf("\n[Recovery] %s panic recovered: \n%s", string(httprequest), err)

			c.Redirect(307, "/errors")
			return
		}
	}()
	c.Next()
}

// NotFound 当页面 404 时的处理
func NotFound(c *gin.Context) {
	c.HTML(404, "errors/notfound.html", c.Keys)
}

// Errors 是错误的页面
func Errors(c *gin.Context) {
	c.HTML(500, "errors/errors.html", c.Keys)
	return
}
