package controllers

import (
	"github.com/broqiang/mdblog/app/mdfile"
	"github.com/gin-gonic/gin"
)

// Home 主页控制器
func Home(c *gin.Context) {

	params := mergeH(c, gin.H{
		"title":    "博客列表",
		"keywords": "博客列表",
		"posts":    mdfile.Model.ArticlesAll(),
	})

	c.HTML(200, "home/index.html", params)
}

// About 关于控制器
func About(c *gin.Context) {
	panic("手动 panic")
	c.String(200, "About")
}

// MergeH 合并默认参数
func mergeH(c *gin.Context, h gin.H) gin.H {
	if c.Keys == nil {
		return h
	}

	if h == nil || len(h) == 0 {
		return c.Keys
	}

	mh := make(gin.H)

	for key, val := range c.Keys {
		mh[key] = val
	}

	for key, val := range h {
		mh[key] = val
	}

	return mh
}
