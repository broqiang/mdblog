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

// PostsByCategory 根据分类显示博客列表
func PostsByCategory(c *gin.Context) {
	name := c.Param("name")

	articles := mdfile.Model.ArticlesByCategory(name)

	if len(articles) > 0 {
		params := mergeH(c, gin.H{
			"title":    "分类 | " + name,
			"keywords": name,
			"posts":    articles,
		})
		c.HTML(200, "posts/index.html", params)
		return
	}

	c.HTML(404, "errors/notfound.html", c.Keys)
}

// PostsByTag 根据标签显示博客列表
func PostsByTag(c *gin.Context) {
	name := c.Param("name")
	articles := mdfile.Model.ArticlesByTag(name)

	if len(articles) > 0 {
		params := mergeH(c, gin.H{
			"title":    "标签 | " + name,
			"keywords": name,
			"posts":    articles,
		})
		c.HTML(200, "posts/index.html", params)
		return
	}

	c.HTML(404, "errors/notfound.html", c.Keys)
}
