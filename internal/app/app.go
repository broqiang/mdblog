package app

import (
	"embed"
	"log"

	"mdblog/internal/config"
	"mdblog/internal/data"
	"mdblog/internal/markdown"
	"mdblog/internal/server"
)

// App 应用结构
type App struct {
	dataManager *data.Manager
	server      *server.Server
}

// NewApp 创建新的应用实例
func NewApp(postsDir string, assets *embed.FS) (*App, error) {
	// 获取posts目录路径
	actualPostsDir, err := config.GetPostsDirectory(postsDir)
	if err != nil {
		return nil, err
	}

	// 验证posts目录是否存在
	if err := config.ValidatePostsDirectory(actualPostsDir); err != nil {
		log.Printf("posts目录不存在: %s", actualPostsDir)
		if postsDir != "" {
			log.Printf("请确保指定的posts目录存在")
		} else {
			log.Printf("请确保posts目录存在或使用 -posts 参数指定正确的路径")
		}
		return nil, err
	}

	log.Printf("使用posts目录: %s", actualPostsDir)

	// 创建数据管理器
	dataManager := data.NewManager(actualPostsDir)

	// 创建Markdown解析器
	parser := markdown.NewParser()

	// 初始化数据
	if err := dataManager.InitializeData(parser); err != nil {
		return nil, err
	}

	// 创建服务器
	srv := server.NewServerWithAssets(config.DefaultHost, config.DefaultPort, dataManager, assets)

	return &App{
		dataManager: dataManager,
		server:      srv,
	}, nil
}

// Start 启动应用
func (a *App) Start() error {
	log.Printf("启动服务器: %s:%d", config.DefaultHost, config.DefaultPort)
	return a.server.Start()
}

// GetDataManager 获取数据管理器
func (a *App) GetDataManager() *data.Manager {
	return a.dataManager
}
