package controllers

import (
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/broqiang/mdblog/app/config"
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

	c.HTML(200, "posts/index.html", params)
}

// About 关于控制器
func About(c *gin.Context) {
	// about 页面就直接展示的项目根目录下的 README.md
	path := filepath.Join(config.Root, "README.md")
	about, err := ioutil.ReadFile(path)

	if err != nil {
		c.Redirect(307, "/errors")
		return
	}

	params := mergeH(c, gin.H{
		"about": string(about),
	})

	c.HTML(200, "layouts/about.html", params)
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

// ToKeywords 用都好分割，拼接关键词
func ToKeywords(works ...string) string {
	return strings.Join(works, ",")
}
