package assets

import (
	"embed"
)

// Templates 嵌入所有模板文件
//
//go:embed web/templates
var Templates embed.FS

// Static 嵌入所有静态文件
//
//go:embed web/static
var Static embed.FS
