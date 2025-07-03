package server

import (
	"net/http"
	"strconv"

	"mdblog/internal/config"
	"mdblog/internal/data"

	"github.com/gin-gonic/gin"
)

// handleGetPosts 处理获取文章列表API
func (s *Server) handleGetPosts(c *gin.Context) {
	// 获取分页参数
	pageStr := c.Query("page")
	pageSizeStr := c.Query("size")

	page := 1
	pageSize := config.PageSize

	if pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	if pageSizeStr != "" {
		if ps, err := strconv.Atoi(pageSizeStr); err == nil && ps > 0 {
			pageSize = ps
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

	// 返回JSON响应
	c.JSON(http.StatusOK, gin.H{
		"posts":     posts,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

// handleGetPost 处理获取单篇文章API
func (s *Server) handleGetPost(c *gin.Context) {
	id := c.Param("id")
	// 移除开头的斜杠
	if len(id) > 0 && id[0] == '/' {
		id = id[1:]
	}

	post, exists := s.manager.GetPost(id)
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "文章未找到"})
		return
	}

	c.JSON(http.StatusOK, post)
}

// handleGetCategories 处理获取分类API
func (s *Server) handleGetCategories(c *gin.Context) {
	categories := s.manager.GetAllCategories()
	c.JSON(http.StatusOK, categories)
}

// handleGetTags 处理获取标签API
func (s *Server) handleGetTags(c *gin.Context) {
	tags := s.manager.GetAllTags()
	c.JSON(http.StatusOK, tags)
}

// handleSearch 处理搜索API
func (s *Server) handleSearch(c *gin.Context) {
	query := c.Query("q")
	pageStr := c.Query("page")
	pageSizeStr := c.Query("size")

	page := 1
	pageSize := config.PageSize

	if pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	if pageSizeStr != "" {
		if ps, err := strconv.Atoi(pageSizeStr); err == nil && ps > 0 {
			pageSize = ps
		}
	}

	result := s.manager.Search(query, page, pageSize)
	c.JSON(http.StatusOK, result)
}
