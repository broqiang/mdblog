package config

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
