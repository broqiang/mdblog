package app

import (
	"github.com/broqiang/mdblog/app/config"
)

const (
	// VERSION 版本信息
	VERSION = "1.0.0"
)

// Cfg 加载配置文件
var Cfg = config.New()

// SysLog 系统日志输出的 Writer
var SysLog = config.NewSysLog()

// Init 初始化
func Init() {
	// 初始化配置
	// engine := gin.Default()
}
