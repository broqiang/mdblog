package routes

import (
	"github.com/broqiang/mdblog/app/http/controllers"
	"github.com/gin-gonic/gin"
)

// New 初始化路由
func New(e *gin.Engine) {
	front := e.Group("/")
	{
		front.GET("/", controllers.Home)
	}
}
