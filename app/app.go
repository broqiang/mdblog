package app

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/broqiang/mdblog/app/config"
	"github.com/broqiang/mdblog/app/helper"
	"github.com/broqiang/mdblog/app/mylog"
	"github.com/broqiang/mdblog/app/routes"
	"github.com/gin-gonic/gin"
)

const (
	// VERSION 版本信息
	VERSION = "1.0.0"
)

var cfg = config.Cfg

// Init 初始化
func Init() *gin.Engine {
	// 初始化日志配置
	setDefaultConfig()

	// 初始化引擎
	engine := gin.New()

	// 初始化模板相关的内容
	loadTemplate(engine)

	// 配置路由
	routes.New(engine)

	return engine
}

func setDefaultConfig() {
	if cfg.Debug {
		gin.SetMode(gin.DebugMode)

		return
	}

	gin.SetMode(gin.ReleaseMode)
	gin.DisableConsoleColor()

	if cfg.Log.Mode == mylog.LogFile && cfg.Log.Access {
		gin.DefaultWriter = mylog.NewAccessLog()
		return
	}

	gin.DefaultWriter = ioutil.Discard
}

func loadTemplate(e *gin.Engine) {
	root := config.Root

	// 初始化静态文件路径
	e.StaticFS("/public", http.Dir(root+"/resources/static"))
	// 初始化 favicon 图标
	e.StaticFile("/favicon.ico", root+"/resources/static/favicon.ico")

	// 初始化模板中的自定义函数
	e.FuncMap = funcMap()

	// 初始化模板
	e.LoadHTMLGlob(root + "/resources/views/*/*")
}

func funcMap() map[string]interface{} {
	return map[string]interface{}{
		// markdown 转 html
		"markdowntohtml": helper.MarkdownToHTML,
		"staticpath": func(path string) string {
			return fmt.Sprintf("%s:%d/%s", cfg.URL, cfg.Port, strings.Trim(path, "/"))
		},
	}
}

// Run 启动服务
func Run(e *gin.Engine) {
	e.Run(config.Cfg.Addr())
}
