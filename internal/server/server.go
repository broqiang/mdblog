package server

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"

	"mdblog/assets"
	"mdblog/internal/config"
	"mdblog/internal/data"

	"io/fs"

	"github.com/gin-gonic/gin"
)

// Server Web服务器
type Server struct {
	host        string
	port        int
	dataManager *data.Manager
	engine      *gin.Engine
	server      *http.Server
}

// NewServer 创建新的Web服务器
func NewServer(host string, port int, manager *data.Manager) *Server {
	// 设置 Gin 为发布模式
	gin.SetMode(gin.ReleaseMode)

	engine := gin.New()

	srv := &Server{
		host:        host,
		port:        port,
		dataManager: manager,
		engine:      engine,
	}

	// 初始化模板
	srv.initTemplates()

	// 设置中间件
	srv.setupMiddleware()

	// 设置路由
	srv.setupRoutes()

	// 创建HTTP服务器
	srv.server = &http.Server{
		Addr:         fmt.Sprintf("%s:%d", host, port),
		Handler:      engine,
		ReadTimeout:  config.ReadTimeoutSeconds * time.Second,
		WriteTimeout: config.WriteTimeoutSeconds * time.Second,
	}

	return srv
}

// initTemplates 初始化模板
func (s *Server) initTemplates() {
	// 定义模板函数
	funcMap := template.FuncMap{
		"add": func(a, b int) int {
			return a + b
		},
		"subtract": func(a, b int) int {
			return a - b
		},
		"multiply": func(a, b int) int {
			return a * b
		},
		"divide": func(a, b int) int {
			if b == 0 {
				return 0
			}
			return a / b
		},
		"ceil": func(a, b int) int {
			if b == 0 {
				return 0
			}
			return (a + b - 1) / b // 向上取整
		},
		"max": func(a, b int) int {
			if a > b {
				return a
			}
			return b
		},
		"min": func(a, b int) int {
			if a < b {
				return a
			}
			return b
		},
		"seq": func(start, end int) []int {
			var result []int
			for i := start; i <= end; i++ {
				result = append(result, i)
			}
			return result
		},
		"eq": func(a, b interface{}) bool {
			return a == b
		},
		"gt": func(a, b int) bool {
			return a > b
		},
		"lt": func(a, b int) bool {
			return a < b
		},
		"len": func(v interface{}) int {
			switch val := v.(type) {
			case map[string]int:
				return len(val)
			case []string:
				return len(val)
			default:
				return 0
			}
		},
		"safeHTML": func(s string) template.HTML {
			return template.HTML(s)
		},
	}

	// 从嵌入的文件系统创建子文件系统
	templatesSubFS, err := fs.Sub(assets.Templates, "web/templates")
	if err != nil {
		log.Fatalf("无法创建模板子文件系统: %v", err)
	}

	// 解析所有模板文件
	tmpl, err := template.New("").Funcs(funcMap).ParseFS(templatesSubFS, "*/*")
	if err != nil {
		log.Fatalf("解析模板失败: %v", err)
	}

	// 设置模板到 Gin 引擎
	s.engine.SetHTMLTemplate(tmpl)
}

// setupMiddleware 设置中间件
func (s *Server) setupMiddleware() {
	// 日志中间件
	s.engine.Use(gin.Logger())

	// 恢复中间件
	s.engine.Use(gin.Recovery())

	// CORS 中间件
	s.engine.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})
}

// setupRoutes 设置路由
func (s *Server) setupRoutes() {
	// 从嵌入的文件系统创建静态文件子系统
	staticSubFS, err := fs.Sub(assets.Static, "web/static")
	if err != nil {
		log.Fatalf("无法创建静态文件子文件系统: %v", err)
	}

	// 静态文件路由 - 使用嵌入的文件系统
	s.engine.StaticFS("/static", http.FS(staticSubFS))

	// API路由组
	api := s.engine.Group("/api")
	{
		api.GET("/posts", s.handleGetPosts)
		api.GET("/posts/*id", s.handleGetPost)
		api.GET("/categories", s.handleGetCategories)
		api.GET("/tags", s.handleGetTags)
		api.GET("/search", s.handleSearch)
	}

	// Webhook路由组
	webhook := s.engine.Group("/webhook")
	{
		webhook.POST("/gitee", s.handleGiteeWebhook)
	}

	// 页面路由
	s.engine.GET("/", s.handleIndex)
	s.engine.GET("/post/*id", s.handlePostDetail)
	s.engine.GET("/category/:category", s.handleCategory)
	s.engine.GET("/tag/:tag", s.handleTag)
	s.engine.GET("/search", s.handleSearchPage)
	s.engine.GET("/about", s.handleAbout)
}

