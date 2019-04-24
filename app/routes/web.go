package routes

import (
	"github.com/broqiang/mdblog/app/http/controllers"
	"github.com/broqiang/mdblog/app/http/midderware"
	"github.com/gin-gonic/gin"
)

// New 初始化路由
func New(e *gin.Engine) {
	// 注册全局的导航栏中间件
	e.Use(midderware.Navigation)

	// 前台页面组，添加右侧标签的中间件
	front := e.Group("/", midderware.Sites, midderware.Tags)
	{
		front.GET("/", controllers.Home)
	}
}
