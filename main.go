package main

import (
	"embed"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"mdblog/internal/data"
	"mdblog/internal/markdown"
	"mdblog/internal/server"
)

// EmbeddedAssets 嵌入的静态资源
//
//go:embed web
var EmbeddedAssets embed.FS

var (
	port     = flag.Int("port", 8080, "服务器端口")
	host     = flag.String("host", "0.0.0.0", "服务器地址")
	postsDir = flag.String("posts", "", "posts目录路径（默认为可执行文件同级目录）")
)

func main() {
	flag.Parse()

	// 设置posts目录
	var actualPostsDir string
	if *postsDir != "" {
		actualPostsDir = *postsDir
	} else {
		// 获取可执行文件所在目录
		execPath, err := os.Executable()
		if err != nil {
			log.Fatalf("获取可执行文件路径失败: %v", err)
		}
		execDir := filepath.Dir(execPath)

		// 检查是否在go run模式下（临时文件目录）
		if strings.Contains(execPath, "go-build") {
			// go run模式，使用当前工作目录
			pwd, err := os.Getwd()
			if err != nil {
				log.Fatalf("获取工作目录失败: %v", err)
			}
			actualPostsDir = filepath.Join(pwd, "posts")
		} else {
			// 正常编译的可执行文件，使用可执行文件同级目录
			actualPostsDir = filepath.Join(execDir, "posts")
		}
	}

	// 检查posts目录是否存在
	if _, err := os.Stat(actualPostsDir); os.IsNotExist(err) {
		log.Printf("posts目录不存在: %s", actualPostsDir)
		log.Printf("请确保posts目录存在或使用 -posts 参数指定正确的路径")
		os.Exit(1)
	}

	log.Printf("使用posts目录: %s", actualPostsDir)

	// 创建数据管理器
	dataManager := data.NewManager(actualPostsDir)

	// 创建Markdown解析器
	parser := markdown.NewParser()

	// 初始化数据
	if err := initializeData(dataManager, parser, actualPostsDir); err != nil {
		log.Fatalf("初始化数据失败: %v", err)
	}

	// 创建并启动服务器，传入嵌入的资源
	srv := server.NewServerWithAssets(*host, *port, dataManager, &EmbeddedAssets)

	log.Printf("启动服务器: %s:%d", *host, *port)
	if err := srv.Start(); err != nil {
		log.Fatalf("启动服务器失败: %v", err)
	}
}

// initializeData 初始化数据
func initializeData(manager *data.Manager, parser *markdown.Parser, postsDir string) error {
	// 遍历posts目录
	err := filepath.Walk(postsDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
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
		post, err := parser.ParseFile(path, postsDir)
		if err != nil {
			log.Printf("解析文件失败 %s: %v", path, err)
			return nil // 继续处理其他文件
		}

		// 添加到数据管理器
		manager.UpdatePost(post)
		log.Printf("加载文章: %s", post.Title)

		return nil
	})

	if err != nil {
		return fmt.Errorf("遍历posts目录失败: %w", err)
	}

	// 清理重复数据
	manager.CleanupDuplicates()

	log.Printf("数据初始化完成，共加载 %d 篇文章", len(manager.GetData().Posts))
	return nil
}
