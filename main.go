package main

import (
	"embed"
	"flag"
	"log"
	"os"

	"mdblog/internal/app"
)

// EmbeddedAssets 嵌入的静态资源
//
//go:embed web
var EmbeddedAssets embed.FS

var (
	postsDir = flag.String("posts", "", "posts目录路径（默认为可执行文件同级目录）")
)

func main() {
	flag.Parse()

	// 创建并启动应用
	application, err := app.NewApp(*postsDir, &EmbeddedAssets)
	if err != nil {
		log.Fatalf("创建应用失败: %v", err)
		os.Exit(1)
	}

	// 启动应用
	if err := application.Start(); err != nil {
		log.Fatalf("启动应用失败: %v", err)
	}
}
