package controllers

import "github.com/gin-gonic/gin"

// Home 主页控制器
func Home(c *gin.Context) {
	c.HTML(200, "index.html", gin.H{})
}
