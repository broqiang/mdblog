package routes

import (
	"github.com/broqiang/mdblog/app/http/controllers"
	"github.com/broqiang/mdblog/app/http/midderware"
	"github.com/gin-gonic/gin"
)

// New 初始化路由
func New(e *gin.Engine) {
	// 注册全局的中间件
	e.Use(gin.Logger(), midderware.Recovery, midderware.Sites, midderware.Navigation)

	// 出现错误的页面
	e.GET("/errors", midderware.Errors)

	// 404 页面
	e.NoRoute(midderware.NotFound)

	// 前台页面组，添加右侧标签的中间件
	front := e.Group("/", midderware.Tags)
	{
		front.GET("/", controllers.Home)
		front.GET("/about", controllers.About)

		// 博客文章详情
		front.GET("/posts/:path", controllers.PostByPath)

		// 根据分类显示文章
		front.GET("/c/:name", controllers.PostsByCategory)

		// 根据标签显示文章
		front.GET("/tag/:name", controllers.PostsByTag)
	}
}
