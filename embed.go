//go:build ignore
// +build ignore

package main

import (
	"embed"
)

// 这个文件用于嵌入静态资源到可执行文件中
// 当前因为embed路径问题暂时禁用

//go:embed web/templates web/static
var assets embed.FS

// TODO: 完善embed集成到server包中
