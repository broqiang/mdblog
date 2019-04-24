package controllers

import (
	"github.com/broqiang/mdblog/app/mdfile"
	"github.com/gin-gonic/gin"
)

// PostByPath 查询指定 Path 文档的详细信息
func PostByPath(c *gin.Context) {
	path := c.Param("path")

	post, err := mdfile.Model.ArticleByPath(path)

	if err != nil {
		c.HTML(404, "errors/notfound.html", c.Keys)
		return
	}

	params := mergeH(c, gin.H{
		"title":    post.Title,
		"keywords": ToKeywords(post.Title, post.Path),
		"post":     post,
	})

	c.HTML(200, "posts/article.html", params)
}
