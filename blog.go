package main

import (
	"fmt"

	"github.com/broqiang/mdblog/app"
	"github.com/gin-gonic/gin"
)

func main() {

	fmt.Println(app.Cfg)

	fmt.Println(app.Cfg.Addr())
	fmt.Println(app.Cfg.Addr())
	fmt.Println(app.Cfg.Addr())
	fmt.Println(app.Cfg.Addr())
	fmt.Println(app.Cfg.Addr())
	fmt.Println(app.Cfg.Addr())

	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.String(200, app.Cfg.Addr())
	})
	r.Run() // 监听并在 0.0.0.0:8080 上启动服务
}
