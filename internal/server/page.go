package server

import (
	"log"
	"net/http"
	"strconv"

	"mdblog/internal/config"
	"mdblog/internal/data"

	"github.com/gin-gonic/gin"
)

// handleIndex 处理首页
func (s *Server) handleIndex(c *gin.Context) {
	// 获取分页参数
	page := 1
	pageSize := config.PageSize

	pageStr := c.Query("page")
	if pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	// 获取所有文章（按更新时间倒序）
	allPosts := s.manager.GetAllPosts()

	// 过滤掉about文章，它只在About页面显示
	var filteredPosts []*data.Post
	for _, post := range allPosts {
		if post.ID != "about" {
			filteredPosts = append(filteredPosts, post)
		}
	}

	// 分页处理
	total := len(filteredPosts)
	start := (page - 1) * pageSize
	end := start + pageSize

	var posts []*data.Post
	if start >= total {
		posts = []*data.Post{}
	} else if end > total {
		posts = filteredPosts[start:total]
	} else {
		posts = filteredPosts[start:end]
	}

	// 渲染首页模板
	s.renderTemplate(c, "posts/index.html", map[string]interface{}{
		"Title":       "首页",
		"Posts":       posts,
		"Total":       total,
		"Page":        page,
		"PageSize":    pageSize,
		"Categories":  s.manager.GetAllCategories(),
		"CurrentPage": "home",
	})
}

// handlePostDetail 处理文章详情页
func (s *Server) handlePostDetail(c *gin.Context) {
	id := c.Param("id")
	// 移除开头的斜杠
	if len(id) > 0 && id[0] == '/' {
		id = id[1:]
	}

	post, exists := s.manager.GetPost(id)
	if !exists {
		c.HTML(http.StatusNotFound, "404.html", gin.H{
			"Title": "页面未找到",
		})
		return
	}

	s.renderTemplate(c, "posts/detail.html", map[string]interface{}{
		"Title":       post.Title,
		"Post":        post,
		"Categories":  s.manager.GetAllCategories(),
		"CurrentPage": "category",
		"Category":    post.Category,
	})
}

// handleCategory 处理分类页
func (s *Server) handleCategory(c *gin.Context) {
	category := c.Param("category")

	// 获取分页参数
	page := 1
	pageSize := config.PageSize

	pageStr := c.Query("page")
	if pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	// 获取该分类的所有文章
	allPosts := s.manager.GetPostsByCategory(category)

	// 添加调试日志
	log.Printf("分类 %s: 总文章数 = %d, 当前页 = %d, 页面大小 = %d", category, len(allPosts), page, pageSize)

	// 分页处理
	total := len(allPosts)
	start := (page - 1) * pageSize
	end := start + pageSize

	var posts []*data.Post
	if start >= total {
		posts = []*data.Post{}
	} else if end > total {
		posts = allPosts[start:total]
	} else {
		posts = allPosts[start:end]
	}

	log.Printf("分页结果: start = %d, end = %d, 实际返回文章数 = %d", start, end, len(posts))

	s.renderTemplate(c, "posts/category.html", map[string]interface{}{
		"Title":       category,
		"Category":    category,
		"Posts":       posts,
		"Total":       total,
		"Page":        page,
		"PageSize":    pageSize,
		"Categories":  s.manager.GetAllCategories(),
		"CurrentPage": "category",
	})
}

// handleSearchPage 处理搜索结果页
func (s *Server) handleSearchPage(c *gin.Context) {
	query := c.Query("q")
	page := 1
	pageSize := config.PageSize

	pageStr := c.Query("page")
	if pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	result := s.manager.Search(query, page, pageSize)

	s.renderTemplate(c, "search.html", map[string]interface{}{
		"Title":       "搜索结果",
		"Result":      result,
		"Categories":  s.manager.GetAllCategories(),
		"CurrentPage": "search",
	})
}

// handleAbout 处理About页面
func (s *Server) handleAbout(c *gin.Context) {
	post, exists := s.manager.GetPost("about")
	if !exists {
		c.HTML(http.StatusNotFound, "404.html", gin.H{
			"Title":       "页面未找到",
			"Message":     "About页面不存在，请创建posts/about.md文件",
			"Categories":  s.manager.GetAllCategories(),
			"CurrentPage": "404",
		})
		return
	}

	data := map[string]interface{}{
		"Title":       post.Title,
		"Post":        post,
		"Categories":  s.manager.GetAllCategories(),
		"CurrentPage": "about",
	}

	s.renderTemplate(c, "posts/detail.html", data)
}

// handle404 处理404错误
func (s *Server) handle404(c *gin.Context) {
	c.HTML(http.StatusNotFound, "404.html", gin.H{
		"Title":       "页面未找到",
		"Message":     "",
		"Categories":  s.manager.GetAllCategories(),
		"CurrentPage": "404",
	})
}
