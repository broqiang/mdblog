// 该文件将被拆分为 server.go、webhook.go、render.go、api.go、page.go
// 具体内容见后续新文件

package server

import (
	"context"
	"embed"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"mdblog/internal/config"
	"mdblog/internal/data"
	"mdblog/internal/markdown"

	"github.com/gin-gonic/gin"
)

// Server Web服务器主结构体
// 负责初始化、启动、关闭服务器，注册路由和中间件
// 其余业务逻辑拆分到其他文件
type Server struct {
	host           string
	port           int
	manager        *data.Manager
	engine         *gin.Engine
	server         *http.Server
	embeddedAssets *embed.FS // 嵌入的资源文件系统
}

// NewServer 创建新的Web服务器
func NewServer(host string, port int, manager *data.Manager) *Server {
	return NewServerWithAssets(host, port, manager, nil)
}

// NewServerWithAssets 创建支持嵌入资源的Web服务器
func NewServerWithAssets(host string, port int, manager *data.Manager, embeddedAssets *embed.FS) *Server {
	gin.SetMode(gin.ReleaseMode)
	engine := gin.New()
	srv := &Server{
		host:           host,
		port:           port,
		manager:        manager,
		engine:         engine,
		embeddedAssets: embeddedAssets,
	}
	srv.initTemplates()
	srv.setupMiddleware()
	srv.setupRoutes()
	srv.server = &http.Server{
		Addr:         fmt.Sprintf("%s:%d", host, port),
		Handler:      engine,
		ReadTimeout:  config.ReadTimeoutSeconds * time.Second,
		WriteTimeout: config.WriteTimeoutSeconds * time.Second,
	}
	return srv
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
	// 优先使用嵌入的文件系统
	if s.embeddedAssets != nil {
		log.Println("使用嵌入的静态文件系统")

		// 从嵌入的文件系统创建静态文件子系统
		staticSubFS, err := fs.Sub(*s.embeddedAssets, "web/static")
		if err != nil {
			log.Fatalf("无法创建嵌入静态文件子文件系统: %v", err)
		}

		// 静态文件路由 - 使用嵌入的文件系统
		s.engine.StaticFS("/static", http.FS(staticSubFS))
	} else {
		log.Println("使用文件系统静态文件")

		// 检查静态文件目录是否存在
		staticDir := "web/static"
		if _, err := os.Stat(staticDir); os.IsNotExist(err) {
			log.Fatalf("静态文件目录不存在: %s。请确保在项目根目录运行或使用嵌入模式", staticDir)
		}

		// 静态文件路由
		s.engine.Static("/static", "web/static")
	}

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

	// 健康检查接口
	s.engine.GET("/health", s.handleHealth)

	// 页面路由
	s.engine.GET("/", s.handleIndex)
	s.engine.GET("/post/*id", s.handlePostDetail)
	s.engine.GET("/category/:category", s.handleCategory)
	s.engine.GET("/tag/:tag", s.handleTag)
	s.engine.GET("/search", s.handleSearchPage)
	s.engine.GET("/about", s.handleAbout)

	// 404处理器 - 必须放在所有路由的最后
	s.engine.NoRoute(s.handle404)
}

// Start 启动服务器
func (s *Server) Start() error {
	return s.server.ListenAndServe()
}

// Shutdown 关闭服务器
func (s *Server) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}

// handleHealth 处理健康检查接口
func (s *Server) handleHealth(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}

// reloadData 重新加载数据
func (s *Server) reloadData() error {
	log.Printf("开始重新加载数据...")

	// 清空现有数据
	s.manager.Clear()

	// 创建Markdown解析器
	parser := markdown.NewParser()
	postsDir := s.manager.GetPostsDir()

	// 遍历posts目录
	err := filepath.Walk(postsDir, func(path string, info os.FileInfo, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}

		// 跳过目录
		if info.IsDir() {
			return nil
		}

		// 只处理.md文件
		if filepath.Ext(path) != ".md" {
			return nil
		}

		// 解析Markdown文件
		post, parseErr := parser.ParseFile(path, postsDir)
		if parseErr != nil {
			log.Printf("解析文件失败 %s: %v", path, parseErr)
			return nil // 继续处理其他文件
		}

		// 添加到数据管理器
		s.manager.UpdatePost(post)
		log.Printf("重新加载文章: %s", post.Title)

		return nil
	})

	if err != nil {
		return fmt.Errorf("遍历posts目录失败: %w", err)
	}

	// 清理重复数据
	s.manager.CleanupDuplicates()

	log.Printf("数据重新加载完成，共加载 %d 篇文章", len(s.manager.GetData().Posts))
	return nil
}