// Start 启动服务器
func (s *Server) Start() error {
	return s.server.ListenAndServe()
}

// Shutdown 关闭服务器
func (s *Server) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}

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
	allPosts := s.dataManager.GetAllPosts()

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
		"Title":      "首页",
		"Posts":      posts,
		"Total":      total,
		"Page":       page,
		"PageSize":   pageSize,
		"Categories": s.dataManager.GetAllCategories(),
		"Tags":       s.dataManager.GetAllTags(),
	})
}

// handlePostDetail 处理文章详情页
func (s *Server) handlePostDetail(c *gin.Context) {
	id := c.Param("id")
	// 移除开头的斜杠
	if len(id) > 0 && id[0] == '/' {
		id = id[1:]
	}

	post, exists := s.dataManager.GetPost(id)
	if !exists {
		c.HTML(http.StatusNotFound, "404.html", gin.H{
			"Title": "页面未找到",
		})
		return
	}

	s.renderTemplate(c, "posts/detail.html", map[string]interface{}{
		"Title":      post.Title,
		"Post":       post,
		"Categories": s.dataManager.GetAllCategories(),
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
	allPosts := s.dataManager.GetPostsByCategory(category)

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
		"Title":      category,
		"Category":   category,
		"Posts":      posts,
		"Total":      total,
		"Page":       page,
		"PageSize":   pageSize,
		"Categories": s.dataManager.GetAllCategories(),
		"Tags":       s.dataManager.GetAllTags(),
	})
}

// handleTag 处理标签页
func (s *Server) handleTag(c *gin.Context) {
	tag := c.Param("tag")

	// 获取分页参数
	page := 1
	pageSize := config.PageSize

	pageStr := c.Query("page")
	if pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	// 获取该标签的所有文章
	allPosts := s.dataManager.GetPostsByTag(tag)

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

	s.renderTemplate(c, "tag.html", map[string]interface{}{
		"Title":      tag,
		"Tag":        tag,
		"Posts":      posts,
		"Total":      total,
		"Page":       page,
		"PageSize":   pageSize,
		"Categories": s.dataManager.GetAllCategories(),
		"Tags":       s.dataManager.GetAllTags(),
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

	result := s.dataManager.Search(query, page, pageSize)

	s.renderTemplate(c, "search.html", map[string]interface{}{
		"Title":      "搜索结果",
		"Result":     result,
		"Categories": s.dataManager.GetAllCategories(),
		"Tags":       s.dataManager.GetAllTags(),
	})
}

// handleAbout 处理About页面
func (s *Server) handleAbout(c *gin.Context) {
	post, exists := s.dataManager.GetPost("about")
	if !exists {
		c.HTML(http.StatusNotFound, "404.html", gin.H{
			"Title":      "页面未找到",
			"Message":    "About页面不存在，请创建posts/about.md文件",
			"Categories": s.dataManager.GetAllCategories(),
		})
		return
	}

	data := map[string]interface{}{
		"Title":      post.Title,
		"Post":       post,
		"Categories": s.dataManager.GetAllCategories(),
		"Tags":       s.dataManager.GetAllTags(),
	}

	s.renderTemplate(c, "posts/detail.html", data)
}

// API处理器
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
	allPosts := s.dataManager.GetAllPosts()

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

func (s *Server) handleGetPost(c *gin.Context) {
	id := c.Param("id")
	// 移除开头的斜杠
	if len(id) > 0 && id[0] == '/' {
		id = id[1:]
	}

	post, exists := s.dataManager.GetPost(id)
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "文章未找到"})
		return
	}

	c.JSON(http.StatusOK, post)
}

func (s *Server) handleGetCategories(c *gin.Context) {
	categories := s.dataManager.GetAllCategories()
	c.JSON(http.StatusOK, categories)
}

func (s *Server) handleGetTags(c *gin.Context) {
	tags := s.dataManager.GetAllTags()
	c.JSON(http.StatusOK, tags)
}

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

	result := s.dataManager.Search(query, page, pageSize)
	c.JSON(http.StatusOK, result)
}

func (s *Server) handleGiteeWebhook(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Webhook received"})
}

// renderTemplate 渲染模板
func (s *Server) renderTemplate(c *gin.Context, templateName string, data interface{}) {
	c.HTML(http.StatusOK, templateName, data)
}
