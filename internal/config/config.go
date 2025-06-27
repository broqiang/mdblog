package config

// 应用配置常量
const (
	// 服务器配置
	DefaultPort         = 8080
	DefaultHost         = "0.0.0.0"
	ReadTimeoutSeconds  = 30
	WriteTimeoutSeconds = 30

	// 文章配置
	SummaryLines = 3
	PageSize     = 10

	// Webhook配置
	WebhookBranch = "main"

	// 搜索配置
	MaxSearchResults = 100
)
