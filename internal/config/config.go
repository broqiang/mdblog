package config

import (
	"os"
	"path/filepath"
	"strings"
)

// 应用配置常量
const (
	// 服务器配置
	DefaultPort         = 8091
	DefaultHost         = "0.0.0.0"
	ReadTimeoutSeconds  = 30
	WriteTimeoutSeconds = 30

	// 文章配置
	SummaryLines = 3
	PageSize     = 10

	// Webhook配置
	WebhookBranch = "main"
	// 这里是我的 key， 实际使用要替换成自己的仓库的 key， 我这里没有做 env 处理
	// 这里暴露出来了，为了方便参考我这个项目的朋友， 请不要搞破坏，谢谢
	WebhookSecret = "c85b7544d37c89280afcf912ec70f4083b18065d9e89966cae6059798f0dadf5"

	// 开发模式：设为 true 跳过签名验证（生产环境请设为 false）
	WebhookDevMode = false // 启用签名验证

	// Git 配置
	GitUseHTTP = true // 使用 HTTP 方式拉取 Git 仓库（避免 SSH 密钥问题）
	// PostsRepoPath 已移除，实际路径通过命令行参数或默认位置（可执行文件同级 posts 目录）动态确定

	// 搜索配置
	MaxSearchResults = 100
)

// GetPostsDirectory 获取posts目录路径
// 如果指定了postsDir参数则使用，否则自动检测
func GetPostsDirectory(postsDir string) (string, error) {
	if postsDir != "" {
		return postsDir, nil
	}

	// 获取可执行文件所在目录
	execPath, err := os.Executable()
	if err != nil {
		return "", err
	}
	execDir := filepath.Dir(execPath)

	// 检查是否在go run模式下（临时文件目录）
	if strings.Contains(execPath, "go-build") {
		// go run模式，使用当前工作目录
		pwd, err := os.Getwd()
		if err != nil {
			return "", err
		}
		return filepath.Join(pwd, "posts"), nil
	}

	// 正常编译的可执行文件，使用可执行文件同级目录
	return filepath.Join(execDir, "posts"), nil
}

// ValidatePostsDirectory 验证posts目录是否存在
func ValidatePostsDirectory(postsDir string) error {
	if _, err := os.Stat(postsDir); os.IsNotExist(err) {
		return err
	}
	return nil
}
