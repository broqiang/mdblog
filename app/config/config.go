// Package config 用于处理配置
package config

import (
	"fmt"
	"path"
	"sync"

	"github.com/broqiang/mdblog/app/helper"

	"github.com/BurntSushi/toml"
)

// 配置文件的名称
const filePath = "config/config.toml"

var (
	once sync.Once

	// Cfg 配置文件
	Cfg Config
)

// Config 配置文件类型
type Config struct {
	// 应用的名称
	Name string `toml:"name"`

	// 应用的域名
	URL string `toml:"url"`

	// 应用的监听地址，为空就是 0.0.0.0
	Host string `toml:"host"`

	// 应用端口号
	Port uint `toml:"port"`

	// 是否开启 debug，
	// 如果 ture ，对应的 gin 就是 debug 模式，
	// 如果 false ， 对应的 gin 就是 release 模式
	Debug bool

	// 日志相关的配置
	Log Log `toml:"log"`
}

// Log 是日志相关的配置
type Log struct {
	// 日志模式，对应这的是 LogClose ，LogFile， LogStdout
	Mode string `toml:"mode"`

	// 日志文件的保存目录
	Dir string `toml:"dir"`

	// 日志文件的格式
	Format string

	// 是否开启访问日志，一般开启，有时可以关闭，比如部署在 Nginx 后面，
	// Nginx 已经开启了，可以考虑把这个日志关闭，记录两份有点多了
	Access bool
}

// InitCfg 初始化配置文件
func initCfg() {

	once.Do(func() {
		path := configPath()
		if _, err := toml.DecodeFile(path, &Cfg); err != nil {
			helper.Panicf("cannot parse config file. %v", err)
		}
	})
}

// Addr 获取服务器需要的监听地址
func (cfg *Config) Addr() string {
	return fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
}

// 获取配置文件的路径
func configPath() string {
	if configDir == "" {
		configDir = path.Join(Root, filePath)
	}

	return configDir
}
